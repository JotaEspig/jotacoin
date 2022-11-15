package database

import (
	"os"

	"github.com/dgraph-io/badger"
)

var (
	// DBPath is the path to the database folder
	DBPath = "./db/"
	DBFile = "./db/MANIFEST"
)

// ConnectDB connects to the database
func ConnectDB(path string) *badger.DB {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}

	return db
}

// DBExists checks if the database already exists
func DBexists() bool {
	_, err := os.Stat(DBFile)
	return !os.IsNotExist(err)
}
