package wallet

import (
	"bytes"
	"crypto/x509"
	"encoding/gob"
	"io/ioutil"
	"jotacoin/pkg/utils"
	"os"
)

var (
	// WalletFilePath is the file path where the wallets will be stored
	WalletFilePath = "./dbwallets/"
	// WalletFile is the file where the wallets will be stored
	WalletFile = "wallets.data"
)

// Wallets represents a map containing the wallets
type Wallets map[string]*Wallet

// walletsToFile is a struct that is the midway between Wallets struct and the file content
type walletAsBytes struct {
	PrivateKey []byte
	PublicKey  []byte
}

// LoadFile load the content of a file and returns the map containing the maps
func LoadFile() (Wallets, error) {
	var wsToLoad []walletAsBytes
	wallets := Wallets{}

	filepath := WalletFilePath + WalletFile
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return Wallets{}, err
	}

	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Wallets{}, err
	}

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wsToLoad)
	if err != nil {
		return Wallets{}, err
	}

	for _, wf := range wsToLoad {
		priv, err := x509.ParseECPrivateKey(wf.PrivateKey)
		if err != nil {
			return Wallets{}, err
		}

		w := &Wallet{priv, wf.PublicKey}
		address, err := w.Address()
		if err != nil {
			return Wallets{}, err
		}

		wallets[address] = w
	}

	return wallets, nil
}

// SaveFile saves the wallets into a file
func (ws *Wallets) SaveFile() error {
	wsToSave := []walletAsBytes{}

	for _, w := range *ws {
		priv, err := x509.MarshalECPrivateKey(w.PrivateKey)
		if err != nil {
			return err
		}

		wsToSave = append(wsToSave, walletAsBytes{
			priv,
			w.PublicKey,
		})
	}

	content, err := utils.Serialize(wsToSave)
	if err != nil {
		return err
	}

	if _, err = os.Stat(WalletFilePath); os.IsNotExist(err) {
		os.Mkdir(WalletFilePath, os.ModePerm)
	}

	filepath := WalletFilePath + WalletFile
	return ioutil.WriteFile(filepath, content, 0644)
}

// AddWallet adds a wallet to the Wallets map (itself)
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

// GetAllAddresses returns the addresses of all wallets stored in Wallets map
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range *ws {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet gets a wallet from the Wallets map according to the address
func (ws *Wallets) GetWallet(address string) *Wallet {
	return (*ws)[address]
}
