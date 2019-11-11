package database

import (
	"github.com/prologic/bitcask"
	helpers "github.com/repcakk/hardware-inventory-server/helpers"
)

// All database objects are key-value maps.

// GpuAllDB stores gpu serial number and GPU name pair.
// It is database of all available GPUs.
var GpuAllDB, _ = OpenDB("data/gpu-info-database")

// GpuInUseDB stores hostname and gpu serial number pair.
// It is database of all currently used GPUs in computers.
var GpuInUseDB, _ = OpenDB("data/gpu-in-use-database")

// UserDB stores hostname and user pair
// It is database of all users and theirs computers
var UserDB, _ = OpenDB("data/user-database")

//BitcaskDB wraps bitcask into structure to extend its functionalities
type BitcaskDB struct {
	*bitcask.Bitcask
}

// TODO: remove this object. It is stored in userDB
var userComputerMap = helpers.ReadJSON("config/usernames-mapping.json")

// OpenDB creates and opens database for given database path
// dbPath - target database path
func OpenDB(dbPath string) (BitcaskDB, error) {

	db, err := bitcask.Open(dbPath)

	dbWrapper := BitcaskDB{db}

	return dbWrapper, err
}

// CloseDB close database file
// db - target database
func (db *BitcaskDB) CloseDB() {
	defer db.Close()
}

// AddOrUpdateRow adds new or updates row in database
// db - target database
func (db *BitcaskDB) AddOrUpdateRow(key string, value string) {
	db.Put([]byte(key), []byte(value))
	db.Merge()
}

// GetRows gets all data from given database
// db - target database
func (db *BitcaskDB) GetRows() map[string]string {
	resultsMap := make(map[string]string)
	for key := range db.Keys() {
		value, _ := db.Get([]byte(key))
		resultsMap[string([]byte(key))] = string(value)
	}

	return resultsMap
}
