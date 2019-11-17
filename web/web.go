package web

import (
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

// Run runs http server and perform http handling
func Run() error {
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/inventory", inventoryHandler)

	// TODO: Add config handling
	port := "8080" //config.Port
	// TODO: Improve error handling
	log.Fatal(http.ListenAndServe(":"+port, nil))

	return nil
}
