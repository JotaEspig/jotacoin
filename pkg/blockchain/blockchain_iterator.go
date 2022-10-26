package blockchain

import (
	"github.com/dgraph-io/badger"
)

// Iterator is a struct that will iterate with all the hashes
// in the database
type Iterator struct {
	CurrentHash []byte
	DB          *badger.DB
}

// Next returns the block according to Iterator.CurrentHash and will set the
// Iterator.CurrentHash to be the Block.PrevHash gotten
func (iter *Iterator) Next() (*Block, error) {
	block := &Block{}

	err := iter.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			block, err = bytesToSerializedBlock(val).Deserialize()
			return err
		})
	})
	if err != nil {
		return nil, err
	}

	iter.CurrentHash = block.PrevHash
	return block, nil
}
