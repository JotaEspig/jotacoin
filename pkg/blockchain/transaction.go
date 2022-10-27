package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const CoinbaseValue = 100

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID  []byte // transacion ID
	Out int    // idx of output in the transaction struct
	Sig string
}

type TxOutput struct {
	Value  int
	PubKey string
}

func (tx *Transaction) SetID() error {
	var result bytes.Buffer
	var hash [sha256.Size]byte

	encode := gob.NewEncoder(&result)
	err := encode.Encode(tx)
	if err != nil {
		return err
	}

	hash = sha256.Sum256(result.Bytes())
	tx.ID = hash[:]
	return nil
}

func CoinbaseTx(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{CoinbaseValue, to}

	tx := &Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	err := tx.SetID()

	return tx, err
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 1 && tx.Inputs[0].Out == -1
}

func (txin *TxInput) CanUnlock(data string) bool {
	return txin.Sig == data
}

func (txout *TxOutput) CanBeUnlocked(data string) bool {
	return txout.PubKey == data
}
