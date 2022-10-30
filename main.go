package main

import (
	"fmt"
	"jotacoin/pkg/blockchain"
	"jotacoin/pkg/cli"
)

func main() {
	fmt.Println("==== JOTA COIN ====")
	chain, err := blockchain.NewBlockChain("Jota")
	if err != nil {
		panic(err)
	}

	cli := &cli.CommandLine{
		Chain: chain,
	}
	cli.Run()
}
