// main.go

package main

// import (
//     "fmt"
//     "log"

//     "wasm_go"
// )


// func main() {
//     // Initialize the WASM module
//     wasmModule, err := wasm_go.NewWasmModule("../artifacts/whitelist_contract.wasm")
//     if err != nil {
//         log.Fatalf("Failed to initialize WASM module: %v", err)
//     }

//     // Prepare input data as a JSON string (input from Smart Contract State)
//     smartContractInput := `{"whitelist_did": {"did": "did123"}}`

//     // Call the function
//     contractInputResult, err := wasmModule.CallFunction(smartContractInput)
//     if err != nil {
//         log.Fatalf("Function call for concatenate_strings failed: %v", err)
//     }

//     concatResult, ok := contractInputResult.(string)
//     if !ok {
//         log.Fatalf("Expected string result, got %T", contractInputResult)
//     }

//     fmt.Printf("whitelist_did Result: %s\n", concatResult)
// }

import (
	"dapp/server"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Server has been started")
	go server.RunServer()

	time.Sleep(300000 * time.Second)
	fmt.Println("Server has stopped!")
}
