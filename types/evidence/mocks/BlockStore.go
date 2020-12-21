// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	types "github.com/kardiachain/go-kardia/types"
	mock "github.com/stretchr/testify/mock"
)

// BlockStore is an autogenerated mock type for the BlockStore type
type BlockStore struct {
	mock.Mock
}

// LoadBlockCommit provides a mock function with given fields: height
func (_m *BlockStore) LoadBlockCommit(height uint64) *types.Commit {
	ret := _m.Called(height)

	var r0 *types.Commit
	if rf, ok := ret.Get(0).(func(uint64) *types.Commit); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Commit)
		}
	}

	return r0
}

// LoadBlockMeta provides a mock function with given fields: height
func (_m *BlockStore) LoadBlockMeta(height uint64) *types.BlockMeta {
	ret := _m.Called(height)

	var r0 *types.BlockMeta
	if rf, ok := ret.Get(0).(func(uint64) *types.BlockMeta); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BlockMeta)
		}
	}

	return r0
}
