package main

import (
	"fmt"
	"jotacoin/pkg/cli"
)

func main() {
	fmt.Println("==== JOTA COIN ====")
	cli := &cli.CommandLine{}
	cli.Run()
}
