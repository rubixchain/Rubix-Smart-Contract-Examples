package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func checkIfDIDIsWhitelisted(did string) (bool, error) {
	resp, err := http.Get("http://localhost:8080/get-whitelisted-dids")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	var whiteListedDIDMaps map[string][]string
	err = json.Unmarshal(respBody, &whiteListedDIDMaps)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	whitelistedDIDs := whiteListedDIDMaps["whitelisted_dids"]

	for _, whitelistedDID := range whitelistedDIDs {
		if did == whitelistedDID {
			return true, nil
		}
	}

	return false, nil
}

func runBiddingContractHandle(w http.ResponseWriter, r *http.Request) {
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

	// Get the dids info from Whitelisting Contract and check if input DID
	// is part of the Whitelist

	// Extract DID from Smart Contract Data
	inputSmartContractDataMap, err := unmarshalSmartContractDataString(latestContractData[0].SmartContractData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// SmartContractData: `{"place_nid": {"did":"did1", "bid": 123}}`
	inputDID, ok := inputSmartContractDataMap["place_bid"]["did"]
	if !ok {
		fmt.Println("invalid key in input smart contract data")
		return
	}

	isDIDWhitelisted, err := checkIfDIDIsWhitelisted(inputDID)
	if err != nil {
		fmt.Printf("error occured while requesting from Whitelisting contract, err: %v\n", err)
		return
	}

	if !isDIDWhitelisted {
		fmt.Printf("The DID %v is not allowed to place bid\n", inputDID)
		return
	}

	outputs, err := executeWASMContract(wasmFileDir, latestContractData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// The output type in WASM is expected to be of String type
	fmt.Println("Outputs are : ", outputs)
	bidInfo, err := unmarshalWASMOutputStringToMap(outputs[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	biddingDID := bidInfo["did"]
	biddingAmount := bidInfo["bid"]

	stateFileDir := os.Getenv("STATE_STORAGE_JSON")
	if err := storage.AddBiddingInfoToState(biddingDID, biddingAmount, stateFileDir); err != nil {
		fmt.Printf("unable to create or update state file, err: %v", err)
	}
}

func getBidsHandler(w http.ResponseWriter, _ *http.Request) {
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
