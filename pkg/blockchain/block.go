package blockchain

// Block represents a block in a blockchain
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
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
