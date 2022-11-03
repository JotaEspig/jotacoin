package tests

import (
	"fmt"
	"jotacoin/pkg/wallet"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddWallet(t *testing.T) {
	ws := wallet.Wallets{}

	fmt.Println(ws.AddWallet())

	err := ws.SaveFile()
	assert.Equal(t, nil, err)

	wsFromLoadFile, err := wallet.LoadFile()
	assert.Equal(t, nil, err)

	addresses := wsFromLoadFile.GetAllAddresses()
	fmt.Println(addresses)
	assert.Equal(t, 1, len(addresses))

	filepath := wallet.WalletFilePath + wallet.WalletFile
	err = os.Remove(filepath)
	assert.Equal(t, nil, err)
}
