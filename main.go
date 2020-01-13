package main

/*
#cgo LDFLAGS: -lWAVM
#include <stdio.h>
#include <stdlib.h>
#include <wavm-c.h>

typedef struct export_t
{
	const char* name;
	size_t num_name_bytes;
	wasm_externtype_t* typ;
} export_t;
*/
import "C"

import (
	"io/ioutil"
	"log"
	"unsafe"
)

// WASMExternKind represents wasm_externkind_enum values
type WASMExternKind int

const (
	// WASMExternFunc is an alias for WASM_EXTERN_FUNC
	WASMExternFunc WASMExternKind = iota
	// WASMExternGlobal is an alias for WASM_EXTERN_GLOBAL
	WASMExternGlobal
	// WASMExternTable is an alias for WASM_EXTERN_TABLE
	WASMExternTable
	// WASMExternMemory is an alias for WASM_EXTERN_MEMORY
	WASMExternMemory
)

// WASMEngineT is an alias for wasm_engine_t
type WASMEngineT C.wasm_engine_t
// WASMCompartmentT is an alias for wasm_compartment_t
type WASMCompartmentT C.wasm_compartment_t
// WASMStoreT is an alias for wasm_store_T
type WASMStoreT C.wasm_store_t
// WASMModuleT is an alias for wasm_module_t
type WASMModuleT C.wasm_module_t
// WASMInstanceT is an alias for wasm_instance_t
type WASMInstanceT C.wasm_instance_t
// WASMExternT is an alias for wasm_extern_t
type WASMExternT C.wasm_extern_t
// WASMTrapT is an alias for wasm_trap_t
type WASMTrapT C.wasm_trap_t
// WASMFuncT is an alias for wasm_trap_t
type WASMFuncT C.wasm_func_t
// WASMValT is an alias for wasm_val_t
type WASMValT C.wasm_val_t
// WASMExportT is an alias for wasm_export_t
type WASMExportT C.wasm_export_t

// WASMEngine wraps wasm_engine_t
type WASMEngine struct {
	ptr *WASMEngineT
}

// Ptr returns a wasm_engine_t pointer
func (e *WASMEngine) Ptr() *C.wasm_engine_t {
	return (*C.wasm_engine_t)(e.ptr)
}

// Delete wraps wasm_engine_delete
func (e *WASMEngine) Delete() {
	C.wasm_engine_delete(e.Ptr())
}

// NewWASMEngine wraps wasm_engine_new
func NewWASMEngine() *WASMEngine {
	engine := C.wasm_engine_new()
	return &WASMEngine{
		ptr: (*WASMEngineT)(engine),
	}
}

// WASMCompartment wraps wasm_compartment_t
type WASMCompartment struct {
	ptr *WASMCompartmentT
}

// Ptr returns a wasm_compartment_t pointer
func (c *WASMCompartment) Ptr() *C.wasm_compartment_t {
	return (*C.wasm_compartment_t)(c.ptr)
}

// Delete wraps wasm_compartment_delete calls
func (c *WASMCompartment) Delete() {
	C.wasm_compartment_delete(c.Ptr())
}

// NewWASMCompartment wraps wasm_compartment_new calls
func NewWASMCompartment(engine *WASMEngine, debugName string) *WASMCompartment {
	cDebugName := C.CString(debugName)
	defer C.free(unsafe.Pointer(cDebugName))
	compartment := C.wasm_compartment_new(engine.Ptr(), cDebugName)
	return &WASMCompartment{
		ptr: (*WASMCompartmentT)(compartment),
	}
}

// WASMStore wraps wasm_store_t
type WASMStore struct {
	ptr *WASMStoreT
}

// Ptr returns a wasm_store_t pointer
func (s *WASMStore) Ptr() *C.wasm_store_t {
	return (*C.wasm_store_t)(s.ptr)
}

// Delete wraps wasm_store_delete
func (s *WASMStore) Delete() {
	C.wasm_store_delete(s.Ptr())
}

// NewWASMStore wraps wasm_store_new calls
func NewWASMStore(compartment *WASMCompartment, debugName string) *WASMStore {
	cDebugName := C.CString(debugName)
	defer C.free(unsafe.Pointer(cDebugName))
	store := C.wasm_store_new(compartment.Ptr(), cDebugName)
	return &WASMStore{
		ptr: (*WASMStoreT)(store),
	}
}

// WASMModule wraps wasm_module_t and provides helpers for retrieving exports 
type WASMModule struct {
	ptr *WASMModuleT
}

// Delete wraps wasm_module_delete
func (m *WASMModule) Delete() {
	C.wasm_module_delete(m.Ptr())
}

// NumExports wraps wasm_module_num_exports
func (m *WASMModule) NumExports() int {
	return int(C.wasm_module_num_exports(m.Ptr()))
}

// GetExport builds a WASMExport data structure based on wasm_module_export output
// It also resolves the extern type
func (m *WASMModule) GetExport(index int) *WASMExport {
	var export C.wasm_export_t
	C.wasm_module_export(m.Ptr(), C.ulong(index), &export)
	name := C.GoString(export.name)
	// Ugly hack to access "type" (it's a reserved word, not allowed in cgo):
	exportT := (*C.export_t)(unsafe.Pointer(&export))
	kind := wasmExternTypeKind((*WASMExternT)(exportT.typ))
	return &WASMExport{
		Name: name,
		Kind: kind,
	}
}

// WASMExport wraps wasm_export_t values
type WASMExport struct {
	Name string
	Kind WASMExternKind
}

// IsFunction is a helper for WASMExport wrappers
func (e *WASMExport) IsFunction() bool {
	return e.Kind == WASMExternFunc
}

// NewWASMModule wraps wasm_module_new
func NewWASMModule(engine *WASMEngine, wasmBytes []byte) *WASMModule {
	bytes := C.CBytes(wasmBytes)
	ptr := C.wasm_module_new(engine.Ptr(), (*C.char)(bytes), (C.ulong)(len(wasmBytes)))
	return &WASMModule{
		ptr: (*WASMModuleT)(ptr),
	}
}

// Ptr returns a wasm_module_t pointer
func (m *WASMModule) Ptr() *C.wasm_module_t {
	return (*C.wasm_module_t)(m.ptr)
}

// WASMInstance wraps wasm_instance_t
type WASMInstance struct {
	ptr *WASMInstanceT
}

// Ptr returns a wasm_instance_t pointer
func (i *WASMInstance) Ptr() *C.wasm_instance_t {
	return (*C.wasm_instance_t)(i.ptr)
}

// Delete wraps wasm_instance_delete
func (i *WASMInstance) Delete() {
	C.wasm_instance_delete(i.Ptr())
}

// NumExports wraps wasm_instance_num_exports
func (i *WASMInstance) NumExports() int {
	return int(C.wasm_instance_num_exports(i.Ptr()))
}

// GetExport wraps wasm_instance_export
func (i *WASMInstance) GetExport(index int) *WASMExtern {
	ptr := C.wasm_instance_export(i.Ptr(), C.ulong(index))
	return &WASMExtern{
		ptr: (*WASMExternT)(ptr),
	}
}

// WASMExtern wraps wasm_extern_t
type WASMExtern struct {
	ptr *WASMExternT
}

// Ptr returns a wasm_extern_t pointer
func (e *WASMExtern) Ptr() *C.wasm_extern_t {
	return (*C.wasm_extern_t)(e.ptr)
}

// AsFunction wraps wasm_extern_as_func
func (e *WASMExtern) AsFunction() *WASMFunction {
	// TODO: handle errors in case the input isn't a function?
	ptr := C.wasm_extern_as_func(e.Ptr())
	return &WASMFunction{
		ptr: (*WASMFuncT)(ptr),
	}
}

// WASMFunction wraps wasm_func_t
type WASMFunction struct {
	ptr *WASMFuncT
}

// Ptr returns a wasm_func_t pointer
func (f *WASMFunction) Ptr() *C.wasm_func_t {
	return (*C.wasm_func_t)(f.ptr)
}

// Call wraps wasm_func_call
func (f *WASMFunction) Call(store *WASMStore, a int, b int) int {
	args := make([]C.wasm_val_t, 2)
	args[0] = C.wasm_val_t{byte(a)}
	args[1] = C.wasm_val_t{byte(b)}
	results := make([]C.wasm_val_t, 1)
	C.wasm_func_call(store.Ptr(), f.Ptr(), &args[0], &results[0])
	return int(results[0][0])
}

// NewWASMInstance wraps wasm_instance_new
func NewWASMInstance(store *WASMStore, module *WASMModule) *WASMInstance {
	ptr := C.wasm_instance_new(store.Ptr(), module.Ptr(), nil, nil, C.CString(""))
	return &WASMInstance{
		ptr: (*WASMInstanceT)(ptr),
	}
}

// wasmExternTypeKind resolves the extern type.
func wasmExternTypeKind(ptr *WASMExternT) WASMExternKind {
	kind := C.wasm_extern_kind((*C.wasm_extern_t)(ptr))
	return WASMExternKind(kind)
}

func main() {
	log.Print("Initializing WAVM")
	engine := NewWASMEngine()
	rawModule, err := ioutil.ReadFile("simple.wasm")
	if err != nil {
		panic(err)
	}
	compartment := NewWASMCompartment(engine, "")
	store := NewWASMStore(compartment, "")
	module := NewWASMModule(engine, rawModule)
	instance := NewWASMInstance(store, module)

	numExports := module.NumExports()
	log.Printf("Found %d module exports.", numExports)

	var fn *WASMFunction
	for i := 0; i < numExports; i++ {
		moduleExport := module.GetExport(i)
		if moduleExport.Name != "sum" {
			continue
		}
		instanceExport := instance.GetExport(i)
		fn = instanceExport.AsFunction()
	}

	if fn == nil {
		panic("sum function not found")
	}

	log.Print("Calling sum(2,2)")
	ret := fn.Call(store, 2, 2)
	log.Printf("Returns %d", ret)

	module.Delete()
	instance.Delete()
	store.Delete()
	compartment.Delete()
	engine.Delete()
}
