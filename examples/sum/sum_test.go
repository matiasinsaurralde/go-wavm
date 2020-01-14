package main

import (
	"testing"

	wavm "github.com/matiasinsaurralde/go-wavm"
)

var (
	defaultFn *wavm.WASMFunction
)

func init() {
	engine := wavm.NewEngine(&wavm.Config{})
	module := engine.LoadModule(rawModule, false)
	instance := engine.NewInstance(module)

	numExports := instance.NumExports()

	for i := 0; i < numExports; i++ {
		moduleExport := module.GetExport(i)
		if moduleExport.Name != fnName {
			continue
		}
		instanceExport := instance.GetExport(i)
		defaultFn = instanceExport.AsFunction()
		break
	}
}

func initAndCallFn() {
	engine := wavm.NewEngine(&wavm.Config{})
	module := engine.LoadModule(rawModule, false)
	instance := engine.NewInstance(module)

	numExports := module.NumExports()

	var fn *wavm.WASMFunction
	for i := 0; i < numExports; i++ {
		moduleExport := module.GetExport(i)
		if moduleExport.Name != fnName {
			continue
		}
		instanceExport := instance.GetExport(i)
		fn = instanceExport.AsFunction()
		break
	}
	_ = fn.Call(2, 2)
	module.Delete()
	instance.Delete()
	engine.Delete()
}

func BenchmarkWAVMSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		initAndCallFn()
	}
}

func BenchmarkWAVMSumReentrant(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = fn.Call(2, 2)
	}
}

func TestWAVMSum(t *testing.T) {
	ret := fn.Call(2, 2)
	if ret != 4 {
		t.Fatal("Unexpected output value")
	}
}
