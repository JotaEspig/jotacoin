package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
)

// CoinbaseValue is the predefined amount of tokens for a coinbase
const CoinbaseValue = 100

// Transaction represents a transaction in a blockchain. For more information:
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch05.html
type Transaction struct {
	Hash    []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// NewTransaction creates a normal transaction (one sender and one receiver)
func NewTransaction(from, to string, amount int, chain *BlockChain) (*Transaction, error) {
	var inputs []TxInput
	var outputs []TxOutput

	acc, spendableTxs := chain.FindSpendableTx(from, amount)
	if acc < amount {
		return nil, errors.New("transaction: not enough balance from the sender")
	}

	for prevTxIDstr, outsIdxs := range spendableTxs {
		prevTxID, err := hex.DecodeString(prevTxIDstr)
		if err != nil {
			return nil, err
		}

		for _, outIdx := range outsIdxs {
			input := TxInput{prevTxID, outIdx, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		// if the accumulated is greater than the payment, there should be a change
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := &Transaction{nil, inputs, outputs}
	err := tx.SetHash()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// NewCoinbaseTx creates a coinbase and it "gives" to a receiver
func NewCoinbaseTx(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{CoinbaseValue, to}

	tx := &Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	err := tx.SetHash()

	return tx, err
}

// SetHash sets the hash to the transaction
func (tx *Transaction) SetHash() error {
	var result bytes.Buffer
	var hash [sha256.Size]byte

	encode := gob.NewEncoder(&result)
	err := encode.Encode(tx)
	if err != nil {
		return err
	}

	hash = sha256.Sum256(result.Bytes())
	tx.Hash = hash[:]
	return nil
}

// IsCoinbase checks if the transaction is a coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].PrevTxHash) == 0 && tx.Inputs[0].OutIdx == -1
}
