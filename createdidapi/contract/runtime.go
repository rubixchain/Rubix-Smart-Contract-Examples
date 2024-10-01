package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func GenerateSmartContract(did string, wasmPath string, schemaPath string, rawCodePath string, port string) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the form fields
	_ = writer.WriteField("did", did)

	// Add the binaryCodePath field
	file, _ := os.Open(wasmPath)
	defer file.Close()
	binaryPart, _ := writer.CreateFormFile("binaryCodePath", wasmPath)
	_, _ = io.Copy(binaryPart, file)

	// Add the rawCodePath field
	rawFile, _ := os.Open(rawCodePath)
	defer rawFile.Close()
	rawPart, _ := writer.CreateFormFile("rawCodePath", rawCodePath)
	_, _ = io.Copy(rawPart, rawFile)

	// Add the schemaFilePath field
	schemaFile, _ := os.Open(schemaPath)
	defer schemaFile.Close()
	schemaPart, _ := writer.CreateFormFile("schemaFilePath", schemaPath)
	_, _ = io.Copy(schemaPart, schemaFile)

	// Close the writer
	writer.Close()

	// Create the HTTP request
	url := fmt.Sprintf("http://localhost:%s/api/generate-smart-contract", port)
	req, _ := http.NewRequest("POST", url, &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in execute Contract :", string(data2))

	// Process the response as needed
	fmt.Println("Response status code:", resp.StatusCode)
}

func GetSmartContractData(port string, token string) []byte {
	data := map[string]interface{}{
		"token":  token,
		"latest": false,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	url := fmt.Sprintf("http://localhost:%s/api/get-smart-contract-token-chain-data", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
	}
	// Process the data as needed
	fmt.Println("Response Body in get smart contract data :", string(data2))

	return data2

}

func DeploySmartContract(comment string, deployerAddress string, quorumType int, rbtAmount int, smartContractToken string, port string) string {
	data := map[string]interface{}{
		"comment":            comment,
		"deployerAddr":       deployerAddress,
		"quorumType":         quorumType,
		"rbtAmount":          rbtAmount,
		"smartContractToken": smartContractToken,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	url := fmt.Sprintf("http://localhost:%s/api/deploy-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
	}
	// Process the data as needed
	fmt.Println("Response Body in deploy smart contract:", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)

	defer resp.Body.Close()
	return id

}

func SignatureResponse(requestId string, port string) string {
	data := map[string]interface{}{
		"id":       requestId,
		"mode":     0,
		"password": "mypassword",
	}

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		//	return
	}
	url := fmt.Sprintf("http://localhost:%s/api/signature-response", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		//return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		//return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		//return
	}
	// Process the data as needed
	fmt.Println("Response Body in signature response :", string(data2))
	//json encode string
	defer resp.Body.Close()
	return string(data2)
}

func ExecuteSmartContract(comment string, executorAddress string, quorumType int, smartContractData string, smartContractToken string, port string) {
	data := map[string]interface{}{
		"comment":            comment,
		"executorAddr":       executorAddress,
		"quorumType":         quorumType,
		"smartContractData":  smartContractData,
		"smartContractToken": smartContractToken,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/execute-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in execute smart contract :", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)
	SignatureResponse(id, port)

	defer resp.Body.Close()

}


func SubscribeSmartContract(contractToken string, port string) {
	fmt.Println(contractToken)
	data := map[string]interface{}{
		"smartContractToken": contractToken,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/subscribe-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in subscribe smart contract :", string(data2))

	defer resp.Body.Close()

}

func FetchSmartContract(smartContractTokenHash string, port string) {
	data := map[string]interface{}{
		"smart_contract_token": smartContractTokenHash,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/fetch-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in fetch smart contract :", string(data2))

	defer resp.Body.Close()

}

func RegisterCallBackUrl(smartContractTokenHash string, urlPort string, endPoint string, nodePort string) {
	callBackUrl := fmt.Sprintf("http://localhost:%s/%s", urlPort, endPoint)
	data := map[string]interface{}{
		"CallBackURL":        callBackUrl,
		"SmartContractToken": smartContractTokenHash,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/register-callback-url", nodePort)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	fmt.Println("Response Body in register callback url :", string(data2))
}

func ReadCurrentState(stateFilePath string) string {
	currentStateJsonFile, err := os.ReadFile(stateFilePath)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string
	currentState := string(currentStateJsonFile)
	return currentState
}

func GetRubixSmartContractPath(contractHash string, smartContractName string, nodeName string) (string, error) {
	rubixcontractPath := fmt.Sprintf("/home/rubix/Sai-Rubix/rubixgoplatform/linux/%s/SmartContract/%s/%s", nodeName, contractHash, smartContractName)

	// Check if the path exists
	if _, err := os.Stat(rubixcontractPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("smart contract path does not exist")
		}
		return "", err // Return other errors as is
	}

	return rubixcontractPath, nil
}

func GetRubixSchemaPath(contractHash string, nodeName string, schemaName string) (string, error) {
	rubixSchemaPath := fmt.Sprintf("/home/rubix/Sai-Rubix/rubixgoplatform/linux/%s/SmartContract/%s/%s", nodeName, contractHash, schemaName)

	// Check if the path exists
	if _, err := os.Stat(rubixSchemaPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("schema path does not exist")
		}
		return "", err // Return other errors as is
	}

	return rubixSchemaPath, nil
}
