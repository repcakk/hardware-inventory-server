package main

import (
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

	web.Run()
}
