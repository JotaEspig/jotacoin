package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Block represents a block in a blockchain
type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

// NewBlock creates a new block struct
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProof(b)
	nonce, hash := pow.Run()

	b.Hash = hash[:]
	b.Nonce = nonce
	return b
}

// Genesis creates a genesis block
func Genesis(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// DeserializeBlock transforms a serialized block ([]byte) into a Block
func DeserializeBlock(data []byte) (*Block, error) {
	block := &Block{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(block)
	return block, err
}

// HashTransactions generates the hash of the combined transactions
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [sha256.Size]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.HashID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
