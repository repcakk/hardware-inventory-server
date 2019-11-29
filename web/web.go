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
func moveToStorage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gpuSN := r.FormValue("gpuSN")

	fmt.Println("MOVE TO STORAGE CALLED")
	fmt.Println(gpuSN)
	db.MoveGpuToStorage(gpuSN)
}

var serveMux *http.ServeMux
var server http.Server

// Init initialize server
func Init(port string) error {
	serveMux = http.NewServeMux()
	serveMux.HandleFunc("/update", updateHandler)
	serveMux.HandleFunc("/inventory", inventoryHandler)
	serveMux.HandleFunc("/move_to_storage", moveToStorage)
	server = http.Server{Addr: ":" + port, Handler: serveMux}
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
