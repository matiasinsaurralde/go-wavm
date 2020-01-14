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
	engine            *wavm.WASMEngine
	module            *wavm.WASMModule
	instance          *wavm.WASMInstance

	fn *wavm.WASMFunction

	rawModule []byte
)

func init() {
	log.Print("Initializing WAVM")
	engine = wavm.NewEngine(&wavm.Config{})
	var err error
	rawModule, err = ioutil.ReadFile("sum.wasm")
	if err != nil {
		panic(err)
	}
	module = engine.LoadModule(rawModule, false)
	instance = engine.NewInstance(module)

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
	ret := fn.Call(2, 2)
	log.Printf("Returns %d", ret)

	module.Delete()
	instance.Delete()
	engine.Delete()
}
