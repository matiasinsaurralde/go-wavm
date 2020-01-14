package main

import (
	"fmt"
	"io/ioutil"
	"log"

	wavm "github.com/matiasinsaurralde/go-wavm"
)

func main() {
	log.Print("Initializing WAVM")
	engine := wavm.NewEngine(&wavm.Config{
		WASIEnabled: true,
	})
	rawModule, _ := ioutil.ReadFile("quickjs.wasm.precompiled")
	module := engine.LoadModule(rawModule, true)
	numExports := module.NumExports()
	for i := 0; i < numExports; i++ {
		export := module.GetExport(i)
		fmt.Println(i, export)
	}
	instance := engine.WASIRun(module)
	log.Print("Instance loaded ", instance)
}
