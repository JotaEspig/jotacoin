package blockchain

import (
	"bytes"
	"crypto/sha256"
	"jotacoin/pkg/utils"
	"math"
	"math/big"
)

// Difficulty defines the difficult of the proof
const Difficulty = 12

// ProofOfWork represents a struct that will be responsable to run the algorithm
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof creates a new Proof of Work struct according to the Difficulty variable value
func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	return &ProofOfWork{block, target}
}

// InitData transform the values of the block plus the nonce and difficulty to generate the hash afterwards
func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			utils.ToHex(int64(nonce)),
			utils.ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
}

// Run runs the proof of work and generates the nonce and the hash
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [sha256.Size]byte
	var nonce int

	for nonce = 0; nonce < math.MaxInt64; nonce++ {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		}
	}
	return nonce, hash[:]
}

// IsValid checks the validation of the block
func (pow *ProofOfWork) IsValid() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
