package blockchain

import "github.com/dgraph-io/badger"

func getLastHash(db *badger.DB) ([]byte, error) {
	var lastHash []byte

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
	})
	return lastHash, err
}

func getBlock(db *badger.DB, hash []byte) (*Block, error) {
	var block *Block

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			block, err = bytesToSerializedBlock(val).Deserialize()
			return err
		})
	})

	return block, err
}

func addBlockToDB(db *badger.DB, b *Block) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Set(b.Hash, b.Serialize())
		if err != nil {
			return err
		}
		err = txn.Set([]byte("lastHash"), b.Hash)
		if err != nil {
			return err
		}

		return nil
	})
}
