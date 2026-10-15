package main

import (
	"os"

	"github.com/mascah/grove-next2/internal/cli"
)

func main() { os.Exit(cli.Main(os.Args[1:])) }
