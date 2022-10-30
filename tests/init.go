package tests

import (
	"io/ioutil"
	"jotacoin/pkg/database"
	"log"
)

func init() {
	log.SetOutput(ioutil.Discard)
	database.DBPath = "./../dbtest/"
	database.DBFile = "./../dbtest/MANIFEST"
}
