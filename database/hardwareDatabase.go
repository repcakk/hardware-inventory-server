package database

import (
	"github.com/prologic/bitcask"
)

var hardwareDB, _ = bitcask.Open("data/hardware-database")

// CloseHardwareDB close users database file
func CloseHardwareDB() {
	defer hardwareDB.Close()
}

// AddOrUpdateRow adds or updates new computerID-gpuName pair to database
func AddOrUpdateRow(computerID string, gpuName string) {
	hardwareDB.Put([]byte(computerID), []byte(gpuName))
	hardwareDB.Merge()
}

// GetRows gets all data from database
func GetRows() map[string]string {
	computersGpuMap := make(map[string]string)
	for key := range hardwareDB.Keys() {
		value, _ := hardwareDB.Get([]byte(key))
		computersGpuMap[string([]byte(key))] = string(value)
	}

	userGpuMap := make(map[string]string)
	for computerID, gpuName := range computersGpuMap {
		userGpuMap[userComputerMap[computerID]] = gpuName
	}

	return userGpuMap
}
