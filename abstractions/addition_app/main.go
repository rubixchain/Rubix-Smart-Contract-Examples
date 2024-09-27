package main

import (
    "fmt"
    "log"

    "wasm_go"
)

type AddThreeNumsReq struct {
    A uint32 `json:"a"`
    B uint32 `json:"b"`
    C uint32 `json:"c"`
}

func main() {
    // Initialize the WASM module
    wasmModule, err := wasm_go.NewWasmModule("./addition_contract.wasm")
    if err != nil {
        log.Fatalf("Failed to initialize WASM module: %v", err)
    }

    // Prepare input data
    input := AddThreeNumsReq{
        A: 10,
        B: 20,
        C: 30,
    }

    // Prepare a variable to receive the output
    var result string

    // Call the function
    err = wasmModule.CallFunction("add_three_nums", input, &result)
    if err != nil {
        log.Fatalf("Function call failed: %v", err)
    }

    fmt.Printf("Result: %s\n", result) // Expected output: "60"
}