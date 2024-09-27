// wasm_module.go

package wasm_go

import (
    "encoding/binary"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"

    "github.com/bytecodealliance/wasmtime-go"
)

// WasmModule encapsulates the WASM module and its associated functions.
type WasmModule struct {
    engine      *wasmtime.Engine
    store       *wasmtime.Store
    instance    *wasmtime.Instance
    memory      *wasmtime.Memory
    allocFunc   *wasmtime.Func
    deallocFunc *wasmtime.Func
}

// NewWasmModule initializes and returns a new WasmModule.
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

// allocate allocates memory in WASM and copies the data.
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

// deallocate frees memory in WASM.
func (w *WasmModule) deallocate(ptr int32, size int32) error {
    _, err := w.deallocFunc.Call(w.store, ptr, size)
    return err
}

// CallFunction accepts a JSON string in the format:
// {"function_name": { ... input struct ... }}
// It invokes the corresponding WASM function and returns the output as interface{}.
func (w *WasmModule) CallFunction(args string) (interface{}, error) {
    // Parse the JSON string
    var inputMap map[string]interface{}
    err := json.Unmarshal([]byte(args), &inputMap)
    if err != nil {
        return nil, fmt.Errorf("failed to parse input JSON: %v", err)
    }

    if len(inputMap) != 1 {
        return nil, errors.New("input JSON must contain exactly one function")
    }

    // Extract function name and input struct
    var funcName string
    var inputStruct interface{}
    for key, value := range inputMap {
        funcName = key
        inputStruct = value
    }

    // Append '_wrapper' suffix to get the actual function name
    wrapperFuncName := funcName + "_wrapper"

    // Serialize the input struct to JSON
    inputJSON, err := json.Marshal(inputStruct)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize input struct: %v", err)
    }

    // Allocate memory for input data
    inputPtr, err := w.allocate(inputJSON)
    if err != nil {
        return nil, fmt.Errorf("failed to allocate memory for input data: %v", err)
    }
    defer w.deallocate(inputPtr, int32(len(inputJSON)))

    // Prepare pointers for output data
    outputPtrPtr, err := w.allocate(make([]byte, 4)) // 4 bytes for pointer
    if err != nil {
        return nil, fmt.Errorf("failed to allocate memory for output_ptr_ptr: %v", err)
    }
    defer w.deallocate(outputPtrPtr, 4)

    outputLenPtr, err := w.allocate(make([]byte, 8)) // 8 bytes for length
    if err != nil {
        return nil, fmt.Errorf("failed to allocate memory for output_len_ptr: %v", err)
    }
    defer w.deallocate(outputLenPtr, 8)

    // Retrieve the wrapper function
    function := w.instance.GetExport(w.store, wrapperFuncName).Func()
    if function == nil {
        return nil, fmt.Errorf("function %s does not exist in the contract", funcName)
    }

    // Call the wrapper function
    ret, err := function.Call(w.store, inputPtr, len(inputJSON), outputPtrPtr, outputLenPtr)
    if err != nil {
        return nil, fmt.Errorf("error calling WASM function: %v", err)
    }

    // Check return code
    retCode, ok := ret.(int32)
    if !ok {
        return nil, errors.New("unexpected return type from WASM function")
    }
    if retCode != 0 {
        return nil, errors.New("WASM function returned an error")
    }

    // Read output_ptr_ptr and output_len_ptr
    memoryData := w.memory.UnsafeData(w.store)
    if len(memoryData) < int(outputPtrPtr)+4 || len(memoryData) < int(outputLenPtr)+8 {
        return nil, errors.New("invalid memory access for output pointers")
    }

    outputPtr := int32(binary.LittleEndian.Uint32(memoryData[outputPtrPtr:]))
    outputLen := int32(binary.LittleEndian.Uint64(memoryData[outputLenPtr:]))

    // Validate memory bounds
    if outputPtr < 0 || outputPtr+outputLen > int32(len(memoryData)) {
        return nil, errors.New("output data exceeds memory bounds")
    }

    // Read output data
    outputData := make([]byte, outputLen)
    copy(outputData, memoryData[outputPtr:outputPtr+outputLen])

    // Deserialize output data
    var output interface{}
    err = json.Unmarshal(outputData, &output)
    if err != nil {
        return nil, fmt.Errorf("failed to deserialize output data: %v", err)
    }

    // Deallocate output data
    err = w.deallocate(outputPtr, outputLen)
    if err != nil {
        return nil, fmt.Errorf("failed to deallocate output data: %v", err)
    }

    return output, nil
}
