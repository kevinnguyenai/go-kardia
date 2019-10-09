/*
 *  Copyright 2018 KardiaChain
 *  This file is part of the go-kardia library.
 *
 *  The go-kardia library is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  The go-kardia library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU Lesser General Public License for more details.
 *
 *  You should have received a copy of the GNU Lesser General Public License
 *  along with the go-kardia library. If not, see <http://www.gnu.org/licenses/>.
 */

package tx_pool

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/kardiachain/go-kardia/kai/events"
	"github.com/kardiachain/go-kardia/kai/state"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/event"
	"github.com/kardiachain/go-kardia/lib/log"
	"github.com/kardiachain/go-kardia/types"
)

const (
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10

	// promotableQueueSize is the size for promotableQueue
	promotableQueueSize = 1000000
)

var (
	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")
)

// blockChain provides the state of blockchain and current gas limit to do
// some pre checks in tx pool and event subscribers.
type blockChain interface {
	CurrentBlock() *types.Block
	GetBlock(hash common.Hash, number uint64) *types.Block
	StateAt(root common.Hash) (*state.StateDB, error)
	DB() types.Database
	SubscribeChainHeadEvent(ch chan<- events.ChainHeadEvent) event.Subscription
}

// TxPoolConfig are the configuration parameters of the transaction pool.
type TxPoolConfig struct {
	NoLocals  bool          // Whether local transaction handling should be disabled
	Journal   string        // Journal of local transactions to survive node restarts
	Rejournal time.Duration // Time interval to regenerate the local transaction journal
	GlobalSlots  uint64 // Maximum number of executable transaction slots for all accounts
	GlobalQueue  uint64 // Maximum number of non-executable transaction slots for all accounts
	NumberOfWorkers int
	WorkerCap       int
	BlockSize       int
}

// DefaultTxPoolConfig contains the default configurations for the transaction
// pool.
var DefaultTxPoolConfig = TxPoolConfig{
	Journal:   "transactions.rlp",
	Rejournal: time.Hour,

	GlobalSlots:  64,
	GlobalQueue:  5120000,

	NumberOfWorkers: 3,
	WorkerCap: 512,
	BlockSize: 7192,
}

// GetDefaultTxPoolConfig returns default txPoolConfig with given dir path
func GetDefaultTxPoolConfig(path string) *TxPoolConfig {
	conf := DefaultTxPoolConfig
	return &conf
}

// TxPool contains all currently known transactions. Transactions
// enter the pool when they are received from the network or submitted
// locally. They exit the pool when they are included in the blockchain.
//
// The pool separates processable transactions (which can be applied to the
// current state) and future transactions. Transactions move between those
// two states over time as they are received and processed.
type TxPool struct {
	logger log.Logger

	config       TxPoolConfig
	chainconfig  *types.ChainConfig
	chain        blockChain
	txFeed       event.Feed
	scope        event.SubscriptionScope

	// txsCh is used for pending txs
	txsCh       chan []interface{}

	// allCh is used to cache all processed txs
	allCh       chan []interface{}

	chainHeadCh  chan events.ChainHeadEvent
	chainHeadSub event.Subscription
	mu           sync.RWMutex

	numberOfWorkers int
	workerCap       int

	currentState *state.StateDB      // Current state in the blockchain head
	pendingState *state.ManagedState // Pending state tracking virtual nonces
	addressState map[common.Address]uint64 // address state will cache current state of addresses with the latest nonce

	currentMaxGas uint64 // Current gas limit for transaction caps
	totalPendingGas uint64

	journal *txJournal  // Journal of local transaction to back up to disk

	pendingSize uint // pendingSize is a counter, increased when adding new txs, decreased when remove txs
	pending  map[common.Address]types.Transactions   // All currently processable transactions
	all      *common.Set                        // All transactions to allow lookups
	promotableQueue *common.Set                 // a queue of addresses that are waiting for processing txs (FIFO)

	wg sync.WaitGroup // for shutdown sync
}

// NewTxPool creates a new transaction pool to gather, sort and filter inbound
// transactions from the network.
func NewTxPool(logger log.Logger, config TxPoolConfig, chainconfig *types.ChainConfig, chain blockChain) *TxPool {
	// Sanitize the input to ensure no vulnerable gas prices are set
	//config = (&config).sanitize(logger)

	// Create the transaction pool with its initial settings
	pool := &TxPool{
		logger:      logger,
		config:      config,
		chainconfig: chainconfig,
		chain:       chain,
		pending:     make(map[common.Address]types.Transactions),
		all:         common.NewSet(int64(config.GlobalQueue)),
		promotableQueue: common.NewSet(promotableQueueSize),
		addressState: make(map[common.Address]uint64),
		chainHeadCh: make(chan events.ChainHeadEvent, chainHeadChanSize),
		totalPendingGas: uint64(0),
		txsCh: make(chan []interface{}, 100),
		allCh: make(chan []interface{}),
		numberOfWorkers: config.NumberOfWorkers,
		workerCap: config.WorkerCap,
		pendingSize: 0,
	}
	//pool.priced = newTxPricedList(logger, pool.all)
	pool.reset(nil, chain.CurrentBlock().Header())

	/*if !config.NoLocals && config.Journal != "" {
		pool.journal = newTxJournal(logger, config.Journal)

		if err := pool.journal.load(pool.AddLocals); err != nil {
			logger.Warn("Failed to load transaction journal", "err", err)
		}
	}*/

	// Subscribe events from blockchain
	pool.chainHeadSub = pool.chain.SubscribeChainHeadEvent(pool.chainHeadCh)

	// Start the event loop and return
	pool.wg.Add(1)
	go pool.loop()

	return pool
}

// loop is the transaction pool's main event loop, waiting for and reacting to
// outside blockchain events as well as for various reporting and transaction
// eviction events.
func (pool *TxPool) loop() {
	// Track the previous head headers for transaction reorgs
	head := pool.chain.CurrentBlock()
	collectTicker := time.NewTicker(2000 * time.Millisecond)
	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle ChainHeadEvent
		case ev := <-pool.chainHeadCh:
			go pool.reset(head.Header(), ev.Block.Header())
		// Be unsubscribed due to system stopped
		case <-pool.chainHeadSub.Err():
			return
		case <-collectTicker.C:
			go pool.collectTxs()
		}
	}
}

// collectTxs is called periodically to add txs from txsCh to pending pool
func (pool *TxPool) collectTxs() {
	for i := 0; i < pool.numberOfWorkers; i++ {
		go pool.work(i, <-pool.txsCh)
	}
}

// work is called by workers to add txs into pool
func (pool *TxPool) work(index int, txs []interface{}) {
	go pool.addTxs(txs)
}

func (pool *TxPool) AddTxs(txs []interface{}) {
	if len(txs) > 0 {
		to := pool.workerCap
		if len(txs) < to {
			to = len(txs)
		}
		pool.txsCh <- txs[0:to]
		go pool.AddTxs(txs[to:])
	}
}

func (pool *TxPool) ResetWorker(workers int, cap int) {
	pool.numberOfWorkers = workers
	pool.workerCap = cap
}

// ClearPending is used to clear pending data. Note: this function is only for testing only
func (pool *TxPool) ClearPending() {
	pool.pending = make(map[common.Address]types.Transactions)
	pool.pendingSize = 0
}

// lockedReset is a wrapper around reset to allow calling it in a thread safe
// manner. This method is only ever used in the tester!
func (pool *TxPool) lockedReset(oldHead, newHead *types.Header) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.reset(oldHead, newHead)
}

// reset retrieves the current state of the blockchain and ensures the content
// of the transaction pool is valid with regard to the chain state.
func (pool *TxPool) reset(oldHead, newHead *types.Header) {
	// Initialize the internal state to the current head
	currentBlock := pool.chain.CurrentBlock()

	if newHead == nil {
		newHead = currentBlock.Header() // Special case during testing
	}

	statedb, err := pool.chain.StateAt(currentBlock.Root())
	pool.logger.Info("TxPool reset state to new head block", "height", newHead.Height, "root", newHead.Root)
	if err != nil {
		pool.logger.Error("Failed to reset txpool state", "err", err)
		return
	}
	pool.currentState = statedb
	pool.pendingState = state.ManageState(statedb)
	pool.currentMaxGas = newHead.GasLimit

	// remove current block's txs from pending
	pool.RemoveTxs(currentBlock.Transactions())
	go pool.saveTxs(currentBlock.Transactions())
}

// Stop terminates the transaction pool.
func (pool *TxPool) Stop() {
	// Unsubscribe all subscriptions registered from txpool
	pool.scope.Close()

	// Unsubscribe subscriptions registered from blockchain
	pool.chainHeadSub.Unsubscribe()
	pool.wg.Wait()
	pool.logger.Info("Transaction pool stopped")
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (pool *TxPool) SubscribeNewTxsEvent(ch chan<- events.NewTxsEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}

// State returns the virtual managed state of the transaction pool.
func (pool *TxPool) State() *state.ManagedState {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.pendingState
}

// ProposeTransactions collects transactions from pending and remove them.
func (pool *TxPool) ProposeTransactions() types.Transactions {
	txs, _ := pool.Pending(pool.config.BlockSize, true)
	return txs
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// GetAddressState gets address's nonce based on statedb or current addressState in txPool
func (pool *TxPool) GetAddressState(address common.Address) uint64 {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	nonce := pool.currentState.GetNonce(address) + 1
	if _, ok := pool.addressState[address]; !ok {
		return nonce
	} else if nonce > pool.addressState[address] {
		pool.addressState[address] = nonce
	}
	return pool.addressState[address] + 1
}

// Pending collects pending transactions with limit number, if removeResult is marked to true then remove results after all.
func (pool *TxPool) Pending(limit int, removeResult bool) (types.Transactions, error) {

	pool.mu.Lock()
	defer pool.mu.Unlock()

	pending := make(types.Transactions, 0)
	addedTx := make(map[common.Hash]struct{})
	promotableAddresses := pool.promotableQueue.List()

	// loop through promotableAddresses to get addresses come first
	for _, addrInterface := range promotableAddresses {

		if len(pending) >= limit {
			break
		}

		addr := addrInterface.(common.Address)
		txs := pool.pending[addr]

		if len(txs) > 0 {
			// latest txs must be the highest nonce
			// update addressState here
			pool.addressState[addr] = txs[len(txs)-1].Nonce()
			for _, tx := range txs {

				if _, ok := addedTx[tx.Hash()]; ok {
					continue
				}

				if pool.all.Has(tx.Hash()) {
					continue
				}

				pending = append(pending, tx)
				addedTx[tx.Hash()] = struct{}{}
			}
			// delete all txs in address if removeResult is true
			if removeResult {
				delete(pool.pending, addr)
				pool.pendingSize -= uint(len(txs))

				// remove addr from queue
				pool.promotableQueue.Remove(addrInterface)
			}
		}
	}
	// sort pending list
	return pending, nil
}

// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) ValidateTx(tx *types.Transaction) (*common.Address, error) {

	// check sender and duplicated pending tx
	sender, err := pool.getSender(tx)
	if err != nil {
		return nil, err
	}

	if pool.pending[*sender] != nil && uint64(len(pool.pending[*sender])) >= pool.config.GlobalSlots {
		return nil, fmt.Errorf("%v has reached its limit %v/%v", sender.Hex(), len(pool.pending[*sender]), pool.config.GlobalSlots)
	}

	if tx.Nonce() <= pool.addressState[*sender] && pool.addressState[*sender] > 0 {
		return nil, fmt.Errorf("invalid nonce with sender %v %v <= %v", sender.Hex(), tx.Nonce(), pool.addressState[*sender])
	}

	// if tx has been added into db then reject it
	if pool.all.Has(tx.Hash()) {
		return nil, fmt.Errorf("transaction %v existed", tx.Hash().Hex())
	}

	return sender, nil
}


// AddTx enqueues a single transaction into the pool if it is valid
func (pool *TxPool) AddTx(tx *types.Transaction) error {
	if err := pool.addTx(tx); err != nil {
		return err
	}
	_, err := types.Sender(tx)
	if err != nil {
		return ErrInvalidSender
	}
	go pool.txFeed.Send(events.NewTxsEvent{Txs: []*types.Transaction{tx}})
	return nil
}

// getSender gets transaction's sender
func (pool *TxPool) getSender(tx *types.Transaction) (*common.Address, error) {
	sender, err := types.Sender(tx)
	if err != nil {
		return nil, ErrInvalidSender
	}
	return &sender, nil
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *TxPool) addTx(tx *types.Transaction) error {

	sender, err := pool.ValidateTx(tx)
	if err != nil {
		return err
	}

	// update address state
	if nonce, ok := pool.addressState[*sender]; !ok || nonce < tx.Nonce() {
		pool.logger.Info("update nonce", "address", sender.Hex(), "nonce", tx.Nonce(), "currentNonce", nonce)
		pool.addressState[*sender] = tx.Nonce()
	}

	pendingTxs := pool.pending[*sender]
	if pendingTxs == nil {
		pendingTxs = make(types.Transactions, 0)
	}

	pendingTxs = append(pendingTxs, tx)
	sort.Sort(types.TxByNonce(pendingTxs))
	pool.pending[*sender] = pendingTxs

	// add sender to queue if it does not exist in queue
	if !pool.promotableQueue.Has(*sender) {
		pool.promotableQueue.Add(*sender)
	}

	pool.pendingSize++

	return nil
}

// addTxs attempts to queue a batch of transactions if they are valid.
func (pool *TxPool) addTxs(txs []interface{}) {
	pool.mu.Lock()
	promoted := make([]*types.Transaction, 0)
	addedTx := make(map[common.Hash]struct{})
	for _, txInterface := range txs {
		if txInterface == nil {
			continue
		}
		tx := txInterface.(*types.Transaction)

		if _, ok := addedTx[tx.Hash()]; ok {
			continue
		}

		// validate and add tx to pool
		if err := pool.addTx(tx); err == nil {
			promoted = append(promoted, tx)
			addedTx[tx.Hash()] = struct{}{}
		}
	}
	pool.mu.Unlock()

	if len(promoted) > 0 {
		go pool.txFeed.Send(events.NewTxsEvent{Txs: promoted})
	}
}

func (pool *TxPool) CachedTxs() *common.Set {
	return pool.all
}

func (pool *TxPool) saveTxs(txs []*types.Transaction) {
	if len(txs) > 0 {
		hashes := make([]interface{}, len(txs))
		for i, tx := range txs {
			hashes[i] = tx.Hash()
		}
		go pool.all.Add(hashes...)
	}
}

// RemoveTx removes transactions from pending queue.
func (pool *TxPool) RemoveTxs(txs types.Transactions) {

	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.logger.Trace("Removing Txs from pending", "txs", len(txs))
	startTime := getTime()
	for _, tx := range txs {
		sender, _ := pool.getSender(tx)
		pendings := pool.pending[*sender]

		if pendings != nil && len(pendings) > 0 {
			newTxs := make(types.Transactions, 0)
			for _, pending := range pendings {
				if pending.Nonce() != tx.Nonce() {
					newTxs = append(newTxs, pending)
				} else {
					pool.pendingSize -= 1
				}
			}
			pool.pending[*sender] = newTxs
		}
	}
	diff := getTime() - startTime
	pool.logger.Trace("total time to finish removing txs from pending", "time", diff)
}

func (pool *TxPool) PendingSize() int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	count := 0
	for _, txs := range pool.pending {
		count += txs.Len()
	}
	return count
}

func (pool *TxPool) GetPendingData() *types.Transactions {
	txs := make(types.Transactions, 0)
	for _, pendings := range pool.pending {
		txs = append(txs, pendings...)
	}
	return &txs
}

// TxInterfaceByNonce implements the sort interface to allow sorting a list of transactions (in interface)
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxInterfaceByNonce []interface{}
func (s TxInterfaceByNonce) Len() int           { return len(s) }
func (s TxInterfaceByNonce) Less(i, j int) bool {
	return s[i].(*types.Transaction).Nonce() < s[j].(*types.Transaction).Nonce()
}
func (s TxInterfaceByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
