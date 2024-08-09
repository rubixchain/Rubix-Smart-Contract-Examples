package main

import (
	"bidding-contract/contract"
	"bufio"
	"encoding/hex"
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
	wasmPath := "./bidding_contract/target/wasm32-unknown-unknown/debug/bidding_contract.wasm"
	schemaPath := "./data/state/bidding_contract.json"
	rawCodePath := "./bidding_contract/src/lib.rs"
	port := "20009"
	contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)
}

// This function is intended to pass the smart contract hash which is retruned while generating smart contract
func smartContractHash() string {
	return "QmaCGpmpTjY184ag9mMVX4LHHqxqrBgb5p5A8cP13Rau9J"
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
func EncryptSmartcontractData(smartcontractData string, pubkeypath string) string {
	smartcontractData_bytes := []byte(smartcontractData)
	encryptedsmartcontractdata := (contract.Ecies_encryption(pubkeypath, smartcontractData_bytes))
	fmt.Println("encrypted smartcontract data byes before encoding", encryptedsmartcontractdata)
	encryptedsmartcontractdata_string := hex.EncodeToString(encryptedsmartcontractdata)

	// Print the hexadecimal string
	//fmt.Println("Hexadecimal format:", encryptedsmartcontractdata_string)
	return encryptedsmartcontractdata_string
}
func ExecuteSmartContractTestNode2() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on TestNode2"
	executorAddress := "bafybmidwfmwrq4mj74usaazwlb3hkhjuqj6wzxcwntvicxoxni3ge47myq"
	quorumType := 2
	smartContractData := `{"did":"bafybmidwfmwrq4mj74usaazwlb3hkhjuqj6wzxcwntvicxoxni3ge47myq","bid":20.00}`
	smartContractToken := smartContractHash()
	pubkey_path := "/home/rubix/Sai-Rubix/rubixgoplatform/linux/node9/Rubix/TestNetDID/bafybmiahqapy3fvpqn4feoawdotnewu3zh4cpne54uivaceb2bpd2ihnja/pubKey.pem"
	encrypted_smartcontractdata := EncryptSmartcontractData(smartContractData, pubkey_path)
	// plaintex_bytes := []byte(smartContractData)
	// encrypted_smartcontractdata := (contract.Ecies_encryption(pubkey_path, plaintex_bytes))
	fmt.Println("encryptedsmartcontract data of node10 is ", encrypted_smartcontractdata)
	port := "20010"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, encrypted_smartcontractdata, smartContractToken, port)
}

func ExecuteSmartContractTestNode3() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on TestNode3"
	executorAddress := "bafybmihqj74dcyi3ipuzbpcqpxhyzxpr5viys6w3ethuzfxrfu37yzs4hu"
	quorumType := 2
	smartContractData := `{"did":"bafybmihqj74dcyi3ipuzbpcqpxhyzxpr5viys6w3ethuzfxrfu37yzs4hu","bid":22.00}`
	smartContractToken := smartContractHash()
	port := "20011"
	pubkey_path := "/home/rubix/Sai-Rubix/rubixgoplatform/linux/node9/Rubix/TestNetDID/bafybmiahqapy3fvpqn4feoawdotnewu3zh4cpne54uivaceb2bpd2ihnja/pubKey.pem"
	// plaintex_bytes := []byte(smartContractData)
	// encrypted_smartcontractdata := (contract.Ecies_encryption(pubkey_path, plaintex_bytes))
	encrypted_smartcontractdata := EncryptSmartcontractData(smartContractData, pubkey_path)
	fmt.Println("encryptedsmartcontract data of node11 is ", encrypted_smartcontractdata)

	contract.ExecuteSmartContract(comment, executorAddress, quorumType, encrypted_smartcontractdata, smartContractToken, port)
}

func ExecuteSmartContractTestNode4() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node3"
	executorAddress := "bafybmihsa7qc5onikjlxvguxifnh7xz7t57q4mqnopee62geheno4iia2m"
	quorumType := 2
	smartContractData := `{"did":"bafybmif4xcderso3aouqhwwdeeffhknbosd4mgtqphi6rznrtrekvdo6fq","bid":24.00}`
	smartContractToken := smartContractHash()
	port := "20012"
	pubkey_path := "/home/rubix/Sai-Rubix/rubixgoplatform/linux/node9/Rubix/TestNetDID/bafybmiahqapy3fvpqn4feoawdotnewu3zh4cpne54uivaceb2bpd2ihnja/pubKey.pem"
	// plaintex_bytes := []byte(smartContractData)
	// encrypted_smartcontractdata := (contract.Ecies_encryption(pubkey_path, plaintex_bytes))
	encrypted_smartcontractdata := EncryptSmartcontractData(smartContractData, pubkey_path)
	fmt.Println("encryptedsmartcontract data of node12 is ", encrypted_smartcontractdata)
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, encrypted_smartcontractdata, smartContractToken, port)
}

func ExecuteSmartContractTestNode5() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node3"
	executorAddress := "12D3KooWQ8smfu9WeGqYDrZmCboSdiPyvVKMcDQoDUdZeb9eeQR1.bafybmidaqxue2u5ayicgyqpiv6yiir2qyiw6lhr3cgq7oru7ng25lupu2y"
	quorumType := 2
	smartContractData := `{"did":"bafybmidaqxue2u5ayicgyqpiv6yiir2qyiw6lhr3cgq7oru7ng25lupu2y","bid":26.01}`
	smartContractToken := smartContractHash()
	port := "20006"
	pubkey_path := "/home/rubix/Sai-Rubix/rubixgoplatform/linux/node9/Rubix/TestNetDID/bafybmiahqapy3fvpqn4feoawdotnewu3zh4cpne54uivaceb2bpd2ihnja/pubKey.pem"
	// plaintex_bytes := []byte(smartContractData)
	// encrypted_smartcontractdata := (contract.Ecies_encryption(pubkey_path, plaintex_bytes))
	encrypted_smartcontractdata := EncryptSmartcontractData(smartContractData, pubkey_path)

	contract.ExecuteSmartContract(comment, executorAddress, quorumType, encrypted_smartcontractdata, smartContractToken, port)
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
	if err != nil {
		fmt.Println("new contract execution failed,err: ", err)
		// return
	}
	smartContractTokenData := contract.GetSmartContractData(port, contractId)
	fmt.Println("Smart Contract Token Data :", (smartContractTokenData))
	var dataReply SmartContractDataReply

	if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data reply in RunSmartContract", dataReply)
	//Decrypting Sct data
	privatekeyBytes := []byte{26, 74, 206, 110, 150, 148, 87, 32, 213, 102, 150, 120, 224, 105, 131, 103, 58, 95, 72, 72, 142, 240, 97, 25, 113, 39, 140, 138, 164, 82, 187, 147}

	decryptedSCTDataReply := dataReply.SCTDataReply
	for i, reply := range decryptedSCTDataReply {
		encryptedsmartcontractdata_byes, err := hex.DecodeString(reply.SmartContractData)
		if err != nil {
			fmt.Println("unable to decode", err)
		}
		reply.SmartContractData = contractModule.Ecies_decryption(privatekeyBytes, encryptedsmartcontractdata_byes)
		fmt.Println("Decrypted smartcontract data look at it:", reply.SmartContractData)
		decryptedSCTDataReply[i] = reply
	}
	action := contractModule.Action{
		Function: "bid",
		Args:     []interface{}{""},
	}
	actions := []contractModule.Action{action}
	smartContractData := decryptedSCTDataReply
	fmt.Println("Smart Contract Data :", smartContractData)
	jsonString, err := json.Marshal(smartContractData)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Print the JSON string
	fmt.Println("sct data in json str:", jsonString)
	contractExec.ProcessActions(actions, jsonString)
}
func main() {

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enlighten me with the function to be executed ")
		fmt.Println(`
		1. Generate Contract 
		2. Subscribe Contract TestNode9 aka Deployer Node
		3. Subscribe Contract TestNode10 
		4. Subscribe Contract TestNode11
		5. Subscribe Contract TestNode12
		6. Subscribe Contract TestNode5
		7. Deploy Contract
		8. Execute Contract TestNode10 
		9. Execute Contract TestNode11
		10.Execute Contract TestNode12
		11.Execute Contract TestNode5
		12.Find the Highest Bidder`)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Generate Contract")
			GenerateSmartContract()
		case "2":
			fmt.Println("Subscribing Smart Contract in TestNode9 aka Deployer Node")
			SubscribeSmartContractTestNode2("20009")
		case "3":
			fmt.Println("Subscribing Smart Contract in TestNode10")
			SubscribeSmartContractTestNode2("20010")
		case "4":
			fmt.Println("Subscribing Smart Contract in TestNode11")
			SubscribeSmartContractTestNode3("20011")
		case "5":
			fmt.Println("Subscribing Smart Contract in TestNode12")
			SubscribeSmartContractTestNode4("20012")
		case "6":
			fmt.Println("Subscribing Smart Contract in TestNode5")
			SubscribeSmartContractTestNode4("20006")
		case "7":
			fmt.Println("Deploying Smart Contract in TestNode9")
			DeploySmartContract()
		case "8":
			fmt.Println("Executing Smart Contract in TestNode10")
			ExecuteSmartContractTestNode2()
		case "9":
			fmt.Println("Executing Smart Contract in TestNode11")
			ExecuteSmartContractTestNode3()
		case "10":
			fmt.Println("Executing Smart Contract in Node 12")
			ExecuteSmartContractTestNode4()
		case "11":
			fmt.Println("Executing Smart Contract in Node 3")
			ExecuteSmartContractTestNode5()
		case "12":
			fmt.Println("Bidding is Completed: Finding the highest bid")
			SetDelayAndTriggerContractExecution("20009", 20)

		default:
			fmt.Println("You entered an unknown number")
		}
	}

}
