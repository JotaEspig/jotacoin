package cli

import (
	"fmt"
	"jotacoin/pkg/blockchain"
	"jotacoin/pkg/wallet"
	"os"
)

// CommandLine is the struct that is responsable for running the commands
type CommandLine struct {
}

func (cli *CommandLine) newWallet() {
	w, err := wallet.NewWallet()
	if err != nil {
		panic(err)
	}
	address, err := w.Address()
	if err != nil {
		panic(err)
	}

	ws := wallet.Wallets{address: w}
	err = ws.SaveFile()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Added Wallet!\nAddress: %s\n", address)
}

func (cli *CommandLine) showWallets() {
	ws, err := wallet.LoadFile()
	if err != nil {
		panic(err)
	}
	for _, w := range ws {
		address, err := w.Address()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Priv: %v\nPub: %x\nAddress: %s\n\n", w.PrivateKey, w.PublicKey, address)
	}
}

func (cli *CommandLine) newBlockchain(address string) {
	_, err := blockchain.NewBlockchain(address)
	if err != nil {
		panic(err)
	}

	fmt.Println("New BlockChain created")
}

func (cli *CommandLine) printAll() {
	var block *blockchain.Block
	var err error

	chain, err := blockchain.ContinueBlockchain()
	if err != nil {
		panic(err)
	}

	iter := chain.Iterator()
	for {
		block, err = iter.Next()
		if err != nil {
			break
		}

		pow := blockchain.NewProof(block)
		isValid := pow.IsValid()

		fmt.Println()
		fmt.Println(*block, isValid)
	}
}

// Run runs the command line
func (cli *CommandLine) Run() {
	switch os.Args[1] {
	case "newwallet":
		cli.newWallet()
	case "showwallets":
		cli.showWallets()
	case "newblockchain":
		cli.newBlockchain(os.Args[2])
	case "print":
		cli.printAll()
	}
}
