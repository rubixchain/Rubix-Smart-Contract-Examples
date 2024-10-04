package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// Requests

type GetSmartContractDataRequest struct {
	Token string `json:"token"`
	Latest bool `json:"latest"`
}

// Reponses

type SCTDataReply struct {
	BlockNo           uint64
	BlockId           string
	SmartContractData string
}

type SmartContractDataReply struct {
	BasicResponse
	SCTDataReply []SCTDataReply
}

type BasicResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func GetLatestSmartContractData() ([]SCTDataReply, error) {
	smartContractHash := os.Getenv("BIDDING_CONTRACT_HASH")
	rubixNodePort := os.Getenv("RUBIX_NODE_PORT")

	requestData := &GetSmartContractDataRequest{
		Token:  smartContractHash,
		Latest: true,
	}


	smartContractDataResponse, err := getSmartContractData(requestData, rubixNodePort)
	if err != nil {
		return nil, err
	}

	if len(smartContractDataResponse.SCTDataReply) != 1 {
		return nil, fmt.Errorf("Invalid format for SCDataReply received for Contract: %v\n", requestData.Token)
	} else {
		return smartContractDataResponse.SCTDataReply, nil
	}
}

func GetAllSmartContractData() ([]SCTDataReply, error) {
	smartContractHash := os.Getenv("BIDDING_CONTRACT_HASH")
	rubixNodePort := os.Getenv("RUBIX_NODE_PORT")

	requestData := &GetSmartContractDataRequest{
		Token:  smartContractHash,
		Latest: false,
	}

	smartContractDataResponse, err := getSmartContractData(requestData, rubixNodePort)
	if err != nil {
		return []SCTDataReply{}, err
	}

	return smartContractDataResponse.SCTDataReply, nil
}

func getSmartContractData(requestData *GetSmartContractDataRequest, rubixNodePort string) (*SmartContractDataReply, error) {
	bodyJSON, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling JSON:", err)
	}
	url := fmt.Sprintf("http://localhost:%s/api/get-smart-contract-token-chain-data", rubixNodePort)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request:", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending HTTP reques to Get Smart contract data:", err)
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s\n", err)
	}

	var smartContractData SmartContractDataReply

	if err := json.Unmarshal(data2, &smartContractData); err != nil {
		return nil, fmt.Errorf("Err unable to unmarshal SmartContractDataReply: %v\n", err)
	}

	return &smartContractData, nil

	// if len(smartContractData.SCTDataReply) != 1 {
	// 	return "", fmt.Errorf("Invalid format for SCDataReply received for Contract: %v\n", requestData.Token)
	// } else {
	// 	return smartContractData.SCTDataReply[0].SmartContractData, nil
	// }
}
