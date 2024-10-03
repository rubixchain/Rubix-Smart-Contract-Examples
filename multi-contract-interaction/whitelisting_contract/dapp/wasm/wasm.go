package wasm

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
	Engine      *wasmtime.Engine
	Store       *wasmtime.Store
	Instance    *wasmtime.Instance
	Memory      *wasmtime.Memory
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
		Engine:      engine,
		Store:       store,
		Instance:    instance,
		Memory:      memory,
		allocFunc:   allocFunc,
		deallocFunc: deallocFunc,
	}, nil
}

func (w *WasmModule) Allocate(data []byte) (int32, error) {
	size := len(data)
	result, err := w.allocFunc.Call(w.Store, size)
	if err != nil {
		return 0, err
	}
	ptr := result.(int32)
	memoryData := w.Memory.UnsafeData(w.Store)
	copy(memoryData[ptr:ptr+int32(size)], data)
	return ptr, nil
}

func (w *WasmModule) Deallocate(ptr int32, size int32) error {
	_, err := w.deallocFunc.Call(w.Store, ptr, size)
	return err
}

func (w *WasmModule) CallWasmFunc(wasmFuncName string, contractInput interface{}) (interface{}, error) {
	marshalledSmartContractDataInput, err := json.Marshal(contractInput)
	if err != nil {
		return nil, err
	}

	contractInputPtr, err := w.Allocate(marshalledSmartContractDataInput)
	if err != nil {
		return nil, err
	}
	defer w.Deallocate(contractInputPtr, int32(len(marshalledSmartContractDataInput)))

	contractOutputPtr, err := w.Allocate(make([]byte, 4))
	if err != nil {
		return nil, err
	}
	defer w.Deallocate(contractOutputPtr, 4)

	contractOutputLenPtr, err := w.Allocate(make([]byte, 8))
	if err != nil {
		return nil, err
	}
	defer w.Deallocate(contractOutputLenPtr, 8)

	wasmFuncExtern := w.Instance.GetExport(w.Store, wasmFuncName)
	if wasmFuncExtern == nil {
		return nil, fmt.Errorf("no function named %v exists in the WASM file", wasmFuncName)
	}

	fn := wasmFuncExtern.Func()
	if fn == nil {
		return nil, fmt.Errorf("unable to define the function object")
	}

	_, fnErr := fn.Call(
		w.Store,
		contractInputPtr,
		len(marshalledSmartContractDataInput),
		contractOutputPtr,
		contractOutputLenPtr,
	)
	if fnErr != nil {
		return nil, fmt.Errorf("failed while calling the function: %v, err: %v", wasmFuncName, fnErr)
	}

	memoryData := w.Memory.UnsafeData(w.Store)
	if len(memoryData) < int(contractOutputPtr)+4 || len(memoryData) < int(contractOutputLenPtr)+8 {
		return nil, fmt.Errorf("invalid memory access for output pointers")
	}

	outputPtr := int32(binary.LittleEndian.Uint32(memoryData[contractOutputPtr:]))
	outputLen := int32(binary.LittleEndian.Uint64(memoryData[contractOutputLenPtr:]))

	// Validate memory bounds
	if outputPtr < 0 || outputPtr+outputLen > int32(len(memoryData)) {
		return nil, fmt.Errorf("output data exceeds memory bounds")
	}

	outputData := make([]byte, outputLen)
	copy(outputData, memoryData[outputPtr:outputPtr+outputLen])

	var output interface{}
	err = json.Unmarshal(outputData, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialise output data, err: %v", err)
	}

	err = w.Deallocate(outputPtr, outputLen)
	if err != nil {
		return nil, fmt.Errorf("failed to deallocate output data, err: %v", err)
	}

	return output, nil
}
