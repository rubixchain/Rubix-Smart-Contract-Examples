package server

import (
	"encoding/json"
	"fmt"
)

func unmarshalWASMOutputStringToMap(v string) (map[string]string, error) {
	var m map[string]string

	err := json.Unmarshal([]byte(v), &m)
	if err != nil {
		return nil, fmt.Errorf("unmarshalWASMOutputStringToMap: %v", err)
	}

	return m, nil
}

func unmarshalSmartContractDataString(v string) (map[string]map[string]string, error) {
	var m map[string]map[string]string 

	err := json.Unmarshal([]byte(v), m)
	if err != nil {
		return nil, fmt.Errorf("unmarshalSmartContractDataString: %v", err)
	}

	return m, nil
}