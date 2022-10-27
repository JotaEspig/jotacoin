package database

import (
	"os"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./db/"
	dbFile = "./db/MANIFEST"
)

func ConnectDB() *badger.DB {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}

	return db
}

func DBexists() bool {
	_, err := os.Stat(dbFile)
	return !os.IsNotExist(err)
}
