package database

import (
	"fmt"
	"strconv"
)

// GpuInfo represents name and serial number of single GPU
type GpuInfo struct {
	GpuSN   string
	GpuName string
}

// UserInfo represents username and hostname of single user
type UserInfo struct {
	Hostname string
	Username string
}

// UserGpuInfo represent usage of GPU by user
type UserGpuInfo struct {
	User UserInfo
	Gpu  GpuInfo
}

// GetGpuInUse returns array of all GPUs currently in use
func GetGpuInUse() []UserGpuInfo {
	usernames := UserDB.GetRows()
	gpuInfo := GpuAllDB.GetRows()
	gpuInUse := GpuInUseDB.GetRows()

	gpuStatus := make([]UserGpuInfo, len(usernames))

	userCount := 0
	for key, value := range usernames {
		gpuStatus[userCount].User.Hostname = key
		gpuStatus[userCount].User.Username = value
		gpuStatus[userCount].Gpu.GpuSN = gpuInUse[key]                             // Access GPU serial number by hostname
		gpuStatus[userCount].Gpu.GpuName = gpuInfo[gpuStatus[userCount].Gpu.GpuSN] // Access GPU name by gpu serial number
		userCount++
	}
	return gpuStatus
}

// GetGpuNotInUse returns array of all GPUs currently NOT in use
func GetGpuNotInUse() []GpuInfo {
	gpuAll := GpuAllDB.GetRows()
	gpuInUse := GpuInUseDB.GetRows()

	// reversedGpuInUse it maps GPU serial numbers to hostnames
	reversedGpuInUse := make(map[string]string)
	for key, value := range gpuInUse {
		reversedGpuInUse[value] = key
	}

	gpuStatus := make([]GpuInfo, len(gpuAll))

	fmt.Println("gpuAll LEN: " + strconv.Itoa(len(gpuAll)))
	fmt.Println("gpuInUse LEN: " + strconv.Itoa(len(gpuInUse)))
	fmt.Println("ARRAY LEN: " + strconv.Itoa(len(gpuAll)-len(gpuInUse)))
	gpuCount := 0
	for key, value := range gpuAll {
		// if gpu serial number does not have assigned hostname it means GPU is not in use
		if reversedGpuInUse[key] == "" {
			gpuStatus[gpuCount].GpuSN = key
			gpuStatus[gpuCount].GpuName = value
			gpuCount++
		}
	}
	return gpuStatus
}

// MoveGpuToStorage moves gpu to storage for given gpu serial number
func MoveGpuToStorage(gpuSN string) {
	gpuInUse := GpuInUseDB.GetRows()

	// reversedGpuInUse it maps GPU serial numbers to hostnames
	reversedGpuInUse := make(map[string]string)
	for key, value := range gpuInUse {
		reversedGpuInUse[value] = key
	}

	GpuInUseDB.Delete([]byte(reversedGpuInUse[gpuSN]))
}
