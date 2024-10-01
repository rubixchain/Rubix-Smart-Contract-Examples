package main

import (
	"bidding-contract/contract"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	contractModule "bidding-contract/contract"
)

type SmartContractDataReply struct {
	BasicResponse
	SCTDataReply []SCTDataReply
}

type BasicResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type SCTDataReply struct {
	BlockNo           uint64
	BlockId           string
	SmartContractData string
}

func GenerateSmartContract() {
	/*
		This did, wasmPath, schemaPath, rawcodePath and Port should be replaced according to your Rubix node configuration and
		the respective paths
	*/
	did := "bafybmibpgv4fe4xr7wwolrymxfphe7o45r4mynnzam6ohqqzvh3usmue2e"
	wasmPath := "/home/rubix/Sai-Rubix/Rubix-Smart-Contract-Examples/Allencontracts/Rubix-Smart-Contract-Examples/bidding-contract-createdidapi/createdid_contract/target/wasm32-unknown-unknown/debug/createdid_contract.wasm"
	schemaPath := "./createdid_contract/target/createdid_contract.json"
	rawCodePath := "/home/rubix/Sai-Rubix/Rubix-Smart-Contract-Examples/Allencontracts/Rubix-Smart-Contract-Examples/bidding-contract-createdidapi/createdid_contract/src/lib.rs"
	port := "20009"
	contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)
}

// This function is intended to pass the smart contract hash which is retruned while generating smart contract
func smartContractHash() string {
	return "QmVF2nghpdfRAbd1o7uqAwTcdd2cFRtmxWDu879sbrJUqH"
}

func DeploySmartContract() {
	/*
		deployerAddress : The peerdId.did combination of the address of the Rubix node which is deploying the contract
		port : The port corresponding to the deployer node.
	*/
	comment := "Deploying Test Bidding Contract"
	deployerAddress := "bafybmibpgv4fe4xr7wwolrymxfphe7o45r4mynnzam6ohqqzvh3usmue2e"
	quorumType := 2
	rbtAmount := 1
	smartContractToken := smartContractHash()
	port := "20009"
	id := contract.DeploySmartContract(comment, deployerAddress, quorumType, rbtAmount, smartContractToken, port)
	fmt.Println("Contract ID: " + id)
	contract.SignatureResponse(id, port)

}

func ExecuteSmartContractTestNode2() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on TestNode2"
	executorAddress := "bafybmidwfmwrq4mj74usaazwlb3hkhjuqj6wzxcwntvicxoxni3ge47myq"
	quorumType := 2
	smartContractData := `{"Type":0,"PrivPWD":"mypassword","ImgFile":"/home/rubix/Sai-Rubix/rubixgoplatform/linux/image.png"}`
	smartContractToken := smartContractHash()
	port := "20010"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractTestNode3() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on TestNode3"
	executorAddress := "bafybmihqj74dcyi3ipuzbpcqpxhyzxpr5viys6w3ethuzfxrfu37yzs4hu"
	quorumType := 2
	smartContractData := `{"Type":0,"PrivPWD":"mypassword","ImgFile":"/home/rubix/Sai-Rubix/rubixgoplatform/linux/image.png"}`
	smartContractToken := smartContractHash()
	port := "20011"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractTestNode4() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node3"
	executorAddress := "bafybmihsa7qc5onikjlxvguxifnh7xz7t57q4mqnopee62geheno4iia2m"
	quorumType := 2
	smartContractData := `{"Type":0,"PrivPWD":"mypassword","ImgFile":"/home/rubix/Sai-Rubix/rubixgoplatform/linux/image.png"}`
	smartContractToken := smartContractHash()
	port := "20012"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

// This function is responsible for subscribing to a particular smart contract.
func SubscribeSmartContractTestNode2(port string) {
	contractToken := smartContractHash()
	//	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractTestNode3(port string) {
	contractToken := smartContractHash()
	//contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractTestNode4(port string) {
	contractToken := smartContractHash()
	//contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractTestNode5(port string) {
	contractToken := smartContractHash()
	//	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractTestNode1(port string) {
	contractToken := smartContractHash()
	//contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

// Function to manually set a delay and trigger ContractExecution
func SetDelayAndTriggerContractExecution(port string, seconds int) {
	contractId := smartContractHash()
	fmt.Printf("Setting a delay of %d seconds before triggering ContractExecution...\n", seconds)
	time.After(time.Duration(seconds))
	contractExec, err := contract.NewContractExecution(contractId, port)
	smartContractTokenData := contract.GetSmartContractData(port, contractId)
	fmt.Println("Smart Contract Token Data :", string(smartContractTokenData))
	var dataReply SmartContractDataReply

	if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data reply in RunSmartContract", dataReply)
	action := contractModule.Action{
		Function: "did",
		Args:     []interface{}{""},
	}
	actions := []contractModule.Action{action}
	fmt.Println("actions in StetDelay function", actions)
	smartContractData := dataReply.SCTDataReply
	fmt.Println("Smart Contract Data :", smartContractData)
	jsonString, err := json.Marshal(smartContractData)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Print the JSON string
	fmt.Println(string(jsonString))
	contractExec.ProcessActions(actions, string(jsonString))
}

func main() {

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enlighten me with the function to be executed ")
		fmt.Println(`
		1. Generate Contract 
		2. Subscribe Contract TestNode1 aka Deployer Node
		3. Subscribe Contract TestNode2 
		4. Subscribe Contract TestNode3
		5. Subscribe Contract TestNode4
		6. Deploy Contract
		7. Execute Contract TestNode2 
		8. Execute Contract TestNode3 
		9.Execute Contract TestNode4	
		10.Creating the Did`)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Generate Contract")
			GenerateSmartContract()
		case "2":
			fmt.Println("Subscribing Smart Contract in TestNode1 aka Deployer Node")
			SubscribeSmartContractTestNode1("20009")
		case "3":
			fmt.Println("Subscribing Smart Contract in TestNode2")
			SubscribeSmartContractTestNode2("20010")
		case "4":
			fmt.Println("Subscribing Smart Contract in TestNode3")
			SubscribeSmartContractTestNode3("20011")
		case "5":
			fmt.Println("Subscribing Smart Contract in TestNode4")
			SubscribeSmartContractTestNode4("20012")
		case "6":
			fmt.Println("Deploying Smart Contract in TestNode1")
			DeploySmartContract()
		case "7":
			fmt.Println("Executing Smart Contract in TestNode2")
			ExecuteSmartContractTestNode2()
		case "8":
			fmt.Println("Executing Smart Contract in TestNode3")
			ExecuteSmartContractTestNode3()
		case "9":
			fmt.Println("Executing Smart Contract in Node 3")
			ExecuteSmartContractTestNode4()
		case "10":
			fmt.Println("Did Creation called: Creating the Did")
			SetDelayAndTriggerContractExecution("20009", 20)

		default:
			fmt.Println("You entered an unknown number")
		}
	}

}
