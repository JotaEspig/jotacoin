package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
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

// HashTransactions generates the hash of the combined transactions
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [sha256.Size]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Hash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// Serialize returns a block struct as a SerializedBlock
func (b *Block) Serialize() SerializedBlock {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}
