package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func SyncStateJSONFile(outputFileDir string, whitelistedDIDs []string) error {
	file, err := os.OpenFile(outputFileDir, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	initJSONData := map[string][]string{"whitelisted_dids": whitelistedDIDs}

	jsonData, err := json.Marshal(initJSONData)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func AddBiddingInfoToState(did string, bid string, outputFileDir string) error {
	_, err := os.Stat(outputFileDir)
	if err == nil {
		// State file exists, read the existing data, append, and write back
		existingState, err := readExistingData(outputFileDir)
		if err != nil {
			return err
		}

        // Format DID and Bid data
        biddingInfo := map[string]string{did: bid}

		existingState["bid_info"] = append(existingState["bid_info"], biddingInfo)
		return writeUpdatedData(existingState, outputFileDir)
	}

	file, err := os.Create(outputFileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	stateInfo := map[string]interface{}{
		"whitelisted_dids": []string{did},
	}

	// Marshal the initial data to JSON
	jsonData, err := json.Marshal(stateInfo)
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func ReadState(stateDir string) ([]byte, error) {
	_, err := os.Stat(stateDir)
	if err != nil {
		return nil, fmt.Errorf("state file is not found at : %v", stateDir)
	} else {
		var stateMap map[string][]string

		fileInfo, err := os.Open(stateDir)
		if err != nil {
			return nil, err
		}
		defer fileInfo.Close()

		err = json.NewDecoder(fileInfo).Decode(&stateMap)
		if err != nil {
			return nil, err
		}

		return json.Marshal(stateMap)
	}
}

func readExistingData(outputFile string) (map[string][]map[string]string, error) {
	// Open the file for reading
	file, err := os.Open(outputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data
	var data map[string][]map[string]string
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeUpdatedData(data map[string][]map[string]string, outputFile string) error {
	// Open the file for writing
	file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshal the updated data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
