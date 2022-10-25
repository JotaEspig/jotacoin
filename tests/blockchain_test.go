package tests

import (
	"jotacoin/pkg/blockchain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	chain := blockchain.NewBlockChain()
	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")
	assert.Equal(t, 3, len(chain))
	for idx, block := range chain {
		if idx+1 == len(chain) {
			break
		}
		nextBlock := chain[idx+1]
		assert.Equal(t, block.Hash, nextBlock.PrevHash)

		pow := blockchain.NewProof(block)
		assert.Equal(t, true, pow.IsValid())
	}

	block := chain[1]
	block.Data = []byte("cavalo") // changes deliberatily the value of data
	assert.Equal(t, block.Hash, chain[2].PrevHash)
	pow := blockchain.NewProof(block)
	isValid := pow.IsValid()
	assert.NotEqual(t, true, isValid)
}
