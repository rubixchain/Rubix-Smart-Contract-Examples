package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ContractInputRequest struct {
	Port              string `json:"port"`
	SmartContractHash string `json:"smart_contract_hash"` //port should also be added here, so that the api can understand which node.
}

type RubixResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type SmartContractInput struct {
	Did            string `json:"did"`
	BinaryCodePath string `json:"binaryCodePath"`
	RawCodePath    string `json:"rawCodePath"`
	SchemaFilePath string `json:"schemaFilePath"`
	Port           string `json:"port"`
}

type DeploySmartContractInput struct {
	Comment            string `json:"comment"`
	DeployerAddress    string `json:"deployerAddress"`
	QuorumType         int    `json:"quorumType"`
	RbtAmount          int    `json:"rbtAmount"`
	SmartContractToken string `json:"smartContractToken"`
	Port               string `json:"port"`
}

type ExecuteSmartContractInput struct {
	Comment            string `json:"comment"`
	ExecutorAddress    string `json:"executorAddress"`
	QuorumType         int    `json:"quorumType"`
	SmartContractData  string `json:"smartContractData"`
	SmartContractToken string `json:"smartContractToken"`
	Port               string `json:"port"`
}

type GetSmartContractDataInput struct {
	Port  string `json:"port"`
	Token string `json:"token"`
}

func ContractInputHandler(w http.ResponseWriter, r *http.Request) {

	var req ContractInputRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	err3 := godotenv.Load()
	if err3 != nil {
		fmt.Println("Error loading .env file:", err3)
		return
	}
	port := req.Port
	folderPath := "voting_contract/target/wasm32-unknown-unknown/debug/voting_contract.wasm"
	stateFilePath := "store_state/vote_contract/votefile.json"

	fmt.Println(folderPath)
	_, err1 := os.Stat(folderPath)
	fmt.Println(err1)
	if os.IsNotExist(err1) {
		fmt.Println("Smart Contract not found")
		RunSmartContract(folderPath, stateFilePath, port, req.SmartContractHash)
	} else if err == nil {
		fmt.Printf("Folder '%s' exists", folderPath)

		RunSmartContract(folderPath, stateFilePath, port, req.SmartContractHash)

	} else {
		fmt.Printf("Error while checking folder: %v\n", err)
	}

	resp := RubixResponse{Status: true, Message: "Callback Successful", Result: "Success"}
	json.NewEncoder(w).Encode(resp)

}
