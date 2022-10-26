package cli

import (
	"fmt"
	"jotacoin/pkg/blockchain"
	"log"
	"os"
)

// CommandLine is the struct that is responsable for running the commands
type CommandLine struct {
	Chain *blockchain.BlockChain
}

func (cli *CommandLine) addBlock(data string) {
	err := cli.Chain.AddBlock(data)
	if err != nil {
		log.Println("An error has occured when adding a block: ", err)
		return
	}

	log.Printf("Block added into the database\n")
}

func (cli *CommandLine) printAll() {
	var block *blockchain.Block
	var err error

	iter := cli.Chain.Iterator()
	for {
		block, err = iter.Next()
		if err != nil {
			break
		}

		pow := blockchain.NewProof(block)
		isValid := pow.IsValid()

		fmt.Println()
		fmt.Printf("Hash: %x\nValue: %s\nPrevious Hash: %x\nPoW: %v\n",
			block.Hash, string(block.Data), block.PrevHash, isValid)
	}
}

// Run runs the command line
func (cli *CommandLine) Run() {
	switch os.Args[1] {
	case "add":
		cli.addBlock(os.Args[2])
	case "print":
		cli.printAll()
	}
}
