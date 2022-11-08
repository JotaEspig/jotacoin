package blockchain

import (
	"bytes"
	"jotacoin/pkg/wallet"

	"github.com/mr-tron/base58"
)

// TxInput represents an input of a transaction. For more information:
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch05.html
type TxInput struct {
	PrevTxHash []byte // previous transacion ID, where the output is stored
	OutIdx     int    // idx of output in the transaction struct
	Signature  []byte
	PubKey     []byte
}

// TxOutput represents an output of a transaction. For more information:
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch05.html
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// UsesKey checks if the hash of TxInput.PubKey is the same as the input
func (txin *TxInput) UsesKey(publicKeyHash []byte) bool {
	lockedHash, err := wallet.PublicKeyHash(txin.PubKey)
	if err != nil {
		return false
	}
	return bytes.Compare(lockedHash, publicKeyHash) == 0
}

// NewTxOutput creates a new output
func NewTxOutput(value int, address string) (*TxOutput, error) {
	txout := &TxOutput{value, nil}
	err := txout.Lock(address)
	return txout, err
}

// Lock locks the output according to the address
func (txout *TxOutput) Lock(address string) error {
	fullHash, err := base58.Decode(address)
	if err != nil {
		return err
	}

	txout.PubKeyHash = fullHash[1 : len(fullHash)-wallet.ChecksumLength]
	return nil
}

// IsLockedWithKey checks if the output is locked with the key passed in the args
func (txout *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txout.PubKeyHash, pubKeyHash) == 0
}
