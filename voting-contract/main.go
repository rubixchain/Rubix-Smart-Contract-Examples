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
	did := "bafybmibzmgude7driixpb2hihrveiwzkvrhsogs7xzijv7zbys7qnuakvy"
	wasmPath := "/mnt/c/Users/allen/Rubix-Smart-Contract-Examples/voting-contract/voting_contract/target/wasm32-unknown-unknown/release/voting_contract.wasm"
	schemaPath := "/mnt/c/Users/allen/Rubix-Smart-Contract-Examples/voting-contract/store_state/vote_contract/votefile.json"
	rawCodePath := "/mnt/c/Users/allen/Rubix-Smart-Contract-Examples/voting-contract/voting_contract/src/lib.rs"
	port := "20003"
	contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)
}

// This function is intended to pass the smart contract hash which is retruned while generating smart contract
func smartContractHash() string {
	return "QmRct5xwgRYaDzg1qtuwRw1PvK8TC6Avwaugw877Mi6hy6"
}

func DeploySmartContract() {
	/*
		deployerAddress : The peerdId.did combination of the address of the Rubix node which is deploying the contract
		port : The port corresponding to the deployer node.
	*/
	comment := "Deploying Test Voting Contract"
	deployerAddress := "12D3KooWA9rYqCRfniLJxyMSxLq9FcNHLg9NUApPPhafiQ76enpX.bafybmibzmgude7driixpb2hihrveiwzkvrhsogs7xzijv7zbys7qnuakvy"
	quorumType := 2
	rbtAmount := 1
	smartContractToken := smartContractHash()
	port := "20003"
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
	executorAddress := "12D3KooWA9rYqCRfniLJxyMSxLq9FcNHLg9NUApPPhafiQ76enpX.bafybmibzmgude7driixpb2hihrveiwzkvrhsogs7xzijv7zbys7qnuakvy"
	quorumType := 2
	smartContractData := "Red"
	smartContractToken := smartContractHash()
	port := "20003"
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

func SubscribeSmartContractMainNode(port string) {
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
		5. Subscribe Contract Main Node
		6. Deploy Contract
		7. Execute Contract Node 1 
		8. Execute Contract Node 2 
		9. Execute Contract Node 3`)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Generate Contract")
			GenerateSmartContract()
		case "2":
			fmt.Println("Subscribing Smart Contract in Node 1")
			SubscribeSmartContractNode1("20003")
		case "3":
			fmt.Println("Subscribing Smart Contract in Node 2")
			SubscribeSmartContractNode2("20015")
		case "4":
			fmt.Println("Subscribing Smart Contract in Node 3")
			SubscribeSmartContractNode3("20017")
		case "5":
			fmt.Println("Subscribing Smart Contract in Main Node")
			SubscribeSmartContractMainNode("20002")
		case "6":
			fmt.Println("Deploying Smart Contract")
			DeploySmartContract()
		case "7":
			fmt.Println("Executing Smart Contract in Node 1")
			ExecuteSmartContractNode1()
		case "8":
			fmt.Println("Executing Smart Contract in Node 2")
			ExecuteSmartContractNode2()
		case "9":
			fmt.Println("Executing Smart Contract in Node 3")
			ExecuteSmartContractNode3()
		default:
			fmt.Println("You entered an unknown number")
		}
	}

}
