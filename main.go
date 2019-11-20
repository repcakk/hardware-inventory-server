package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	db "github.com/repcakk/hardware-inventory-server/database"
	web "github.com/repcakk/hardware-inventory-server/web"
)

// var config = loadConfig("config/config.json")

// // Config structure for storing server properties
// type Config struct {
// 	Port string `json:"port"`
// }

// func loadConfig(file string) Config {
// 	var config Config
// 	configFile, err := os.Open(file)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	jsonParser := json.NewDecoder(configFile)
// 	jsonParser.Decode(&config)
// 	return config
// }

var port string = "8080" //config.Port

func loadOrSaveDB(operationType string) {
	fmt.Printf(" Select database:\n" +
		"  users   - database contains map of all hostname and username pairs\n" +
		"  gpu-all - database contains map of all GPU-serial-number and gpu-name pairs\n" +
		"  gpu-in-use - database contains map of all used in computers GPU-serial-number and hostname pairs\n\n")

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(" ") // whitespace before command
	dbType, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(" Provide path to load from:\n")
	fmt.Printf(" ") // whitespace before command
	path, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	switch dbType {
	case "users":
		if operationType == "save" {
			db.UserDB.SaveDatabaseToJSON(path)
		} else {
			db.UserDB.LoadDatabaseFromJSON(path)
		}
	case "gpu-all":
		if operationType == "save" {
			db.GpuAllDB.SaveDatabaseToJSON(path)
		} else {
			db.GpuAllDB.LoadDatabaseFromJSON(path)
		}
	case "gpu-in-use":
		if operationType == "save" {
			db.GpuInUseDB.SaveDatabaseToJSON(path)
		} else {
			db.GpuInUseDB.LoadDatabaseFromJSON(path)
		}
	}
}

func appController() {
	var quit bool = false
	reader := bufio.NewReader(os.Stdin)

	var maintenanceMode bool = false
	for !quit {
		if maintenanceMode {
			fmt.Printf(" Select option:\n" +
				"  load <path to json file> - load database from json file\n" +
				"  save <path to save json> - save database to json file\n" +
				"  run - run server again and stop maintenance mode\n" +
				"  quit - quit\n")
		} else {
			fmt.Printf(" Select option:\n" +
				"  stop - stop server and enter maintenance mode\n" +
				"  quit - quit\n")
		}
		fmt.Printf(" ") // whitespace before command
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		command = strings.TrimSpace(command)

		switch command {
		case "quit":
			quit = true
		case "stop":
			web.Shutdown()
			maintenanceMode = true
			fmt.Printf(" ### Http server stopped - maintenance mode enabled ###\n\n")
		case "run":
			maintenanceMode = false
			fmt.Printf(" ### Http server is running - maintenance mode disabled ###\n\n")
			web.Init(port)
			web.Run()
		case "load":
			loadOrSaveDB(command)
		case "save":
			loadOrSaveDB(command)
		}
	}
}

func main() {
	defer db.GpuAllDB.Close()
	defer db.GpuInUseDB.Close()
	defer db.UserDB.Close()

	db.UserDB.LoadDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/userDetails.json")
	db.GpuAllDB.LoadDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/gpuDetails.json")
	db.GpuInUseDB.LoadDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/userGpuDetails.json")

	web.Init(port)
	web.Run()

	appController()

	web.Shutdown()

}
