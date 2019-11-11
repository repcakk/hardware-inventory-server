package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSON loads key-value pair JSON
// into map[string]string data structure
func ReadJSON(file string) map[string]string {
	keyValueMap := make(map[string]string)

	srcFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(srcFile)
	jsonParser.Decode(&keyValueMap)

	return keyValueMap
}

// WriteJSON writes key-value pair JSON
// from map[string]string data structure
func WriteJSON(file string) {
	keyValueMap := make(map[string]string)

	srcFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewEncoder(srcFile)
	jsonParser.Encode(&keyValueMap)
}