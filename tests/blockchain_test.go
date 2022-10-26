package tests

import (
	"fmt"
	"jotacoin/pkg/blockchain"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	chain, err := blockchain.NewBlockChain()
	if err != nil {
		panic(err)
	}
	defer chain.DB.Close()
	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")

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

	iter := chain.Iterator()
	blocksAmount := 0
	for {
		block, err := iter.Next()
		if err != nil {
			break
		}
		fmt.Printf("Hash: %x\nValue: %s\nPrevious Hash: %x\n\n",
			block.Hash, string(block.Data), block.PrevHash)

		blocksAmount++
	}
	assert.GreaterOrEqual(t, blocksAmount, 3)
}
