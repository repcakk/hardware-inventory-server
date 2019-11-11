package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSON loads key-value pair JSON
// into map[string]string data structure
func ReadJSON(filePath string) map[string]string {
	keyValueMap := make(map[string]string)

	srcFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(srcFile)
	jsonParser.Decode(&keyValueMap)

	return keyValueMap
}

// WriteJSON writes key-value pair JSON
// from map[string]string data structure
func WriteJSON(filePath string, mapToDump map[string]string) {
	srcFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewEncoder(srcFile)
	jsonParser.Encode(&mapToDump)
}
