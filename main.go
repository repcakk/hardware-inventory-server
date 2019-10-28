package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	db "github.com/repcakk/hardware-inventory-server/database"
)

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

func updateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	computerID := r.FormValue("computerId")
	gpuName := r.FormValue("gpuName")

	db.AddOrUpdateRow(computerID, gpuName)
}

func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	hardwareInventoryTemplate, _ := template.ParseFiles("hardware-inventory.html")
	hardwareInventoryTemplate.Execute(w, db.GetRows())
}

func main() {
	defer db.CloseHardwareDB()
	defer db.CloseUserDB()

	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/inventory", inventoryHandler)
	port := config.Port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
