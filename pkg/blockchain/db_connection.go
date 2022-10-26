package blockchain

import "github.com/dgraph-io/badger"

func connectDB() *badger.DB {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}

	return db
}
