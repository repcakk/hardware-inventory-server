package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	db "github.com/repcakk/hardware-inventory-server/database"
	web "github.com/repcakk/hardware-inventory-server/web"
)

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

func main() {

	// Initialize GORM connection to postgreSQL database for all ORM related stuff and inserts.
	err := db.ConnectGormDB()
	if err != nil {
		fmt.Println(err)
	}
	defer db.CloseGormDB()
	db.MigrateSchema()

	// Initialize web server.
	web.Init()
	web.Run()
	defer web.Shutdown()

	//db.GormTest()

	fmt.Printf("\n### Hardware Inventory server is running ###\n")
	fmt.Printf(" Commands:\n")
	fmt.Printf("    quit    closes server\n")
	fmt.Printf("\n")
	var quit bool = false
	for !quit {
		command := readCommand()

		switch command {
		case "quit":
			quit = true
		}
	}
}
