package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Block represents a block in a blockchain
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// SerializedBlock represents a block that has passes through the process
// of serializing
type SerializedBlock []byte

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

// Deserialize transforms a SerializedBlock into a Block
func (data SerializedBlock) Deserialize() *Block {
	block := &Block{}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(block)
	if err != nil {
		log.Panic(err)
	}

	return block
}

// NewBlock creates a new block struct
func NewBlock(data string, prevHash []byte) *Block {
	b := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(b)
	nonce, hash := pow.Run()

	b.Hash = hash[:]
	b.Nonce = nonce
	return b
}
