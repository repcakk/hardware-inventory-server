package database

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prologic/bitcask"
)

var userDB, _ = bitcask.Open("data/user-database")
var userComputerMap = loadUsernames("config/usernames-mapping.json")

func loadUsernames(file string) map[string]string {
	usernamesMap := make(map[string]string)

	usernamesFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(usernamesFile)
	jsonParser.Decode(&usernamesMap)

	return usernamesMap
}

// CloseUserDB close users database file
func CloseUserDB() {
	defer userDB.Close()
}
