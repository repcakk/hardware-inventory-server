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

func readCommand() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(" >> ") // command indicator
	command, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	command = strings.TrimSpace(command)
	fmt.Printf("\n")
	return command
}

func loadOrSaveDB(operationType string) {
	fmt.Printf("\n Select database:\n" +
		"     users           database contains map of all hostname and username pairs\n" +
		"     gpu-all         database contains map of all GPU-serial-number and gpu-name pairs\n" +
		"     gpu-in-use      database contains map of all used in computers GPU-serial-number and hostname pairs\n")

	dbType := readCommand()

	fmt.Printf(" Provide path for " + operationType + "\n")
	path := readCommand()

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
	web.Init(port)
	web.Run()

	var quit bool = false
	var maintenanceMode bool = false
	for !quit {
		if maintenanceMode {
			fmt.Printf(" Select option:\n" +
				"     load            load database from json file\n" +
				"     save            save database to json file\n" +
				"     resume          run server again and stop maintenance mode\n" +
				"     quit            quit\n")
		} else {
			fmt.Printf(" Select option:\n" +
				"     stop            stop server and enter maintenance mode\n" +
				"     quit            quit\n")
		}

		command := readCommand()
		switch command {
		case "quit":
			quit = true
		case "stop":
			web.Shutdown()
			maintenanceMode = true
			fmt.Printf("" +
				" ##########################################################\n" +
				" ##### Http server stopped - maintenance mode enabled #####\n" +
				" ##########################################################" +
				"\n\n")
		case "resume":
			maintenanceMode = false
			fmt.Printf("" +
				" ##########################################################\n" +
				" ### Http server is running - maintenance mode disabled ###\n" +
				" ##########################################################" +
				"\n\n")
			web.Init(port)
			web.Run()
		case "load":
			loadOrSaveDB(command)
		case "save":
			loadOrSaveDB(command)
		}
	}

	web.Shutdown()
}

func main() {
	defer db.GpuAllDB.Close()
	defer db.GpuInUseDB.Close()
	defer db.UserDB.Close()

	appController()

}
