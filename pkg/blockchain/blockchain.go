package blockchain

import (
	"encoding/hex"
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
		cbtx, err := CoinbaseTx(address, genesisData)
		if err != nil {
			return err
		}

		genesis := Genesis(cbtx)
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
func (chain *BlockChain) AddBlock(txs []*Transaction) error {
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

	newBlock := NewBlock(txs, lastHash)
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

func (chain *BlockChain) FindUnspentTransactions(address string) []*Transaction {
	var unspentTxs []*Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()
	for {
		block, err := iter.Next()
		if err != nil {
			break
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
				if !tx.IsCoinbase() {
					for _, txin := range tx.Inputs {
						if txin.CanUnlock(address) {
							txinID := hex.EncodeToString(txin.ID)
							spentTXOs[txinID] = append(spentTXOs[txinID], txin.Out)
						}
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// TODO MUST FINISH: https://youtu.be/HNID8W2jgRM?t=1122
