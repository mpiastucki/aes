package main

import (
	"os"

	"github.com/happymanju/aes/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
