package main

import (
	"io/ioutil"
	"testing"
)

const (
	fnName = "sum"
)

var (
	globalEngine      *WASMEngine
	globalFn          *WASMFunction
	globalCompartment *WASMCompartment
	globalStore       *WASMStore
	globalModule      *WASMModule
	globalInstance    *WASMInstance

	wasmBytes []byte
)

func init() {
	globalEngine = NewWASMEngine()
	wasmBytes, _ = ioutil.ReadFile("simple.wasm")
	globalCompartment = NewWASMCompartment(globalEngine, "")
	globalStore = NewWASMStore(globalCompartment, "")
	globalModule = NewWASMModule(globalEngine, wasmBytes)

	globalInstance = NewWASMInstance(globalStore, globalModule)

	numExports := globalModule.NumExports()

	for i := 0; i < numExports; i++ {
		moduleExport := globalModule.GetExport(i)
		if moduleExport.Name != fnName {
			continue
		}
		instanceExport := globalInstance.GetExport(i)
		globalFn = instanceExport.AsFunction()
		break
	}
}
func BenchmarkWAVMSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		engine := NewWASMEngine()
		compartment := NewWASMCompartment(engine, "")
		store := NewWASMStore(compartment, "")
		module := NewWASMModule(engine, wasmBytes)

		instance := NewWASMInstance(store, module)

		numExports := module.NumExports()

		var fn *WASMFunction
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
		_ = globalFn.Call(globalStore, 2, 2)
	}
}
