package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"voting-contract/contract"
	"voting-contract/server"
)

func GenerateSmartContract() {
	/*
		This did, wasmPath, schemaPath, rawcodePath and Port should be replaced according to your Rubix node configuration and
		the respective paths
	*/
	did := "bafybmid5yd3obfpcximuaoe37a564ycpzdpxqy5o7ephu3m5fsug2wqgyy"
	wasmPath := "voting_contract/target/wasm32-unknown-unknown/release/voting_contract.wasm"
	schemaPath := "store_state/vote_contract/votefile.json"
	rawCodePath := "voting_contract/src/lib.rs"
	port := "20016"
	contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)
}

// This function is intended to pass the smart contract hash which is retruned while generating smart contract
func smartContractHash() string {
	return "QmfSEMA6f9aqZa5cfQmXXfmHDogfpyEpjPxKoWcYSv5xSj"
}

func DeploySmartContract() {
	/*
		deployerAddress : The peerdId.did combination of the address of the Rubix node which is deploying the contract
		port : The port corresponding to the deployer node.
	*/
	comment := "Deploying Test Voting Contract"
	deployerAddress := "12D3KooWKdMPfF58y5mndzXhz9dp92QgWjo82s9oWfsJAZRkw7zY.bafybmid5yd3obfpcximuaoe37a564ycpzdpxqy5o7ephu3m5fsug2wqgyy"
	quorumType := 2
	rbtAmount := 1
	smartContractToken := smartContractHash()
	port := "20016"
	id := contract.DeploySmartContract(comment, deployerAddress, quorumType, rbtAmount, smartContractToken, port)
	fmt.Println("Contract ID: " + id)
	contract.SignatureResponse(id, port)

}

func ExecuteSmartContractNode1() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node1"
	executorAddress := "12D3KooWEP8wPaUGofjRSQRfc6fXaEXuwSAUnx7G7omdcZrAfJ5B.bafybmigbwhqopv5gr6i22dqhryg35kvlvokn7nxlss4nnzdprfpgxnxsga"
	quorumType := 2
	smartContractData := "Red"
	smartContractToken := smartContractHash()
	port := "20014"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractNode2() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node2"
	executorAddress := "12D3KooWBnryPQk82qvK6it4aTUqcGhk13YkxKFZejkEpgpAG9X6.bafybmidmdayavtdzolasw3wye7laese6m2fi3mcfa5o3cfkedtpe6utklq"
	quorumType := 2
	smartContractData := "Blue"
	smartContractToken := smartContractHash()
	port := "20015"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractNode3() {
	/*
		executorAddress : The peerdId.did combination of the address of the Rubix node which is execcting the contract
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node3"
	executorAddress := "12D3KooWGMbtc77iZs5Aw59vg97M6teTmLam9wyWCfti7ybo2ZdD.bafybmiejhevgpzgbxqugjl7w7mrrhafefvvolel7jfshbybtgazdqych2u"
	quorumType := 2
	smartContractData := "Red"
	smartContractToken := smartContractHash()
	port := "20017"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

// This function is responsible for subscribing to a particular smart contract.
func SubscribeSmartContractNode1(port string) {
	contractToken := smartContractHash()
	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractNode2(port string) {
	contractToken := smartContractHash()
	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractNode3(port string) {
	contractToken := smartContractHash()
	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func main() {
	go server.Bootup()

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enlighten me with the function to be executed ")
		fmt.Println(`
		1. Generate Contract 
		2. Subscribe Contract Node 1
		3. Subscribe Contract Node 2 
		4. Subscribe Contract Node 3 
		5. Deploy Contract
		6. Execute Contract Node 1 
		7. Execute Contract Node 2 
		8. Execute Contract Node 3`)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Generate Contract")
			GenerateSmartContract()
		case "2":
			fmt.Println("Subscribing Smart Contract in Node 1")
			SubscribeSmartContractNode1("20014")
		case "3":
			fmt.Println("Subscribing Smart Contract in Node 2")
			SubscribeSmartContractNode2("20015")
		case "4":
			fmt.Println("Subscribing Smart Contract in Node 3")
			SubscribeSmartContractNode3("20017")
		case "5":
			fmt.Println("Deploying Smart Contract")
			DeploySmartContract()
		case "6":
			fmt.Println("Executing Smart Contract in Node 1")
			ExecuteSmartContractNode1()
		case "7":
			fmt.Println("Executing Smart Contract in Node 2")
			ExecuteSmartContractNode2()
		case "8":
			fmt.Println("Executing Smart Contract in Node 3")
			ExecuteSmartContractNode3()
		default:
			fmt.Println("You entered an unknown number")
		}
	}

}
