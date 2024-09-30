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

    // Prepare input data as a JSON string (input from Smart Contract State)
    smartContractData_Block_2 := `{"concatenate_strings": {"a": "Arnab", "b": "Ghose"}}`
    smartContractData_Block_3 := `{"add_three_nums": {"a": 1, "b": 2, "c": 3}}`

    // Call the function
    contractInput1Result, err := wasmModule.CallFunction(smartContractData_Block_2)
    if err != nil {
        log.Fatalf("Function call for concatenate_strings failed: %v", err)
    }

    contractInput2Result, err := wasmModule.CallFunction(smartContractData_Block_3)
    if err != nil {
        log.Fatalf("Function call for add_three_nums failed: %v", err)
    }

    concatResult, ok := contractInput1Result.(string)
    if !ok {
        log.Fatalf("Expected string result, got %T", contractInput1Result)
    }

    fmt.Printf("concatenate_strings Result: %s\n", concatResult) // Expected output: "60"

    sum, ok := contractInput2Result.(string)
    if !ok {
        log.Fatalf("Expected string result, got %T", contractInput2Result)
    }

    fmt.Printf("add_three_nums Result: %s\n", sum) // Expected output: "60"

    
}
