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

package evidence

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/kardiachain/go-kardia/kai/kaidb"
	"github.com/kardiachain/go-kardia/kai/state/cstate"

	"github.com/kardiachain/go-kardia/lib/clist"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/log"
	kproto "github.com/kardiachain/go-kardia/proto/kardiachain/types"
	"github.com/kardiachain/go-kardia/types"
)

const (
	baseKeyCommitted = "evidence-committed"
	baseKeyPending   = "evidence-pending"
)

// Pool maintains a pool of valid evidence
// in an Store.
type Pool struct {
	logger log.Logger

	evidenceList *clist.CList // concurrent linked-list of evidence
	evidenceSize uint32       // amount of pending evidence

	// needed to load validators to verify evidence
	blockStore BlockStore
	stateDB    cstate.Store
	evidenceDB kaidb.Database

	// latest state
	mtx   sync.Mutex
	state cstate.LatestBlockState

	pruningHeight uint64
	pruningTime   time.Time
}

// NewPool creates an evidence pool. If using an existing evidence store,
// it will add all pending evidence to the concurrent list.
func NewPool(stateDB cstate.Store, evidenceDB kaidb.Database, blockStore BlockStore) (*Pool, error) {
	evpool := &Pool{
		stateDB:      stateDB,
		state:        stateDB.Load(),
		logger:       log.New(),
		evidenceList: clist.New(),
		blockStore:   blockStore,
		evidenceDB:   evidenceDB,
	}

	// if pending evidence already in db, in event of prior failure, then check for expiration,
	// update the size and load it back to the evidenceList
	evpool.pruningHeight, evpool.pruningTime = evpool.removeExpiredPendingEvidence()
	evList, _, err := evpool.listEvidence([]byte(baseKeyPending), -1)
	if err != nil {
		return nil, err
	}
	atomic.StoreUint32(&evpool.evidenceSize, uint32(len(evList)))
	for _, ev := range evList {
		evpool.evidenceList.PushBack(ev)
	}

	return evpool, nil
}

// PendingEvidence is used primarily as part of block proposal and returns up to maxNum of uncommitted evidence.
func (evpool *Pool) PendingEvidence(maxBytes int64) ([]types.Evidence, int64) {
	if evpool.Size() == 0 {
		return nil, 0
	}
	evidence, size, err := evpool.listEvidence([]byte(baseKeyPending), maxBytes)
	if err != nil {
		evpool.logger.Error("Unable to retrieve pending evidence", "err", err)
	}
	return evidence, size
}

// Update pulls the latest state to be used for expiration and evidence params and then prunes all expired evidence
func (evpool *Pool) Update(state cstate.LatestBlockState, ev types.EvidenceList) {
	// sanity check
	if state.LastBlockHeight <= evpool.state.LastBlockHeight {
		panic(fmt.Sprintf(
			"Failed EvidencePool.Update new state height is less than or equal to previous state height: %d <= %d",
			state.LastBlockHeight,
			evpool.state.LastBlockHeight,
		))
	}
	evpool.logger.Info("Updating evidence pool", "last_block_height", state.LastBlockHeight,
		"last_block_time", state.LastBlockTime)

	// update the state
	evpool.updateState(state)

	evpool.markEvidenceAsCommitted(ev)

	// prune pending evidence when it has expired. This also updates when the next evidence will expire
	if evpool.Size() > 0 && state.LastBlockHeight > evpool.pruningHeight &&
		state.LastBlockTime.After(evpool.pruningTime) {
		evpool.pruningHeight, evpool.pruningTime = evpool.removeExpiredPendingEvidence()
	}
}

// markEvidenceAsCommitted processes all the evidence in the block, marking it as
// committed and removing it from the pending database.
func (evpool *Pool) markEvidenceAsCommitted(evidence types.EvidenceList) {
	blockEvidenceMap := make(map[string]struct{}, len(evidence))
	for _, ev := range evidence {
		if evpool.isPending(ev) {
			evpool.removePendingEvidence(ev)
			blockEvidenceMap[evMapKey(ev)] = struct{}{}
		}

		// Add evidence to the committed list. As the evidence is stored in the block store
		// we only need to record the height that it was saved at.
		key := keyCommitted(ev)

		h := gogotypes.UInt64Value{Value: ev.Height()}
		evBytes, err := proto.Marshal(&h)
		if err != nil {
			evpool.logger.Error("failed to marshal committed evidence", "err", err, "key(height/hash)", key)
			continue
		}

		if err := evpool.evidenceDB.Put(key, evBytes); err != nil {
			evpool.logger.Error("Unable to save committed evidence", "err", err, "key(height/hash)", key)
		}
	}

	// remove committed evidence from the clist
	if len(blockEvidenceMap) != 0 {
		evpool.removeEvidenceFromList(blockEvidenceMap)
	}
}

func (evpool *Pool) Size() uint32 {
	return atomic.LoadUint32(&evpool.evidenceSize)
}

// IsExpired checks whether evidence or a polc is expired by checking whether a height and time is older
// than set by the evidence consensus parameters
func (evpool *Pool) isExpired(height uint64, time time.Time) bool {
	var (
		params       = evpool.State().ConsensusParams.Evidence
		ageDuration  = evpool.State().LastBlockTime.Sub(time)
		ageNumBlocks = evpool.State().LastBlockHeight - height
	)
	return ageNumBlocks > uint64(params.MaxAgeNumBlocks) &&
		ageDuration > params.MaxAgeDuration
}

// IsPending checks whether the evidence is already pending. DB errors are passed to the logger.
func (evpool *Pool) isPending(evidence types.Evidence) bool {
	key := keyPending(evidence)
	ok, err := evpool.evidenceDB.Has(key)
	if err != nil {
		evpool.logger.Error("Unable to find pending evidence", "err", err)
	}
	return ok
}

// IsCommitted returns true if we have already seen this exact evidence and it is already marked as committed.
func (evpool *Pool) isCommitted(evidence types.Evidence) bool {
	key := keyCommitted(evidence)
	ok, err := evpool.evidenceDB.Has(key)
	if err != nil {
		evpool.logger.Error("Unable to find committed evidence", "err", err)
	}
	return ok
}

// EvidenceFront ...
func (evpool *Pool) EvidenceFront() *clist.CElement {
	return evpool.evidenceList.Front()
}

// EvidenceWaitChan ...
func (evpool *Pool) EvidenceWaitChan() <-chan struct{} {
	return evpool.evidenceList.WaitChan()
}

// SetLogger sets the Logger.
func (evpool *Pool) SetLogger(l log.Logger) {
	evpool.logger = l
}

func (evpool *Pool) updateState(state cstate.LatestBlockState) {
	evpool.mtx.Lock()
	defer evpool.mtx.Unlock()
	evpool.state = state
}

// State returns the current state of the evpool.
func (evpool *Pool) State() cstate.LatestBlockState {
	evpool.mtx.Lock()
	defer evpool.mtx.Unlock()
	return evpool.state
}

func (evpool *Pool) removeEvidenceFromList(
	blockEvidenceMap map[string]struct{}) {

	for e := evpool.evidenceList.Front(); e != nil; e = e.Next() {
		// Remove from clist
		ev := e.Value.(types.Evidence)
		if _, ok := blockEvidenceMap[evMapKey(ev)]; ok {
			evpool.evidenceList.Remove(e)
			e.DetachPrev()
		}
	}
}

func (evpool *Pool) removePendingEvidence(evidence types.Evidence) {
	key := keyPending(evidence)
	if err := evpool.evidenceDB.Delete(key); err != nil {
		evpool.logger.Error("Unable to delete pending evidence", "err", err)
	} else {
		atomic.AddUint32(&evpool.evidenceSize, ^uint32(0))
		evpool.logger.Info("Deleted pending evidence", "evidence", evidence)
	}
}

// listEvidence retrieves lists evidence from oldest to newest within maxBytes.
// If maxBytes is -1, there's no cap on the size of returned evidence.
func (evpool *Pool) listEvidence(prefixKey []byte, maxBytes int64) ([]types.Evidence, int64, error) {
	var evidence []types.Evidence
	var evList kproto.EvidenceData // used for calculating the bytes size
	var evSize int64
	var totalSize int64
	iter := evpool.evidenceDB.NewIterator(prefixKey, nil)
	for iter.Next() {
		var evp kproto.Evidence
		if err := evp.Unmarshal(iter.Value()); err != nil {
			return evidence, totalSize, err
		}

		evList.Evidence = append(evList.Evidence, evp)
		evSize = int64(evList.Size())
		if maxBytes != -1 && evSize > maxBytes {
			if err := iter.Error(); err != nil {
				return evidence, totalSize, err
			}
			return evidence, totalSize, nil
		}
		ev, err := types.EvidenceFromProto(&evp)
		if err != nil {
			return nil, totalSize, err
		}

		totalSize = evSize
		evidence = append(evidence, ev)
	}
	return evidence, totalSize, nil
}

func (evpool *Pool) removeExpiredPendingEvidence() (uint64, time.Time) {
	iter := evpool.evidenceDB.NewIterator([]byte(baseKeyPending), nil)
	blockEvidenceMap := make(map[string]struct{})
	for iter.Next() {
		ev, err := bytesToEv(iter.Value())
		if err != nil {
			evpool.logger.Error("Error in transition evidence from protobuf", "err", err)
			continue
		}

		if !evpool.isExpired(ev.Height(), ev.Time()) {
			if len(blockEvidenceMap) != 0 {
				evpool.removeEvidenceFromList(blockEvidenceMap)
			}
			maxAgeNumBlocks := evpool.State().ConsensusParams.Evidence.MaxAgeNumBlocks
			// return the height and time with which this evidence will have expired so we know when to prune next
			return ev.Height() + uint64(maxAgeNumBlocks) + 1,
				ev.Time().Add(evpool.State().ConsensusParams.Evidence.MaxAgeDuration).Add(time.Second)
		}
		evpool.removePendingEvidence(ev)
		blockEvidenceMap[evMapKey(ev)] = struct{}{}
	}
	// We either have no pending evidence or all evidence has expired
	if len(blockEvidenceMap) != 0 {
		evpool.removeEvidenceFromList(blockEvidenceMap)
	}
	return evpool.State().LastBlockHeight, evpool.State().LastBlockTime
}

// AddEvidence checks the evidence is valid and adds it to the pool.
func (evpool *Pool) AddEvidence(ev types.Evidence) error {
	evpool.logger.Debug("Attempting to add evidence", "ev", ev)

	// We have already verified this piece of evidence - no need to do it again
	if evpool.isPending(ev) {
		evpool.logger.Info("Evidence already pending, ignoring this one", "ev", ev)
		return nil
	}

	// check that the evidence isn't already committed
	if evpool.isCommitted(ev) {
		// this can happen if the peer that sent us the evidence is behind so we shouldn't
		// punish the peer.
		evpool.logger.Debug("Evidence was already committed, ignoring this one", "ev", ev)
		return nil
	}

	if err := evpool.verify(ev); err != nil {
		return types.NewErrInvalidEvidence(ev, err)
	}

	// 2) Save to store.
	if err := evpool.addPendingEvidence(ev); err != nil {
		return fmt.Errorf("can't add evidence to pending list: %w", err)
	}

	// 3) Add evidence to clist.
	evpool.evidenceList.PushBack(ev)

	evpool.logger.Info("Verified new evidence of byzantine behaviour", "evidence", ev)
	return nil
}

// AddEvidenceFromConsensus should be exposed only to the consensus so it can add evidence to the pool
// directly without the need for verification.
func (evpool *Pool) AddEvidenceFromConsensus(ev types.Evidence) error {
	// we already have this evidence, log this but don't return an error.
	if evpool.isPending(ev) {
		evpool.logger.Info("Evidence already pending, ignoring this one", "ev", ev)
		return nil
	}

	if err := evpool.addPendingEvidence(ev); err != nil {
		return fmt.Errorf("can't add evidence to pending list: %w", err)
	}
	// add evidence to be gossiped with peers
	evpool.evidenceList.PushBack(ev)

	evpool.logger.Info("Verified new evidence of byzantine behavior", "evidence", ev)

	return nil
}

// CheckEvidence takes an array of evidence from a block and verifies all the evidence there.
// If it has already verified the evidence then it jumps to the next one. It ensures that no
// evidence has already been committed or is being proposed twice. It also adds any
// evidence that it doesn't currently have so that it can quickly form ABCI Evidence later.
func (evpool *Pool) CheckEvidence(evList types.EvidenceList) error {
	hashes := make([]common.Hash, len(evList))
	for idx, ev := range evList {
		ok := evpool.fastCheck(ev)

		if !ok {
			// check that the evidence isn't already committed
			if evpool.isCommitted(ev) {
				return types.NewErrInvalidEvidence(ev, errors.New("evidence was already committed"))
			}

			if err := evpool.verify(ev); err != nil {
				return types.NewErrInvalidEvidence(ev, err)
			}

			if err := evpool.addPendingEvidence(ev); err != nil {
				// Something went wrong with adding the evidence but we already know it is valid
				// hence we log an error and continue
				evpool.logger.Error("Can't add evidence to pending list", "err", err, "ev", ev)
			}

			evpool.logger.Info("Verified new evidence of byzantine behavior", "evidence", ev)
		}

		// check for duplicate evidence. We cache hashes so we don't have to work them out again.
		hashes[idx] = ev.Hash()
		for i := idx - 1; i >= 0; i-- {
			if hashes[i].Equal(hashes[idx]) {
				return types.NewErrInvalidEvidence(ev, errors.New("duplicate evidence"))
			}
		}
	}
	return nil
}

func (evpool *Pool) addPendingEvidence(ev types.Evidence) error {
	evpb, err := types.EvidenceToProto(ev)
	if err != nil {
		return fmt.Errorf("unable to convert to proto, err: %w", err)
	}

	evBytes, err := evpb.Marshal()
	if err != nil {
		return fmt.Errorf("unable to marshal evidence: %w", err)
	}

	key := keyPending(ev)

	err = evpool.evidenceDB.Put(key, evBytes)
	if err != nil {
		return fmt.Errorf("can't persist evidence: %w", err)
	}
	atomic.AddUint32(&evpool.evidenceSize, 1)
	return nil
}

// fastCheck leverages the fact that the evidence pool may have already verified the evidence to see if it can
// quickly conclude that the evidence is already valid.
func (evpool *Pool) fastCheck(ev types.Evidence) bool {
	// for all other evidence the evidence pool just checks if it is already in the pending db
	return evpool.isPending(ev)
}

func bytesToEv(evBytes []byte) (types.Evidence, error) {
	var evpb kproto.Evidence
	err := evpb.Unmarshal(evBytes)
	if err != nil {
		return &types.DuplicateVoteEvidence{}, err
	}

	return types.EvidenceFromProto(&evpb)
}

func evMapKey(ev types.Evidence) string {
	return ev.Hash().String()
}

// big endian padded hex
func bE(h uint64) string {
	return fmt.Sprintf("%0.16X", h)
}

func keyCommitted(evidence types.Evidence) []byte {
	return append([]byte(baseKeyCommitted), keySuffix(evidence)...)
}

func keyPending(evidence types.Evidence) []byte {
	return append([]byte(baseKeyPending), keySuffix(evidence)...)
}

func keySuffix(evidence types.Evidence) []byte {
	return []byte(fmt.Sprintf("%s/%X", bE(evidence.Height()), evidence.Hash()))
}
