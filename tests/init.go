package tests

import (
	"io/ioutil"
	"jotacoin/pkg/database"
	"jotacoin/pkg/wallet"
	"log"
)

func init() {
	log.SetOutput(ioutil.Discard)
	database.DBPath = "./../dbtest/"
	database.DBFile = "./../dbtest/MANIFEST"
	wallet.WalletFilePath = "./../dbtestwallets/"
}
