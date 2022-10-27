package blockchain

import (
	"errors"
	"jotacoin/pkg/database"

	"github.com/dgraph-io/badger"
)

const genesisData = "Genesis Transaction"

// BlockChain Represents a chain of blocks
type BlockChain struct {
	LastHash []byte
	DB       *badger.DB
}

// NewBlockChain creates a new blockchain, starting with 'genesis'
func NewBlockChain(address string) (*BlockChain, error) {
	var lastHash []byte

	if database.DBexists() {
		return nil, errors.New("Blockchain already exists")
	}

	db := database.ConnectDB()
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		if err == badger.ErrKeyNotFound {
			genesis := NewBlock("genesis", []byte{})
			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}

			err = txn.Set([]byte("lastHash"), genesis.Hash)
			if err != nil {
				return err
			}

			lastHash = genesis.Hash
			return nil
		}

		err = txn.Set([]byte("lastHash"), genesis.Hash)
		if err != nil {
			return err
		}

		lastHash = genesis.Hash
		return nil
	})

	return &BlockChain{lastHash, db}, err
}

func ContinueBlockChain(address string) (*BlockChain, error) {
	var lastHash []byte

	if !database.DBexists() {
		return nil, errors.New("Blockchain doesn't exist")
	}

	db := database.ConnectDB()
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
	if err != nil {
		return nil, err
	}

	return &BlockChain{lastHash, db}, nil
}

// Iterator creates a BlockChain Iterador
func (chain *BlockChain) Iterator() *Iterator {
	return &Iterator{chain.LastHash, chain.DB}
}

// AddBlock adds a block into the chain of blocks
func (chain *BlockChain) AddBlock(data string) error {
	var lastHash []byte

	err := chain.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
	})
	if err != nil {
		return err
	}

	newBlock := NewBlock(data, lastHash)
	return chain.DB.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = txn.Set([]byte("lastHash"), newBlock.Hash)
		if err != nil {
			return err
		}

		chain.LastHash = newBlock.Hash
		return nil
	})
}
