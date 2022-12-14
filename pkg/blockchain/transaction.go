package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"jotacoin/pkg/utils"
	"jotacoin/pkg/wallet"
	"math/big"
)

// CoinbaseValue is the predefined amount of tokens for a coinbase
const CoinbaseValue = 100

// Transaction represents a transaction in a blockchain. For more information:
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch05.html
type Transaction struct {
	HashID  []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// NewCoinbaseTx creates a coinbase and it "gives" to a receiver
func NewCoinbaseTx(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout, err := NewTxOutput(CoinbaseValue, to)
	if err != nil {
		return nil, err
	}

	tx := &Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	hash, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	tx.HashID = hash
	return tx, err
}

// IsCoinbase checks if the transaction is a coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].PrevTxHash) == 0 && tx.Inputs[0].OutIdx == -1
}

// NewTransaction creates a normal transaction (one sender and one receiver)
func NewTransaction(from, to string, amount int, chain *Blockchain) (*Transaction, error) {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.LoadFile()
	if err != nil {
		return nil, err
	}
	w := wallets.GetWallet(from)
	if w == nil {
		return nil, errors.New("wallet: wallet not found")
	}
	pubKeyHash, err := wallet.PublicKeyHash(w.PublicKey)
	if err != nil {
		return nil, err
	}

	acc, spendableTxs := chain.FindSpendableTxOutputs(pubKeyHash, amount)
	if acc < amount {
		return nil, errors.New("transaction: not enough balance from the sender")
	}

	for prevTxIDstr, outsIdxs := range spendableTxs {
		prevTxID, err := hex.DecodeString(prevTxIDstr)
		if err != nil {
			return nil, err
		}

		for _, outIdx := range outsIdxs {
			input := TxInput{prevTxID, outIdx, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	newOutput, err := NewTxOutput(amount, to)
	if err != nil {
		return nil, err
	}
	outputs = append(outputs, *newOutput)
	if acc > amount {
		// if the accumulated is greater than the payment, there should be a change
		newOutput, err = NewTxOutput(acc-amount, from)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, *newOutput)
	}

	tx := &Transaction{nil, inputs, outputs}
	err = tx.Sign(w.PrivateKey)
	if err != nil {
		return nil, err
	}
	hash, err := tx.Hash()
	if err != nil {
		return nil, err
	}

	tx.HashID = hash
	return tx, nil
}

// Hash generates and returns the hash of the transaction disregarding the value
// setted for the hash field.
func (tx *Transaction) Hash() ([]byte, error) {
	var hash [sha256.Size]byte

	txCopy := *tx
	txCopy.HashID = []byte{}

	txSerialized, err := utils.Serialize(txCopy)
	if err != nil {
		return []byte{}, err
	}
	hash = sha256.Sum256(txSerialized)

	return hash[:], nil
}

// Sign signs the transaction
func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey) error {
	if tx.IsCoinbase() {
		return nil
	}

	txCopy, err := tx.TrimmedCopy()
	if err != nil {
		return err
	}

	for txinIdx := range txCopy.Inputs {
		signature, err := ecdsa.SignASN1(rand.Reader, privKey, txCopy.HashID)
		if err != nil {
			return err
		}
		tx.Inputs[txinIdx].Signature = signature
	}

	return nil
}

// Verify verifies if the transaction is valid according to the signature and public key
func (tx *Transaction) Verify() bool {
	if tx.IsCoinbase() {
		return true
	}

	curve := elliptic.P256()
	txCopy, err := tx.TrimmedCopy()
	if err != nil {
		return false
	}

	for _, txin := range tx.Inputs {
		var x, y big.Int
		keyLen := len(txin.PubKey)
		x.SetBytes(txin.PubKey[:(keyLen / 2)])
		y.SetBytes(txin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.VerifyASN1(&rawPubKey, txCopy.HashID, txin.Signature) {
			return false
		}
	}

	return true
}

// TrimmedCopy returns a copy of the transaction but without the signature
func (tx *Transaction) TrimmedCopy() (Transaction, error) {
	var txInputs []TxInput
	var txOutputs []TxOutput

	for _, txin := range tx.Inputs {
		txInputs = append(txInputs, TxInput{txin.PrevTxHash, txin.OutIdx, nil, txin.PubKey})
	}
	for _, txout := range tx.Outputs {
		txOutputs = append(txOutputs, TxOutput{txout.Value, txout.PubKeyHash})
	}

	txCopy := Transaction{nil, txInputs, txOutputs}
	hash, err := txCopy.Hash()
	if err != nil {
		return Transaction{}, err
	}
	txCopy.HashID = hash
	return txCopy, nil
}
