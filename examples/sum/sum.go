package main

import (
	"io/ioutil"
	"log"

	wavm "github.com/matiasinsaurralde/go-wavm"
)

const (
	wasmFilename = "sum.wasm"
	fnName       = "sum"
)

var (
	engine      *wavm.WASMEngine
	compartment *wavm.WASMCompartment
	store       *wavm.WASMStore
	module      *wavm.WASMModule
	instance    *wavm.WASMInstance

	fn *wavm.WASMFunction

	rawModule []byte
)

func init() {
	log.Print("Initializing WAVM")
	engine = wavm.NewWASMEngine()
	var err error
	rawModule, err = ioutil.ReadFile(wasmFilename)
	if err != nil {
		panic(err)
	}
	compartment = wavm.NewWASMCompartment(engine, "")
	store = wavm.NewWASMStore(compartment, "")
	module = wavm.NewWASMModule(engine, rawModule)
	instance = wavm.NewWASMInstance(store, module)

	numExports := module.NumExports()
	log.Printf("Found %d module exports.", numExports)

	for i := 0; i < numExports; i++ {
		moduleExport := module.GetExport(i)
		if moduleExport.Name != fnName {
			continue
		}
		instanceExport := instance.GetExport(i)
		fn = instanceExport.AsFunction()
	}

	if fn == nil {
		panic("sum function not found")
	}
}

func main() {
	log.Print("Calling sum(2,2)")
	ret := fn.Call(store, 2, 2)
	log.Printf("Returns %d", ret)

	module.Delete()
	instance.Delete()
	store.Delete()
	compartment.Delete()
	engine.Delete()
}
