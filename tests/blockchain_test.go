package tests

import (
	"fmt"
	"jotacoin/pkg/blockchain"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	chain, err := blockchain.ContinueBlockChain()
	if err != nil {
		panic(err)
	}
	defer chain.DB.Close()
	tx, err := blockchain.NewTransaction(address1, address2, 10, chain)
	assert.Equal(t, nil, err)
	chain.AddBlock([]*blockchain.Transaction{tx})

	chain.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		if err != nil {
			panic(err)
		}

		return item.Value(func(val []byte) error {
			assert.Equal(t, chain.LastHash, val)
			return nil
		})
	})

	// Go through all the blocks created
	iter := chain.Iterator()
	blocksAmount := 0
	for {
		block, err := iter.Next()
		if err != nil {
			break
		}

		fmt.Printf("Block hash: %x\n\n", block.Hash)
		for _, tx := range block.Transactions {
			fmt.Printf("Transaction Hash: %x\n\n", tx.HashID)

			fmt.Println("INPUTS:")
			for _, in := range tx.Inputs {
				fmt.Printf("PrevTxHash: %x\nOutIdx: %d\nSig: %x\n",
					in.PrevTxHash, in.OutIdx, in.Signature)
			}

			fmt.Println("\nOUTPUTS:")
			for _, out := range tx.Outputs {
				fmt.Printf("Amount: %d\nPubKey: %x\n", out.Value, out.PubKeyHash)
			}
			fmt.Printf("\n==================\n\n")
		}

		blocksAmount++
	}
	assert.GreaterOrEqual(t, blocksAmount, 2)
}
