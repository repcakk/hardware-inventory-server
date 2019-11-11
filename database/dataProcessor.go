package database

// GpuInfo represents name and serial number of single GPU
type GpuInfo struct {
	gpuName string
	gpuSN   string
}

// UserInfo represents username and hostname of single user
type UserInfo struct {
	hostname string
	username string
}

// UserGpuInfo represent usage of GPU by user
type UserGpuInfo struct {
	userInfo UserInfo
	gpuInfo  GpuInfo
}

// GetGpuInUse returns array of all GPUs currently in use
func GetGpuInUse() []UserGpuInfo {
	usernames := UserDB.GetRows()
	gpuInfo := GpuAllDB.GetRows()
	gpuInUse := GpuInUseDB.GetRows()

	gpuStatus := make([]UserGpuInfo, len(usernames))

	userCount := 0
	for key, value := range usernames {
		gpuStatus[userCount].userInfo.hostname = key
		gpuStatus[userCount].userInfo.username = value
		gpuStatus[userCount].gpuInfo.gpuSN = gpuInUse[key]                                 // Access GPU serial number by hostname
		gpuStatus[userCount].gpuInfo.gpuName = gpuInfo[gpuStatus[userCount].gpuInfo.gpuSN] // Access GPU name by gpu serial number
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

	gpuStatus := make([]GpuInfo, len(gpuAll)-len(gpuInUse))

	gpuCount := 0
	for key, value := range gpuAll {
		// if gpu serial number does not have assigned hostname it means GPU is not in use
		if reversedGpuInUse[key] == "" {
			gpuStatus[gpuCount].gpuSN = key
			gpuStatus[gpuCount].gpuName = value
			gpuCount++
		}
	}
	return gpuStatus
}
