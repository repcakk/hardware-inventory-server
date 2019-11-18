package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"

	db "github.com/repcakk/hardware-inventory-server/database"
)

// GpuInventoryStatus represents status
type GpuInventoryStatus struct {
	GpuInUse    []db.UserGpuInfo
	GpuNotInUse []db.GpuInfo
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hostname := r.FormValue("hostname")
	username := r.FormValue("username")
	gpuSN := r.FormValue("gpuSN")
	gpuName := r.FormValue("gpuName")

	db.UserDB.AddOrUpdateRow(hostname, username)
	db.GpuAllDB.AddOrUpdateRow(gpuSN, gpuName)
	db.GpuInUseDB.AddOrUpdateRow(hostname, gpuSN)
}

func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	hardwareInventoryTemplate, err := template.ParseFiles("static/html/hardware-inventory.html")

	if err != nil {
		log.Fatal(err)
	}

	var inventoryStatus GpuInventoryStatus
	inventoryStatus.GpuInUse = db.GetGpuInUse()
	inventoryStatus.GpuNotInUse = db.GetGpuNotInUse()

	err = hardwareInventoryTemplate.Execute(w, inventoryStatus)

	if err != nil {
		log.Fatal(err)
	}
}

var srv http.Server

// Run runs http server and perform http handling
func Run(port string) error {

	m := http.NewServeMux()
	srv = http.Server{Addr: ":" + port, Handler: m}

	m.HandleFunc("/update", updateHandler)
	m.HandleFunc("/inventory", inventoryHandler)

	// TODO: Improve error handling
	log.Fatal(srv.ListenAndServe())

	fmt.Println("AFTER LISTEN AND SERVE")

	return nil
}

// Shutdown server
func Shutdown() {
	fmt.Println("SHUTDOWN CALLED")
	srv.Shutdown(context.Background())
}
