package cli

import (
	"errors"
	"fmt"
	"jotacoin/pkg/blockchain"
	"jotacoin/pkg/wallet"
	"os"
	"strconv"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// CommandLine is the struct that is responsable for running the commands
type CommandLine struct {
}

func (cli *CommandLine) newWallet() {
	ws, err := wallet.LoadFile()
	if err != nil {
		ws = wallet.Wallets{}
	}
	w, err := wallet.NewWallet()
	handleError(err)
	address, err := w.Address()
	handleError(err)

	ws[address] = w
	err = ws.SaveFile()
	handleError(err)

	fmt.Printf("Added Wallet!\nAddress: %s\n", address)
}

func (cli *CommandLine) showWallets() {
	ws, err := wallet.LoadFile()
	handleError(err)
	for _, w := range ws {
		address, err := w.Address()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Priv: %v\nPub: %x\nAddress: %s\n\n", w.PrivateKey, w.PublicKey, address)
	}
}

func (cli *CommandLine) getBalance(address string) {
	chain, err := blockchain.ContinueBlockchain()
	handleError(err)
	ws, err := wallet.LoadFile()
	handleError(err)
	w := ws.GetWallet(address)
	if w == nil {
		panic(errors.New("wallet does not exists"))
	}
	pubHash, err := wallet.PublicKeyHash(w.PublicKey)
	handleError(err)
	balance := chain.GetBalance(pubHash)
	fmt.Printf("Balance: %d\n", balance)
}

func (cli *CommandLine) newTransaction(from, to string, amount int) {
	chain, err := blockchain.ContinueBlockchain()
	handleError(err)

	tx, err := blockchain.NewTransaction(from, to, amount, chain)
	handleError(err)
	err = chain.AddBlock([]*blockchain.Transaction{tx})
	handleError(err)

	fmt.Printf("Transaction done!\nTx Hash: %x\nInputs: %v\nOutputs: %v\n\n",
		tx.HashID, tx.Inputs, tx.Outputs)
}

func (cli *CommandLine) newBlockchain(address string) {
	_, err := blockchain.NewBlockchain(address)
	handleError(err)

	fmt.Println("New BlockChain created")
}

func (cli *CommandLine) printAll() {
	chain, err := blockchain.ContinueBlockchain()
	handleError(err)

	// Go through all the blocks created
	iter := chain.Iterator()
	for {
		block, err := iter.Next()
		if err != nil {
			break
		}

		fmt.Printf("Block hash: %x\n\n", block.Hash)
		for _, tx := range block.Transactions {
			fmt.Printf("Transaction Hash: %x\n\n", tx.HashID)

			fmt.Println("INPUTS:")
			for _, in := range tx.Inputs {
				fmt.Printf("PrevTxHash: %x\nOutIdx: %d\nSig: %x\n",
					in.PrevTxHash, in.OutIdx, in.Signature)
			}

			fmt.Println("\nOUTPUTS:")
			for _, out := range tx.Outputs {
				fmt.Printf("Amount: %d\nPubKey: %x\n", out.Value, out.PubKeyHash)
			}

			pow := blockchain.NewProof(block)
			isValid := pow.IsValid()

			fmt.Printf("is valid?: %t", isValid)
			fmt.Printf("\n==================\n\n")
		}
	}
}

// Run runs the command line
func (cli *CommandLine) Run() {
	switch os.Args[1] {
	case "newwallet":
		cli.newWallet()
	case "showwallets":
		cli.showWallets()
	case "getbalance":
		cli.getBalance(os.Args[2])
	case "newtransaction":
		amount, err := strconv.Atoi(os.Args[4])
		if err != nil {
			panic(err)
		}
		cli.newTransaction(os.Args[2], os.Args[3], amount)
	case "newblockchain":
		cli.newBlockchain(os.Args[2])
	case "print":
		cli.printAll()
	default:
		fmt.Println("Command not found")
	}
}
