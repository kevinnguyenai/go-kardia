package cstate

import (
	"fmt"

	"github.com/kardiachain/go-kardiamain/mainchain/genesis"

	"github.com/kardiachain/go-kardiamain/lib/common"

	"github.com/kardiachain/go-kardiamain/lib/rlp"
	"github.com/kardiachain/go-kardiamain/types"

	"github.com/kardiachain/go-kardiamain/kai/kaidb"
	"github.com/kardiachain/go-kardiamain/lib/math"
)

const (
	// persist validators every valSetCheckpointInterval blocks to avoid
	// LoadValidators taking too much time.
	// https://github.com/tendermint/tendermint/pull/3438
	// 100000 results in ~ 100ms to get 100 validators (see BenchmarkLoadValidators)
	valSetCheckpointInterval = 100000
)

//------------------------------------------------------------------------

func calcValidatorsKey(height uint64) []byte {
	return []byte(fmt.Sprintf("validatorsKey:%v", height))
}

func calcConsensusParamsKey(height uint64) []byte {
	return []byte(fmt.Sprintf("consensusParamsKey:%v", height))
}

// LoadStateFromDBOrGenesisDoc loads the most recent state from the database,
// or creates a new one from the given genesisDoc and persists the result
// to the database.
func LoadStateFromDBOrGenesisDoc(stateDB kaidb.Database, genesisDoc *genesis.Genesis) (LastestBlockState, error) {
	state := LoadState(stateDB)
	if state.IsEmpty() {
		var err error
		state, err = MakeGenesisState(genesisDoc)
		if err != nil {
			return state, err
		}
		SaveState(stateDB, state)
	}

	return state, nil
}

// SaveState persists the State, the ValidatorsInfo, and the ConsensusParamsInfo to the database.
// This flushes the writes (e.g. calls SetSync).
func SaveState(db kaidb.KeyValueStore, state LastestBlockState) {
	saveState(db, state, stateKey)
}

func saveState(db kaidb.KeyValueStore, state LastestBlockState, key []byte) {
	nextHeight := state.LastBlockHeight.Uint64() + 1
	// If first block, save validators for block 1.
	if nextHeight == 1 {
		// This extra logic due to Tendermint validator set changes being delayed 1 block.
		// It may get overwritten due to InitChain validator updates.
		lastHeightVoteChanged := uint64(1)
		saveValidatorsInfo(db, nextHeight, lastHeightVoteChanged, state.Validators)
	}
	// Save next validators.
	saveValidatorsInfo(db, nextHeight+1, state.LastHeightValidatorsChanged.Uint64(), state.NextValidators)
	// Save next consensus params.
	saveConsensusParamsInfo(db, uint64(nextHeight), state.LastHeightConsensusParamsChanged, state.ConsensusParams)
	db.Put(key, state.Bytes())
}

// LoadState loads the State from the database.
func LoadState(db kaidb.Database) LastestBlockState {
	return loadState(db, stateKey)
}

func loadState(db kaidb.Database, key []byte) (state LastestBlockState) {
	buf, _ := db.Get(key)

	if len(buf) == 0 {
		return state
	}

	err := rlp.DecodeBytes(buf, &state)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		panic(fmt.Sprintf(`LoadState: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return state
}

//-----------------------------------------------------------------------------

// ValidatorsInfo represents the latest validator set, or the last height it changed
type ValidatorsInfo struct {
	ValidatorSet      *types.ValidatorSet `rlp:"nil"`
	LastHeightChanged uint64
}

// Bytes serializes the ValidatorsInfo
func (valInfo *ValidatorsInfo) Bytes() []byte {
	b, err := rlp.EncodeToBytes(valInfo)
	if err != nil {
		panic(err)
	}
	return b
}

// LoadValidators loads the ValidatorSet for a given height.
// Returns ErrNoValSetForHeight if the validator set can't be found for this height.
func LoadValidators(db kaidb.KeyValueStore, height uint64) (*types.ValidatorSet, error) {
	valInfo := loadValidatorsInfo(db, uint64(height))
	if valInfo == nil {
		return nil, ErrNoValSetForHeight{height}
	}
	if valInfo.ValidatorSet == nil {
		lastStoredHeight := lastStoredHeightFor(height, valInfo.LastHeightChanged)
		valInfo2 := loadValidatorsInfo(db, uint64(lastStoredHeight))
		if valInfo2 == nil || valInfo2.ValidatorSet == nil {
			panic(
				fmt.Sprintf("Couldn't find validators at height %d (height %d was originally requested)",
					lastStoredHeight,
					height,
				),
			)
		}
		valInfo2.ValidatorSet.AdvanceProposer(int64(height - uint64(lastStoredHeight))) // mutate
		valInfo = valInfo2
	}

	return valInfo.ValidatorSet, nil
}

func lastStoredHeightFor(height, lastHeightChanged uint64) int64 {
	checkpointHeight := height - height%valSetCheckpointInterval
	return math.MaxInt64(int64(checkpointHeight), int64(lastHeightChanged))
}

// CONTRACT: Returned ValidatorsInfo can be mutated.
func loadValidatorsInfo(db kaidb.Database, height uint64) *ValidatorsInfo {
	buf, err := db.Get(calcValidatorsKey(height))
	if err != nil {
		panic(err)
	}
	if len(buf) == 0 {
		return nil
	}

	v := new(ValidatorsInfo)

	err = rlp.DecodeBytes(buf, v)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		panic(fmt.Sprintf(`LoadValidators: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v
}

// saveValidatorsInfo persists the validator set.
//
// `height` is the effective height for which the validator is responsible for
// signing. It should be called from s.Save(), right before the state itself is
// persisted.
func saveValidatorsInfo(db kaidb.Database, height, lastHeightChanged uint64, valSet *types.ValidatorSet) {
	if lastHeightChanged > height {
		panic("LastHeightChanged cannot be greater than ValidatorsInfo height")
	}
	valInfo := &ValidatorsInfo{
		LastHeightChanged: uint64(lastHeightChanged),
	}
	// Only persist validator set if it was updated or checkpoint height (see
	// valSetCheckpointInterval) is reached.
	if height == lastHeightChanged || height%valSetCheckpointInterval == 0 {
		valInfo.ValidatorSet = valSet
	}
	db.Put(calcValidatorsKey(height), valInfo.Bytes())
}

//-----------------------------------------------------------------------------

// ConsensusParamsInfo represents the latest consensus params, or the last height it changed
type ConsensusParamsInfo struct {
	ConsensusParams   types.ConsensusParams
	LastHeightChanged uint64
}

// LoadConsensusParams loads the ConsensusParams for a given height.
func LoadConsensusParams(db kaidb.Database, height uint64) (types.ConsensusParams, error) {
	empty := types.ConsensusParams{}

	paramsInfo := loadConsensusParamsInfo(db, height)
	if paramsInfo == nil {
		return empty, ErrNoConsensusParamsForHeight{height}
	}

	if paramsInfo.ConsensusParams.Equals(&empty) {
		paramsInfo2 := loadConsensusParamsInfo(db, paramsInfo.LastHeightChanged)
		if paramsInfo2 == nil {
			panic(
				fmt.Sprintf(
					"Couldn't find consensus params at height %d as last changed from height %d",
					paramsInfo.LastHeightChanged,
					height,
				),
			)
		}
		paramsInfo = paramsInfo2
	}

	return paramsInfo.ConsensusParams, nil
}

func loadConsensusParamsInfo(db kaidb.Database, height uint64) *ConsensusParamsInfo {
	buf, err := db.Get(calcConsensusParamsKey(uint64(height)))
	if err != nil {
		panic(err)
	}
	if len(buf) == 0 {
		return nil
	}

	paramsInfo := new(ConsensusParamsInfo)
	err = rlp.DecodeBytes(buf, paramsInfo)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		panic(fmt.Sprintf(`LoadConsensusParams: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return paramsInfo
}

// Bytes serializes the ConsensusParamsInfo
func (params ConsensusParamsInfo) Bytes() []byte {
	b, err := rlp.EncodeToBytes(params)
	if err != nil {
		panic(err)
	}
	return b
}

// saveConsensusParamsInfo persists the consensus params for the next block to disk.
// It should be called from s.Save(), right before the state itself is persisted.
// If the consensus params did not change after processing the latest block,
// only the last height for which they changed is persisted.
func saveConsensusParamsInfo(db kaidb.Database, nextHeight, changeHeight uint64, params types.ConsensusParams) {
	paramsInfo := &ConsensusParamsInfo{
		LastHeightChanged: uint64(changeHeight),
	}
	if changeHeight == nextHeight {
		paramsInfo.ConsensusParams = params
	}
	db.Put(calcConsensusParamsKey(nextHeight), paramsInfo.Bytes())
}

// MakeGenesisState creates state from types.GenesisDoc.
func MakeGenesisState(genDoc *genesis.Genesis) (LastestBlockState, error) {

	var validatorSet, nextValidatorSet *types.ValidatorSet
	if genDoc.Validators == nil {
		validatorSet = types.NewValidatorSet(nil, 0, 1)
		nextValidatorSet = types.NewValidatorSet(nil, 0, 1)
	} else {
		validators := make([]*types.Validator, len(genDoc.Validators))
		for i, val := range genDoc.Validators {
			validators[i] = types.NewValidator(common.HexToAddress(val.Address), val.Power)
		}
		validatorSet = types.NewValidatorSet(validators, 0, 1)
		nextValidatorSet = types.NewValidatorSet(validators, 0, 1)
		nextValidatorSet.AdvanceProposer(1)
	}

	return LastestBlockState{
		LastBlockHeight: common.NewBigInt64(0),
		LastBlockID:     types.BlockID{},
		LastBlockTime:   genDoc.Timestamp,

		NextValidators:              nextValidatorSet,
		Validators:                  validatorSet,
		LastValidators:              types.NewValidatorSet(nil, 0, 1),
		LastHeightValidatorsChanged: common.NewBigInt64(0),

		//ConsensusParams:                  *genDoc.ConsensusParams,
		LastHeightConsensusParamsChanged: 1,
	}, nil
}