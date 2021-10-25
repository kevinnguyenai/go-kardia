/*
 *  Copyright 2021 KardiaChain
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

package tracers

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/kardiachain/go-kardia/configs"
	"github.com/kardiachain/go-kardia/kai/state"
	"github.com/kardiachain/go-kardia/kvm"
	"github.com/kardiachain/go-kardia/lib/common"
)

type account struct{}

func (account) SubBalance(amount *big.Int)                          {}
func (account) AddBalance(amount *big.Int)                          {}
func (account) SetAddress(common.Address)                           {}
func (account) Value() *big.Int                                     { return nil }
func (account) SetBalance(*big.Int)                                 {}
func (account) SetNonce(uint64)                                     {}
func (account) Balance() *big.Int                                   { return nil }
func (account) Address() common.Address                             { return common.Address{} }
func (account) SetCode(common.Hash, []byte)                         {}
func (account) ForEachStorage(cb func(key, value common.Hash) bool) {}

type dummyStatedb struct {
	state.StateDB
}

func (*dummyStatedb) GetRefund() uint64                       { return 1337 }
func (*dummyStatedb) GetBalance(addr common.Address) *big.Int { return new(big.Int) }

type vmContext struct {
	blockCtx kvm.BlockContext
	txCtx    kvm.TxContext
}

func runTrace(tracer *Tracer, vmctx *vmContext, chaincfg *configs.ChainConfig) (json.RawMessage, error) {
	env := kvm.NewKVM(vmctx.blockCtx, vmctx.txCtx, &dummyStatedb{}, chaincfg, kvm.Config{Debug: true, Tracer: tracer})
	var (
		startGas uint64 = 10000
		value           = big.NewInt(0)
	)
	contract := kvm.NewContract(account{}, account{}, value, startGas)
	contract.Code = []byte{byte(kvm.PUSH1), 0x1, byte(kvm.PUSH1), 0x1, 0x0}

	tracer.CaptureStart(env, contract.Caller(), contract.Address(), false, []byte{}, startGas, value)
	ret, err := env.Interpreter().Run(contract, []byte{}, false)
	tracer.CaptureEnd(ret, startGas-contract.Gas, 1, err)
	if err != nil {
		return nil, err
	}
	return tracer.GetResult()
}

func TestTracer(t *testing.T) {
	execTracer := func(code string) ([]byte, string) {
		t.Helper()
		tracer, err := New(code, new(Context))
		if err != nil {
			t.Fatal(err)
		}
		ret, err := runTrace(tracer, &vmContext{
			blockCtx: kvm.BlockContext{BlockHeight: big.NewInt(1)},
			txCtx:    kvm.TxContext{GasPrice: big.NewInt(100000)},
		}, configs.TestChainConfig)
		if err != nil {
			return nil, err.Error() // Stringify to allow comparison without nil checks
		}
		return ret, ""
	}
	for i, tt := range []struct {
		code string
		want string
		fail string
	}{
		{ // tests that we don't panic on bad arguments to memory access
			code: "{depths: [], step: function(log) { this.depths.push(log.memory.slice(-1,-2)); }, fault: function() {}, result: function() { return this.depths; }}",
			want: `[{},{},{}]`,
		}, { // tests that we don't panic on bad arguments to stack peeks
			code: "{depths: [], step: function(log) { this.depths.push(log.stack.peek(-1)); }, fault: function() {}, result: function() { return this.depths; }}",
			want: `["0","0","0"]`,
		}, { //  tests that we don't panic on bad arguments to memory getUint
			code: "{ depths: [], step: function(log, db) { this.depths.push(log.memory.getUint(-64));}, fault: function() {}, result: function() { return this.depths; }}",
			want: `["0","0","0"]`,
		}, { // tests some general counting
			code: "{count: 0, step: function() { this.count += 1; }, fault: function() {}, result: function() { return this.count; }}",
			want: `3`,
		}, { // tests that depth is reported correctly
			code: "{depths: [], step: function(log) { this.depths.push(log.stack.length()); }, fault: function() {}, result: function() { return this.depths; }}",
			want: `[0,1,2]`,
		}, { // tests to-string of opcodes
			code: "{opcodes: [], step: function(log) { this.opcodes.push(log.op.toString()); }, fault: function() {}, result: function() { return this.opcodes; }}",
			want: `["PUSH1","PUSH1","STOP"]`,
		}, { // tests intrinsic gas
			code: "{depths: [], step: function() {}, fault: function() {}, result: function(ctx) { return ctx.gasPrice+'.'+ctx.gasUsed+'.'+ctx.intrinsicGas; }}",
			want: `"100000.6.29000"`,
		}, { // tests too deep object / serialization crash
			code: "{step: function() {}, fault: function() {}, result: function() { var o={}; var x=o; for (var i=0; i<1000; i++){	o.foo={}; o=o.foo; } return x; }}",
			fail: "RangeError: json encode recursion limit    in server-side tracer function 'result'",
		},
	} {
		if have, err := execTracer(tt.code); tt.want != string(have) || tt.fail != err {
			t.Errorf("testcase %d: expected return value to be '%s' got '%s', error to be '%s' got '%s'\n\tcode: %v", i, tt.want, string(have), tt.fail, err, tt.code)
		}
	}
}

func TestHalt(t *testing.T) {
	t.Skip("duktape doesn't support abortion")

	timeout := errors.New("stahp")
	tracer, err := New("{step: function() { while(1); }, result: function() { return null; }}", new(Context))
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(1 * time.Second)
		tracer.Stop(timeout)
	}()
	vmTest := &vmContext{blockCtx: kvm.BlockContext{BlockHeight: big.NewInt(1)}, txCtx: kvm.TxContext{GasPrice: big.NewInt(100000)}}
	if _, err = runTrace(tracer, vmTest, configs.TestChainConfig); err.Error() != "stahp    in server-side tracer function 'step'" {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

func TestHaltBetweenSteps(t *testing.T) {
	tracer, err := New("{step: function() {}, fault: function() {}, result: function() { return null; }}", new(Context))
	if err != nil {
		t.Fatal(err)
	}
	env := kvm.NewKVM(kvm.BlockContext{BlockHeight: big.NewInt(1)}, kvm.TxContext{}, &dummyStatedb{}, configs.TestChainConfig, kvm.Config{Debug: true, Tracer: tracer})
	scope := &kvm.ScopeContext{
		Contract: kvm.NewContract(&account{}, &account{}, big.NewInt(0), 0),
	}
	tracer.CaptureState(env, 0, 0, 0, 0, scope, nil, 0, nil)
	timeout := errors.New("stahp")
	tracer.Stop(timeout)
	tracer.CaptureState(env, 0, 0, 0, 0, scope, nil, 0, nil)

	if _, err := tracer.GetResult(); err.Error() != timeout.Error() {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

// TestNoStepExec tests a regular value transfer (no exec), and accessing the statedb
// in 'result'
func TestNoStepExec(t *testing.T) {
	runEmptyTrace := func(tracer *Tracer, vmctx *vmContext) (json.RawMessage, error) {
		env := kvm.NewKVM(vmctx.blockCtx, vmctx.txCtx, &dummyStatedb{}, configs.TestChainConfig, kvm.Config{Debug: true, Tracer: tracer})
		startGas := uint64(10000)
		contract := kvm.NewContract(account{}, account{}, big.NewInt(0), startGas)
		tracer.CaptureStart(env, contract.Caller(), contract.Address(), false, []byte{}, startGas, big.NewInt(0))
		tracer.CaptureEnd(nil, startGas-contract.Gas, 1, nil)
		return tracer.GetResult()
	}
	execTracer := func(code string) []byte {
		t.Helper()
		tracer, err := New(code, new(Context))
		if err != nil {
			t.Fatal(err)
		}
		ret, err := runEmptyTrace(tracer, &vmContext{
			blockCtx: kvm.BlockContext{BlockHeight: big.NewInt(1)},
			txCtx:    kvm.TxContext{GasPrice: big.NewInt(100000)},
		})
		if err != nil {
			t.Fatal(err)
		}
		return ret
	}
	for i, tt := range []struct {
		code string
		want string
	}{
		{ // tests that we don't panic on accessing the db methods
			code: "{depths: [], step: function() {}, fault: function() {},  result: function(ctx, db){ return db.getBalance(ctx.to)} }",
			want: `"0"`,
		},
	} {
		if have := execTracer(tt.code); tt.want != string(have) {
			t.Errorf("testcase %d: expected return value to be %s got %s\n\tcode: %v", i, tt.want, string(have), tt.code)
		}
	}
}

func TestIsPrecompile(t *testing.T) {
	chaincfg := &configs.ChainConfig{
		Kaicon:  configs.MainnetChainConfig.Kaicon,
		ChainID: big.NewInt(0),
	}
	txCtx := kvm.TxContext{GasPrice: big.NewInt(100000)}
	tracer, err := New("{addr: toAddress('0000000000000000000000000000000000000009'), res: null, step: function() { this.res = isPrecompiled(this.addr); }, fault: function() {}, result: function() { return this.res; }}", new(Context))
	if err != nil {
		t.Fatal(err)
	}

	blockCtx := kvm.BlockContext{BlockHeight: big.NewInt(1)}
	res, err := runTrace(tracer, &vmContext{blockCtx, txCtx}, chaincfg)
	if err != nil {
		t.Error(err)
	}
	if string(res) != "false" {
		t.Errorf("Tracer should not consider as precompile")
	}
}

func TestEnterExit(t *testing.T) {
	// test that either both or none of enter() and exit() are defined
	if _, err := New("{step: function() {}, fault: function() {}, result: function() { return null; }, enter: function() {}}", new(Context)); err == nil {
		t.Fatal("tracer creation should've failed without exit() definition")
	}
	if _, err := New("{step: function() {}, fault: function() {}, result: function() { return null; }, enter: function() {}, exit: function() {}}", new(Context)); err != nil {
		t.Fatal(err)
	}

	// test that the enter and exit method are correctly invoked and the values passed
	tracer, err := New("{enters: 0, exits: 0, enterGas: 0, gasUsed: 0, step: function() {}, fault: function() {}, result: function() { return {enters: this.enters, exits: this.exits, enterGas: this.enterGas, gasUsed: this.gasUsed} }, enter: function(frame) { this.enters++; this.enterGas = frame.getGas(); }, exit: function(res) { this.exits++; this.gasUsed = res.getGasUsed(); }}", new(Context))
	if err != nil {
		t.Fatal(err)
	}

	scope := &kvm.ScopeContext{
		Contract: kvm.NewContract(&account{}, &account{}, big.NewInt(0), 0),
	}

	tracer.CaptureEnter(kvm.CALL, scope.Contract.Caller(), scope.Contract.Address(), []byte{}, 1000, new(big.Int))
	tracer.CaptureExit([]byte{}, 400, nil)

	have, err := tracer.GetResult()
	if err != nil {
		t.Fatal(err)
	}
	want := `{"enters":1,"exits":1,"enterGas":1000,"gasUsed":400}`
	if string(have) != want {
		t.Errorf("Number of invocations of enter() and exit() is wrong. Have %s, want %s\n", have, want)
	}
}
