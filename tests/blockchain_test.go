package tests

import (
	"fmt"
	"jotacoin/pkg/blockchain"
	"jotacoin/pkg/wallet"
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

// Run it after TestAddBlock
func TestGetBalance(t *testing.T) {
	wallets, err := wallet.LoadFile()
	if err != nil {
		panic(err)
	}
	w1 := wallets.GetWallet(address1)
	w2 := wallets.GetWallet(address2)
	pubKeyHash1, err := wallet.PublicKeyHash(w1.PublicKey)
	if err != nil {
		panic(err)
	}
	pubKeyHash2, err := wallet.PublicKeyHash(w2.PublicKey)
	if err != nil {
		panic(err)
	}

	chain, err := blockchain.ContinueBlockChain()
	if err != nil {
		panic(err)
	}
	defer chain.DB.Close()

	balance1 := chain.GetBalance(pubKeyHash1)
	balance2 := chain.GetBalance(pubKeyHash2)
	assert.Equal(t, 90, balance1)
	assert.Equal(t, 10, balance2)
}
