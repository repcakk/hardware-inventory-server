package database

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // blank import for adding postgres support to gorm
)

// databaseConfig structure for storing database properties
type databaseConfig struct {
	DbHost         string `json:"dbHost"`
	DbPort         string `json:"dbPort"`
	DbName         string `json:"dbName"`
	DbUser         string `json:"dbUser"`
	DbUserPassword string `json:"dbUserPassword"`
}

// Gpu represents name and serial number of single GPU.
// SN				GPU Serial Number, used as unique ID.
// Name				GPU Name.
type Gpu struct {
	gorm.Model
	SN      string
	GpuName string
}

// User represents single user.
// Name				Name of User.
// Surname			Surname of User.
// Email			Unique user email.
type User struct {
	gorm.Model
	Username string
	Surname  string
	Email    string
}

// Computer represents current status of single computer.
// MacAddress		MAC address of computer, used as unique ID.
// Hostname			Name of computer.
// CurrentGpuID		ID of GPU  currently being used in this computer.
// LastGpuID		ID of GPU previusly used in this computer.
// UserID			ID of current user of this computer.
type Computer struct {
	gorm.Model
	MacAddress   string
	Hostname     string
	CurrentGpuID uint
	LastGpuID    uint
	UserID       uint
}

// ComputerInfo represents computer, current gpu last gpu and user.
type ComputerInfo struct {
	Computer   Computer
	CurrentGpu Gpu
	LastGpu    Gpu
	User       User
}

// gorm library database driver used for ORM and SQL inserts.
var gormDB *gorm.DB

var dbConfig databaseConfig = loadConfig("./config/database-config.json")

// loadConfig loads database config
func loadConfig(file string) databaseConfig {
	var config databaseConfig
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return config
}

// ConnectGormDB connects to database.
func ConnectGormDB() error {
	var err error
	connectionString := "host=" + dbConfig.DbHost +
		" port=" + dbConfig.DbPort +
		" dbname=" + dbConfig.DbName +
		" user=" + dbConfig.DbUser +
		" password=" + dbConfig.DbUserPassword +
		" sslmode=disable"
	gormDB, err = gorm.Open("postgres", connectionString)

	return err
}

// CloseGormDB closes connection.
func CloseGormDB() error {
	return gormDB.Close()
}

// MigrateSchema creates schema according to provided structures if table not exists.
func MigrateSchema() {
	gormDB.AutoMigrate(&Gpu{}, &User{}, &Computer{})

	// gormDB.Model(&Gpu{}).AddUniqueIndex("idx_sn", "sn")
	// gormDB.Model(&User{}).AddUniqueIndex("idx_email", "email")
	// gormDB.Model(&Computer{}).AddUniqueIndex("idx_mac_address", "mac_address")
	gormDB.Model(&Computer{}).AddForeignKey("current_gpu_id", "gpus(id)", "NO ACTION", "NO ACTION")
	gormDB.Model(&Computer{}).AddForeignKey("last_gpu_id", "gpus(id)", "NO ACTION", "NO ACTION")
	gormDB.Model(&Computer{}).AddForeignKey("user_id", "users(id)", "NO ACTION", "NO ACTION")
}

// AddUser adds user to database.
// returns true if new record was added, false if not.
func AddUser(user User) (User, bool) {
	previousValue := user.Email
	gormDB.FirstOrCreate(&user, User{Email: user.Email})
	return user, previousValue != user.Email
}

// AddGpu adds gpu to database.
// returns true if new record was added, false if not.
func AddGpu(gpu Gpu) (Gpu, bool) {
	previousValue := gpu.SN
	gormDB.FirstOrCreate(&gpu, Gpu{SN: gpu.SN})
	return gpu, previousValue != gpu.SN
}

// AddComputer adds computer to database.
// returns true if new record was added, false if not.
func AddComputer(computer Computer) (Computer, bool) {
	previousValue := computer.MacAddress
	gormDB.FirstOrCreate(&computer, Computer{MacAddress: computer.MacAddress})
	return computer, previousValue != computer.MacAddress
}

// ChangeGpuInComputer updates currently GPU and last GPU for given computer.
func ChangeGpuInComputer(computerID uint, gpuID uint) {
	// If there is computer with currentGpu == gpuID remove GPU from this computer first.
	var computerWithGPU Computer
	gormDB.Where("current_gpu_id = ?", gpuID).First(&computerWithGPU)
	if computerWithGPU.CurrentGpuID == gpuID {
		fmt.Println("INSIDE IF")
		gormDB.Model(&computerWithGPU).Update("current_gpu_id", gorm.Expr("NULL"), "last_gpu_id", computerWithGPU.CurrentGpuID)
	}

	// Then update GPU in selected computer
	var computer Computer
	gormDB.Where("id = ?", computerID).First(&computer)
	if computer.CurrentGpuID != gpuID {
		gormDB.Model(&computer).Update(Computer{CurrentGpuID: gpuID, LastGpuID: computer.CurrentGpuID})
	}
}

// GetComputersGpusUsers return all computers with current gpu info and last gpu info.
// offset 			how many rows should be skipped.
// limit			how many rows should be fetched.
// return			struct of computers, current gpus, last gpus and users.
func GetComputersGpusUsers(offset uint, limit uint) []ComputerInfo {
	computersInfo := make([]ComputerInfo, 0)
	queryString := "SELECT c.*, g1.*, g2.*, u.* FROM computers AS c " +
		"LEFT JOIN gpus AS g1 ON c.current_gpu_id = g1.id " +
		"LEFT JOIN gpus AS g2 ON c.last_gpu_id = g2.id " +
		"LEFT JOIN users AS u ON c.user_id = u.id"

	rows, _ := gormDB.Raw(queryString).Rows()
	defer rows.Close()

	for rows.Next() {
		var computerInfo ComputerInfo

		gormDB.ScanRows(rows, &computerInfo.Computer)
		gormDB.ScanRows(rows, &computerInfo.CurrentGpu)
		gormDB.ScanRows(rows, &computerInfo.LastGpu)
		gormDB.ScanRows(rows, &computerInfo.User)

		fmt.Println(computerInfo.Computer)
		fmt.Println(computerInfo.CurrentGpu)
		fmt.Println(computerInfo.LastGpu)
		fmt.Println(computerInfo.User)
		fmt.Println("")
		fmt.Println("")

		computersInfo = append(computersInfo, computerInfo)
	}

	return computersInfo
}
