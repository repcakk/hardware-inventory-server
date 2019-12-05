package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	db "github.com/repcakk/hardware-inventory-server/database"
)

// serverConfig structure for storing server properties
type serverConfig struct {
	ServerPort string `json:"serverProt"`
}

var config = loadConfig("./config/server-config.json")

// loadConfig loads web server config
func loadConfig(file string) serverConfig {
	var config serverConfig
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
	macAddress := r.FormValue("macAddress")
	hostname := r.FormValue("hostname")
	userName := r.FormValue("username")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	gpuSN := r.FormValue("gpuSN")
	gpuName := r.FormValue("gpuName")

	gpu := db.Gpu{SN: gpuSN, GpuName: gpuName}
	gpu, _ = db.AddGpu(gpu)

	user := db.User{Username: userName, Surname: surname, Email: email}
	user, _ = db.AddUser(user)

	computer := db.Computer{MacAddress: macAddress, Hostname: hostname, CurrentGpuID: gpu.ID, LastGpuID: gpu.ID, UserID: user.ID}
	computer, isComputerNew := db.AddComputer(computer)

	if !isComputerNew {
		db.ChangeGpuInComputer(computer.ID, gpu.ID)
	}
}

func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	hardwareInventoryTemplate, err := template.ParseFiles("static/html/hardware-inventory.html")

	if err != nil {
		log.Fatal(err)
	}

	computersInfo := db.GetComputersGpusUsers(0, 100)

	err = hardwareInventoryTemplate.Execute(w, computersInfo)

	if err != nil {
		log.Fatal(err)
	}
}

// func moveToStorage(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	gpuSN := r.FormValue("gpuSN")

// 	fmt.Println("MOVE TO STORAGE CALLED")
// 	fmt.Println(gpuSN)
// 	db.MoveGpuToStorage(gpuSN)
// }

var serveMux *http.ServeMux
var server http.Server

// Init initialize server
func Init() error {
	serveMux = http.NewServeMux()
	serveMux.HandleFunc("/update", updateHandler)
	serveMux.HandleFunc("/inventory", inventoryHandler)
	// serveMux.HandleFunc("/move_to_storage", moveToStorage)
	server = http.Server{Addr: ":" + config.ServerPort, Handler: serveMux}
	// TODO: Improve error handling
	return nil
}

// Run runs http server and perform http handling
func Run() error {
	go server.ListenAndServe()
	// TODO: Improve error handling
	return nil
}

// Shutdown server
func Shutdown() {
	server.Shutdown(context.Background())
}
