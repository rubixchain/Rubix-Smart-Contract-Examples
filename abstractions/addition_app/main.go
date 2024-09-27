// main.go

package main

import (
    "fmt"
    "log"

    "wasm_go"
)


func main() {
    // Initialize the WASM module
    wasmModule, err := wasm_go.NewWasmModule("./addition_contract.wasm")
    if err != nil {
        log.Fatalf("Failed to initialize WASM module: %v", err)
    }

    // Prepare input data as a JSON string
    args := `{"concatenate_strings": {"a": "Arnab", "b": "Ghose"}}`

    // Call the function
    result, err := wasmModule.CallFunction(args)
    if err != nil {
        log.Fatalf("Function call failed: %v", err)
    }

    // Type assertion based on expected output type
    sum, ok := result.(string)
    if !ok {
        log.Fatalf("Expected string result, got %T", result)
    }

    fmt.Printf("Result: %s\n", sum) // Expected output: "60"
}
