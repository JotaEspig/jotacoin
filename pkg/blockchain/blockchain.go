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

// NewBlockChain creates a new blockchain, starting with coinbase
func NewBlockChain(address string) (*BlockChain, error) {
	var lastHash []byte

	if database.DBexists() {
		return nil, errors.New("Blockchain already exists")
	}

	db := database.ConnectDB()
	err := db.Update(func(txn *badger.Txn) error {
		cbtx, err := NewCoinbaseTx(address, genesisData)
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

// ContinueBlockChain continues the previous BlockChain if it already exists
func ContinueBlockChain() (*BlockChain, error) {
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

// FindUnspentTransactions returns the Transactions where the output hasn't been spent yet
// by the address
func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []*Transaction {
	var unspentTxs []*Transaction
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()
	for {
		block, err := iter.Next()
		if err != nil {
			break
		}

	Transactions:
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				// If it's not a Coinbase, it means that's a normal transaction
				// so it has a previous transaction
				for _, txin := range tx.Inputs {
					if txin.UsesKey(pubKeyHash) {
						// if address can unlock, it means that address spent the previous output
						prevTxHash := hex.EncodeToString(txin.PrevTxHash)
						spentTXOs[prevTxHash] = append(spentTXOs[prevTxHash], txin.OutIdx)
					}
				}
			}

			txHash := hex.EncodeToString(tx.HashID)
		Outputs:
			for outIdx, out := range tx.Outputs {
				if !out.IsLockedWithKey(pubKeyHash) {
					continue
				}
				// Check if the address has already spent the output
				if spentTXOs[txHash] != nil {
					for _, spentOutIdx := range spentTXOs[txHash] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				unspentTxs = append(unspentTxs, tx)
				// if this transaction was already spent, it can go to another Tx
				continue Transactions
			}
		}

		// if it's the genesis block, break
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

// FindSpendableOutputs returns the tokens accumulated by the spendable outputs and a map where
// the keys are hte Transactions IDs and the values are slices containing the indexes
// of the outputs of that Transaction
func (chain *BlockChain) FindSpendableOutputs(pubKeyHash []byte, requiredAmount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txHash := hex.EncodeToString(tx.HashID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				accumulated += out.Value
				unspentOuts[txHash] = append(unspentOuts[txHash], outIdx)

				if accumulated >= requiredAmount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

// FindUTXO find the unspent outputs of the address. This function is useful
// to get the address balance
func (chain *BlockChain) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	for _, tx := range unspentTxs {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// GetBalance returns the balance of the address
func (chain *BlockChain) GetBalance(pubKeyHash []byte) int {
	unspentOutput := chain.FindUTXO(pubKeyHash)
	total := 0

	for _, out := range unspentOutput {
		total += out.Value
	}

	return total
}
