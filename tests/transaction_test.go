package tests

import (
	"jotacoin/pkg/blockchain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	chain, err := blockchain.ContinueBlockchain()
	if err != nil {
		panic(err)
	}
	defer chain.DB.Close()
	tx, err := blockchain.NewTransaction(address1, address2, 1, chain)
	assert.Equal(t, nil, err)

	isValid := tx.Verify()
	assert.Equal(t, true, isValid)
}
