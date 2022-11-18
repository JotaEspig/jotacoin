package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"os"
)

var (
	// WalletFilePath is the file path where the wallets will be stored
	WalletFilePath = "./dbwallets/"
	// WalletFile is the file where the wallets will be stored
	WalletFile = "wallets.data" // TODO make it just store trhe filename
)

// Wallets represents a map containing the wallets
type Wallets map[string]*Wallet

// LoadFile load the content of a file and returns the map containing the maps
func LoadFile() (Wallets, error) {
	var wallets Wallets

	filepath := WalletFilePath + WalletFile
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return Wallets{}, err
	}

	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Wallets{}, err
	}

	// TODO trying to use elliptic.Marshall and Unmarshall
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return Wallets{}, err
	}

	return wallets, nil
}

// SaveFile saves the wallets into a file
func (ws *Wallets) SaveFile() error {
	var content bytes.Buffer

	// TODO trying to use elliptic.Marshall and Unmarshall
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		return err
	}

	if _, err = os.Stat(WalletFilePath); os.IsNotExist(err) {
		os.Mkdir(WalletFilePath, os.ModePerm)
	}

	filepath := WalletFilePath + WalletFile
	return ioutil.WriteFile(filepath, content.Bytes(), 0644)
}

func (ws *Wallets) AddWallet() (string, error) {
	w, err := NewWallet()
	if err != nil {
		return "", err
	}

	address, err := w.Address()
	if err != nil {
		return "", nil
	}

	(*ws)[address] = w

	return address, nil
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range *ws {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) GetWallet(address string) *Wallet {
	return (*ws)[address]
}
