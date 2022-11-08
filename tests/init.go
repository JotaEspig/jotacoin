package tests

import (
	"io/ioutil"
	"jotacoin/pkg/blockchain"
	"jotacoin/pkg/database"
	"jotacoin/pkg/wallet"
	"log"
)

var (
	address1 string
	address2 string
)

func init() {
	var err error

	log.SetOutput(ioutil.Discard)
	database.DBPath = "./../dbtest/"
	database.DBFile = "./../dbtest/MANIFEST"
	wallet.WalletFilePath = "./../dbtestwallets/"

	// sets the addresses for test and creates the blockchain
	ws := wallet.Wallets{}
	address1, err = ws.AddWallet()
	if err != nil {
		panic(err)
	}
	address2, err = ws.AddWallet()
	if err != nil {
		panic(err)
	}
	err = ws.SaveFile()
	if err != nil {
		panic(err)
	}
	chain, err := blockchain.NewBlockChain(address1)
	if err != nil {
		panic(err)
	}
	defer chain.DB.Close()
}
