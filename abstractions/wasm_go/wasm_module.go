package wasm_go

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/bytecodealliance/wasmtime-go"
)

type WasmModule struct {
    engine      *wasmtime.Engine
    store       *wasmtime.Store
    instance    *wasmtime.Instance
    memory      *wasmtime.Memory
    allocFunc   *wasmtime.Func
    deallocFunc *wasmtime.Func
}

func NewWasmModule(wasmPath string) (*WasmModule, error) {
    // Read the WASM file
    wasmBytes, err := ioutil.ReadFile(wasmPath)
    if err != nil {
        return nil, err
    }

    engine := wasmtime.NewEngine()
    module, err := wasmtime.NewModule(engine, wasmBytes)
    if err != nil {
        return nil, err
    }

    store := wasmtime.NewStore(engine)
    linker := wasmtime.NewLinker(engine)

    instance, err := linker.Instantiate(store, module)
    if err != nil {
        return nil, err
    }

    memory := instance.GetExport(store, "memory").Memory()
    if memory == nil {
        return nil, errors.New("failed to find memory export")
    }

    allocFunc := instance.GetExport(store, "alloc").Func()
    if allocFunc == nil {
        return nil, errors.New("failed to find alloc function")
    }

    deallocFunc := instance.GetExport(store, "dealloc").Func()
    if deallocFunc == nil {
        return nil, errors.New("failed to find dealloc function")
    }

    return &WasmModule{
        engine:      engine,
        store:       store,
        instance:    instance,
        memory:      memory,
        allocFunc:   allocFunc,
        deallocFunc: deallocFunc,
    }, nil
}

func (w *WasmModule) allocate(data []byte) (int32, error) {
    size := len(data)
    result, err := w.allocFunc.Call(w.store, size)
    if err != nil {
        return 0, err
    }
    ptr := result.(int32)
    memoryData := w.memory.UnsafeData(w.store)
    copy(memoryData[ptr:ptr+int32(size)], data)
    return ptr, nil
}

func (w *WasmModule) deallocate(ptr int32, size int32) error {
    _, err := w.deallocFunc.Call(w.store, ptr, size)
    return err
}

func (w *WasmModule) CallFunction(funcName string, input interface{}, output interface{}) error {
    // Serialize input data
    inputData, err := json.Marshal(input)
    if err != nil {
        return err
    }

    // Allocate input data in WASM memory
    inputPtr, err := w.allocate(inputData)
    if err != nil {
        return err
    }
    defer w.deallocate(inputPtr, int32(len(inputData)))

    suffixedFuncName := funcName + "_wrapper"
    
    // Prepare pointers for output data
    outputPtrPtr, err := w.allocate(make([]byte, 4))
    if err != nil {
        return err
    }
    defer w.deallocate(outputPtrPtr, 4)

    outputLenPtr, err := w.allocate(make([]byte, 8))
    if err != nil {
        return err
    }
    defer w.deallocate(outputLenPtr, 8)

    // Retrieve the function
    functionExternObj := w.instance.GetExport(w.store, suffixedFuncName)
    if functionExternObj == nil {
        return fmt.Errorf("function %v does not exists in the contract", funcName)
    }

    function := w.instance.GetExport(w.store, suffixedFuncName).Func()
    if function == nil {
        return errors.New("failed to find function export")
    }

    // Call the function
    ret, err := function.Call(w.store, inputPtr, len(inputData), outputPtrPtr, outputLenPtr)
    if err != nil {
        return err
    }

    // Check return code
    retCode := ret.(int32)
    if retCode != 0 {
        return errors.New("function returned an error")
    }

    // Read output pointers
    memoryData := w.memory.UnsafeData(w.store)
    outputPtr := int32(binary.LittleEndian.Uint32(memoryData[outputPtrPtr:]))
    outputLen := int32(binary.LittleEndian.Uint64(memoryData[outputLenPtr:]))

    // Read output data
    outputData := make([]byte, outputLen)
    copy(outputData, memoryData[outputPtr:outputPtr+outputLen])

    // Deserialize output data
    err = json.Unmarshal(outputData, output)
    if err != nil {
        return err
    }

    // Deallocate output data
    err = w.deallocate(outputPtr, outputLen)
    if err != nil {
        return err
    }

    return nil
}