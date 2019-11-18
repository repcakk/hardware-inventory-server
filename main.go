package main

import (
	"bufio"
	"fmt"
	"os"

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

func main() {
	defer db.GpuAllDB.Close()
	defer db.GpuInUseDB.Close()
	defer db.UserDB.Close()

	db.GpuAllDB.OverwriteDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/gpuDetails.json")
	db.UserDB.OverwriteDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/userDetails.json")
	db.GpuInUseDB.OverwriteDatabaseFromJSON("C:/Users/repca/Desktop/inv_test_data/userGpuDetails.json")

	port := "8080" //config.Port

	go web.Run(port)

	var quit bool = false
	reader := bufio.NewReader(os.Stdin)

	for !quit {
		fmt.Printf("Select option:\n q - quit\n s - suspend server and enter maintenance mode\n c - continuue and return from maintenance mode\n")

		result, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println(err)
			return
		}

		if result == 'q' {
			quit = true
		}
	}

	web.Shutdown()

}
