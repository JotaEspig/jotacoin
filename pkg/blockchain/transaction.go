package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
)

const CoinbaseValue = 100

type Transaction struct {
	Hash    []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	PrevTxHash []byte // previous transacion ID, where the output ("balance") is stored
	OutIdx     int    // idx of output in the transaction struct
	Sig        string
}

type TxOutput struct {
	Value  int
	PubKey string
}

func NewTransaction(from, to string, amount int, chain *BlockChain) (*Transaction, error) {
	var inputs []TxInput
	var outputs []TxOutput

	acc, spendableTxs := chain.FindSpendableTx(from, amount)
	if acc < amount {
		return nil, errors.New("transaction: not enough balance")
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

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].PrevTxHash) == 0 && tx.Inputs[0].OutIdx == -1
}

// IsMadeBy checks if the address has made the input
func (txin *TxInput) IsMadeBy(address string) bool {
	return txin.Sig == address
}

// IsFor checks if the address is the receiver of the output
func (txout *TxOutput) IsFor(address string) bool {
	return txout.PubKey == address
}
