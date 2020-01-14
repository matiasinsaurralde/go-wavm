package wavm

/*
#cgo CFLAGS: -I/usr/local/include/WAVM
#cgo LDFLAGS: -lWAVM
#include <stdio.h>
#include <stdlib.h>
#include <wavm-c/wavm-c.h>
typedef struct export_t
{
	const char* name;
	size_t num_name_bytes;
	wasm_externtype_t* typ;
} export_t;
*/
import "C"
import (
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

// Config holds the main configuration structure
type Config struct {
	WASIEnabled bool
}

// WASMEngine wraps wasm_engine_t
type WASMEngine struct {
	ptr *WASMEngineT
	cfg *Config

	compartment *WASMCompartment
	store *WASMStore
}

// Ptr returns a wasm_engine_t pointer
func (e *WASMEngine) Ptr() *C.wasm_engine_t {
	return (*C.wasm_engine_t)(e.ptr)
}

// Delete wraps wasm_engine_delete
func (e *WASMEngine) Delete() {
	C.wasm_engine_delete(e.Ptr())
}

// LoadModule loads a module from []byte, precompiled modules are supported too.
// Returns a wasm_module_instance_t pointer.
// TODO: see if wasm_engine_t needs to be used at some level...
func (e *WASMEngine) LoadModule(wasmBytes []byte, precompiled bool) *WASMModule {
	bytes := C.CBytes(wasmBytes)
	var ptr *C.wasm_module_t
	if precompiled {
		ptr = C.wasm_module_precomp_new((*C.char)(bytes), (C.ulong)(len(wasmBytes)))
	} else {
		ptr = C.wasm_module_std_new((*C.char)(bytes), (C.ulong)(len(wasmBytes)))
	}
	return &WASMModule{
		ptr: (*WASMModuleT)(ptr),
	}
}


// WASIRun calls the entrypoint function and returns an instance of the current module
func(e *WASMEngine) WASIRun(m *WASMModule) *WASMInstance {
	instancePtr := C.wasi_run(m.Ptr(), e.store.Ptr())
	return &WASMInstance{
		ptr: (*WASMInstanceT)(instancePtr),
	}
}

// NewInstance initializes a new module instance
// TODO: clarify whether to use different stores at instance level.
func(e *WASMEngine) NewInstance(module *WASMModule) *WASMInstance {
	ptr := C.wasm_instance_new(e.store.Ptr(), module.Ptr(), nil, nil, C.CString(""))
	return &WASMInstance{
		ptr: (*WASMInstanceT)(ptr),
		engine: e,
	}
}

// NewEngine wraps wasm_engine_new
func NewEngine(cfg *Config) *WASMEngine {
	enginePtr := C.wasm_engine_new()
	engine := &WASMEngine{
		ptr: (*WASMEngineT)(enginePtr),
		cfg: cfg,
	}
	engine.compartment = NewWASMCompartment(engine, "")
	engine.store = NewWASMStore(engine.compartment, "")
	return engine
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
	// defer C.free(unsafe.Pointer(cDebugName))
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
	// defer C.free(unsafe.Pointer(cDebugName))
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

// NewWASMModule wraps wasm_module_new
func NewWASMModule(engine *WASMEngine, wasmBytes []byte, precompiled bool) *WASMModule {
	bytes := C.CBytes(wasmBytes)
	var ptr *C.wasm_module_t
	if precompiled {
		ptr = C.wasm_module_std_new((*C.char)(bytes), (C.ulong)(len(wasmBytes)))
	} else {
		ptr = C.wasm_module_precomp_new((*C.char)(bytes), (C.ulong)(len(wasmBytes)))
	}
	return &WASMModule{
		ptr: (*WASMModuleT)(ptr),
	}
}

// Ptr returns a wasm_module_t pointer
func (m *WASMModule) Ptr() *C.wasm_module_t {
	return (*C.wasm_module_t)(m.ptr)
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

// WASMInstance wraps wasm_instance_t
type WASMInstance struct {
	ptr *WASMInstanceT
	engine *WASMEngine
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
		instance: i,
	}
}

// WASMExtern wraps wasm_extern_t
type WASMExtern struct {
	ptr *WASMExternT
	instance *WASMInstance
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
		extern: e,
	}
}

// WASMFunction wraps wasm_func_t
type WASMFunction struct {
	ptr *WASMFuncT
	extern *WASMExtern
}

// Ptr returns a wasm_func_t pointer
func (f *WASMFunction) Ptr() *C.wasm_func_t {
	return (*C.wasm_func_t)(f.ptr)
}

// Call wraps wasm_func_call
func (f *WASMFunction) Call(a int, b int) int {
	args := make([]C.wasm_val_t, 2)
	args[0] = C.wasm_val_t{byte(a)}
	args[1] = C.wasm_val_t{byte(b)}
	results := make([]C.wasm_val_t, 1)
	C.wasm_func_call(f.extern.instance.engine.store.Ptr(), f.Ptr(), &args[0], &results[0])
	return int(results[0][0])
}

// NewWASMInstance wraps wasm_instance_new
// TODO: handle imports
func NewWASMInstance(store *WASMStore, module *WASMModule, compartment *WASMCompartment) *WASMInstance {
	// imports := make([]*C.wasm_extern_t, 1)
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

// Precompile precompiles a WASM module, can be held in memory or written to a file. Useful for large modules
// TODO: allocate in a better way
func Precompile(filename string, input []byte) (data []byte) {
	cRawModule := C.CBytes(input)
	len := C.ulong(len(input))
	var out unsafe.Pointer = C.malloc(4000000) // ?
	defer C.free(out)
	outlen := C.wasm_precompile2((*C.char)(cRawModule), len, (*C.char)(out))
	data = C.GoBytes(out, outlen)
	return data
}