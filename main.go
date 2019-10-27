package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/prologic/bitcask"
)

var db, _ = bitcask.Open("hardware-database")
var userComputerMap = loadUsernames("config/usernames-mapping.json")
var config = loadConfig("config/config.json")

// Config structure for storing server properties
type Config struct {
	Port string `json:"port"`
}

func loadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

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

func updateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	computerID := r.FormValue("computerId")
	gpuName := r.FormValue("gpuName")

	addOrUpdateRow(computerID, gpuName)
}

func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	hardwareInventoryTemplate, _ := template.ParseFiles("hardware-inventory.html")
	hardwareInventoryTemplate.Execute(w, getRows())
}

func addOrUpdateRow(computerID string, gpuName string) {
	db.Put([]byte(computerID), []byte(gpuName))
	db.Merge()
}

func getRows() map[string]string {
	computersGpuMap := make(map[string]string)
	for key := range db.Keys() {
		value, _ := db.Get([]byte(key))
		computersGpuMap[string([]byte(key))] = string(value)
	}

	userGpuMap := make(map[string]string)
	for computerID, gpuName := range computersGpuMap {
		userGpuMap[userComputerMap[computerID]] = gpuName
	}

	return userGpuMap
}

func main() {
	defer db.Close() // db will close after main returns

	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/inventory", inventoryHandler)
	port := config.Port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
