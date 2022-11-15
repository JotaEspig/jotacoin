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
	block, err := getBlock(iter.DB, iter.CurrentHash)
	if err != nil {
		return nil, err
	}

	iter.CurrentHash = block.PrevHash
	return block, nil
}
