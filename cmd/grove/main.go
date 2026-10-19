package main

import (
	"os"

	"github.com/mascah/grove/internal/cli"
)

func main() { os.Exit(cli.Main(os.Args[1:])) }
