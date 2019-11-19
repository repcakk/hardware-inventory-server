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

func appController() {
	var quit bool = false
	reader := bufio.NewReader(os.Stdin)

	var maintenanceMode bool = false
	for !quit {
		if !maintenanceMode {
			fmt.Printf(" Select option:\n" +
				"  s - stop server and enter maintenance mode\n" +
				"  q - quit\n")
		}

		result, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		result = strings.TrimSpace(result)

		switch result {
		case "q":
			quit = true
		case "s":
			web.Shutdown()
			maintenanceMode = true
			fmt.Printf(" ### Http server stopped - maintenance mode enabled ###\n\n")
		case "r":
			maintenanceMode = false
			fmt.Printf(" ### Http server is running - maintenance mode disabled ###\n\n")
			web.Init(port)
			web.Run()
		}

		if maintenanceMode {
			fmt.Printf(" Select option:\n" +
				"  l <path to json file> - load database from json file\n" +
				"  s <path to save json> - save database to json file\n" +
				"  r - run server again and stop maintenance mode\n" +
				"  q - quit\n")
		}
	}
}

func main() {
	defer db.GpuAllDB.Close()
	defer db.GpuInUseDB.Close()
	defer db.UserDB.Close()

	db.GpuAllDB.OverwriteDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/gpuDetails.json")
	db.UserDB.OverwriteDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/userDetails.json")
	db.GpuInUseDB.OverwriteDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/userGpuDetails.json")

	web.Init(port)
	web.Run()

	appController()

	web.Shutdown()

}
