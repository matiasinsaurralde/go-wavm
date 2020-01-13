package main

import (
	"testing"

	wavm "github.com/matiasinsaurralde/go-wavm"
)

func BenchmarkWAVMSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		engine := wavm.NewWASMEngine()
		compartment := wavm.NewWASMCompartment(engine, "")
		store := wavm.NewWASMStore(compartment, "")
		module := wavm.NewWASMModule(engine, rawModule)

		instance := wavm.NewWASMInstance(store, module)

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
		_ = fn.Call(store, 2, 2)
		module.Delete()
		instance.Delete()
		store.Delete()
		compartment.Delete()
		engine.Delete()
	}
}

func BenchmarkWAVMSumReentrant(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = fn.Call(store, 2, 2)
	}
}

func TestWAVMSum(t *testing.T) {
	ret := fn.Call(store, 2, 2)
	if ret != 4 {
		t.Fatal("Unexpected output value")
	}
}
