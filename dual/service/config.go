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

package service

import (
	"github.com/kardiachain/go-kardia/dual/blockchain"
)

type DualConfig struct {
	// Protocol options
	NetworkId uint64 // Network

	// The genesis block of dual blockchain, which is inserted if the database is empty.
	// If nil, the Dual main net block is used.
	DualGenesis *blockchain.DualGenesis `toml:",omitempty"`

	// Dual's event pool options
	DualEventPool blockchain.EventPoolConfig

	// chaindata
	ChainData string

	// DB caches
	DbCaches int

	// DB handles
	DbHandles int
}
