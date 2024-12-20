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
	did := "bafybmid23tc2ddfsmxbb5h26cyiljkylbp7fjsgpbucxlm43vzbm5p2b7y"
	wasmPath := "./data/contracts/bidding_contract.wasm"
	schemaPath := "./data/state/bidding_contract.json"
	rawCodePath := "./data/raw-code/lib.rs"
	port := "20002"
	contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)
}

// This function is intended to pass the smart contract hash which is retruned while generating smart contract
func smartContractHash() string {
	return "QmYsAL5cQpXB6zSLyJKoJcedBofMsYwFpUajbXJ3JWwqy9"
}

func DeploySmartContract() {
	/*
		deployerAddress : The peerdId.did combination of the address of the Rubix node which is deploying the contract
		port : The port corresponding to the deployer node.
	*/
	comment := "Deploying Test Bidding Contract"
	deployerAddress := "12D3KooWQK7SKtBo91gesKpe71gWAJQqsKvM1JLAsDTe8fBC6hFZ.bafybmid23tc2ddfsmxbb5h26cyiljkylbp7fjsgpbucxlm43vzbm5p2b7y"
	quorumType := 2
	rbtAmount := 1
	smartContractToken := smartContractHash()
	port := "20002"
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
	executorAddress := "12D3KooWBJw9Ei6L7FkSLyVX3PBfaVGRGbmBjaHaGgUNCMGTHMWp.bafybmidpuc3vcegwxjrr6xnvvwvt3txmhkjnvdgb2kbyvemilw4lptxlvm"
	quorumType := 2
	smartContractData := `{"did":"bafybmidpuc3vcegwxjrr6xnvvwvt3txmhkjnvdgb2kbyvemilw4lptxlvm","bid":40.01}`
	smartContractToken := smartContractHash()
	port := "20003"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractTestNode3() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on TestNode3"
	executorAddress := "12D3KooWJp7uNRG9Cv1YPkernpw53WaXTELWFBged3gqePswguDB.bafybmicothuuyezh3rffamavaghxhcvgwisg5f5dyuhz26ljtpvzueclby"
	quorumType := 2
	smartContractData := `{"did":"bafybmicothuuyezh3rffamavaghxhcvgwisg5f5dyuhz26ljtpvzueclby","bid":50.00}`
	smartContractToken := smartContractHash()
	port := "20004"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractTestNode4() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node3"
	executorAddress := "12D3KooWMJ6hBM9R8uF8Lfo9QUM6YgV522DsJkeiZLTgfcQmbgrU.bafybmidlbjr35spo2wdyhonoslew22jmoagdcxklgqacho5bie4bl5ksiy"
	quorumType := 2
	smartContractData := `{"did":"bafybmidlbjr35spo2wdyhonoslew22jmoagdcxklgqacho5bie4bl5ksiy","bid":60.01}`
	smartContractToken := smartContractHash()
	port := "20005"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractTestNode5() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node3"
	executorAddress := "12D3KooWQ8smfu9WeGqYDrZmCboSdiPyvVKMcDQoDUdZeb9eeQR1.bafybmidaqxue2u5ayicgyqpiv6yiir2qyiw6lhr3cgq7oru7ng25lupu2y"
	quorumType := 2
	smartContractData := `{"did":"bafybmidaqxue2u5ayicgyqpiv6yiir2qyiw6lhr3cgq7oru7ng25lupu2y","bid":60.01}`
	smartContractToken := smartContractHash()
	port := "20006"
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
	contractExec, err := contractModule.NewContractExecution(contractId, port)
	smartContractTokenData := contract.GetSmartContractData(port, contractId)
	fmt.Println("Smart Contract Token Data :", string(smartContractTokenData))
	var dataReply SmartContractDataReply

	if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data reply in RunSmartContract", dataReply)
	action := contractModule.Action{
		Function: "bid",
		Args:     []interface{}{""},
	}
	actions := []contractModule.Action{action}
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
		6. Subscribe Contract TestNode5
		7. Deploy Contract
		8. Execute Contract TestNode2 
		9. Execute Contract TestNode3 
		10.Execute Contract TestNode4
		11.Execute Contract TestNode5
		12.Find the Highest Bidder`)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Generate Contract")
			GenerateSmartContract()
		case "2":
			fmt.Println("Subscribing Smart Contract in TestNode1 aka Deployer Node")
			SubscribeSmartContractTestNode2("20002")
		case "3":
			fmt.Println("Subscribing Smart Contract in TestNode2")
			SubscribeSmartContractTestNode2("20003")
		case "4":
			fmt.Println("Subscribing Smart Contract in TestNode3")
			SubscribeSmartContractTestNode3("20004")
		case "5":
			fmt.Println("Subscribing Smart Contract in TestNode4")
			SubscribeSmartContractTestNode4("20005")
		case "6":
			fmt.Println("Subscribing Smart Contract in TestNode5")
			SubscribeSmartContractTestNode4("20006")
		case "7":
			fmt.Println("Deploying Smart Contract in TestNode1")
			DeploySmartContract()
		case "8":
			fmt.Println("Executing Smart Contract in TestNode2")
			ExecuteSmartContractTestNode2()
		case "9":
			fmt.Println("Executing Smart Contract in TestNode3")
			ExecuteSmartContractTestNode3()
		case "10":
			fmt.Println("Executing Smart Contract in Node 3")
			ExecuteSmartContractTestNode4()
		case "11":
			fmt.Println("Executing Smart Contract in Node 3")
			ExecuteSmartContractTestNode5()
		case "12":
			fmt.Println("Bidding is Completed: Finding the highest bid")
			SetDelayAndTriggerContractExecution("20002", 20)

		default:
			fmt.Println("You entered an unknown number")
		}
	}

}
