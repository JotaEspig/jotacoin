package blockchain

import (
	"bytes"
	"crypto/sha256"
	"jotacoin/pkg/utils"
	"log"
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

// InitData transform the values of the block plus the nonce and difficulty to generate the hash afterwards
func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
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

	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		log.Printf("\r%x", hash)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}
	log.Println()
	return nonce, hash[:]
}

// Validate checks the validation of the block
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// NewProof creates a new Proof of Work struct
func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	return &ProofOfWork{block, target}
}
