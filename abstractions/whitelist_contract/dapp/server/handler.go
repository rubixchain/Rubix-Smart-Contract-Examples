package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"dapp/storage"
	wasm "wasm_go"

	_ "github.com/joho/godotenv/autoload"
)

type ContractInputRequest struct {
	Port              string `json:"port"`
	SmartContractHash string `json:"smart_contract_hash"` //port should also be added here, so that the api can understand which node.
}


// Note: here is the output is of type []string is because of the design on the WASM contract.
// Hence, depending on the implementation on WASM contract, output type may vary.
func executeWASMContract(wasmDir string, states []SCTDataReply) ([]string, error) {
	wasmModule, err := wasm.NewWasmModule(wasmDir)
	if err != nil {
		return nil, fmt.Errorf("err occured while creating WASM object\n")
	}

	var wasmExecutionOutputs []string = make([]string, 0)
	for idx, state := range states {
		if state.BlockNo == 0 {
			continue
		}

		smartContractDataInput := state.SmartContractData

		fmt.Printf("Smart contract data being fed %s at idx %v", smartContractDataInput, idx)

		wasmOutput, err := wasmModule.CallFunction(smartContractDataInput)
		if err != nil {
			return nil, fmt.Errorf("err occured while calling for input %v, err: %v", smartContractDataInput, err)
		}

		wasmOutputStr, ok := wasmOutput.(string)
		if !ok {
			return nil, fmt.Errorf("unable to assert interface to string for WASM function output")
		}
		wasmExecutionOutputs = append(wasmExecutionOutputs, wasmOutputStr)
	}
	
	return wasmExecutionOutputs, nil
}

func runWhitelistContractHandle(w http.ResponseWriter, r *http.Request) {
	// Handle Input from Execute API
	var req ContractInputRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	latestContractData, err := GetLatestSmartContractData()
	fmt.Println("Latest Smart Contract data: ", latestContractData)
	if err != nil {
		fmt.Printf("err occured while getting latest smart contract data: %v\n", err)
		return
	}

	if len(latestContractData) != 1 {
		fmt.Printf("expected length of Contract Data (SCDataReply) slice to be 1, got %v\n", len(latestContractData))
		return
	}
	
	// Setup and call WASM contract function
	wasmFileDir := os.Getenv("WASM_FILE_DIR")

	outputs, err := executeWASMContract(wasmFileDir, latestContractData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// The output type in WASM is expected to be of String type
	fmt.Println("Outputs are : ", outputs)
	whitelistedDid := outputs[0]

	stateFileDir := os.Getenv("STATE_STORAGE_JSON")
	if err := storage.AddWhitelistedDIDToState(whitelistedDid, stateFileDir); err != nil {
		fmt.Printf("unable to create or update state file, err: %v", err)
	}
}

func getWhitelistedDidsHandler(w http.ResponseWriter, _ *http.Request) {
	stateFileDir := os.Getenv("STATE_STORAGE_JSON")

	stateData, err := storage.ReadState(stateFileDir)
	if err != nil {
		fmt.Printf("unable to read state, err: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, writeErr := w.Write(stateData)
	if writeErr != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}


func syncStateHandler(w http.ResponseWriter, _ *http.Request) {
	smartContractHash := os.Getenv("WHITELISTING_CONTRACT_HASH")

	states, err := GetAllSmartContractData()
	if err != nil {
		fmt.Printf("unable to fetch all states for contract %v, err: %v\n", smartContractHash, err)
		return
	}

	wasmFileDir := os.Getenv("WASM_FILE_DIR")

	whitelistedDIDs, err := executeWASMContract(wasmFileDir, states)
	if err != nil {
		fmt.Println(err)
		return
	}

	stateFileDir := os.Getenv("STATE_STORAGE_JSON")
	
	err = storage.SyncStateJSONFile(stateFileDir, whitelistedDIDs)
	if err != nil {
		fmt.Printf("unable to sync the state: %v", err)
		return
	}
}