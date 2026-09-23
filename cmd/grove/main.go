package main

import (
	"fmt"
	"os"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/cli"
)

func main() {
	// The launcher of an attempt starts this binary again as the attempt's
	// owner, named by the environment rather than an argument so that a test
	// binary can be its own owner the same way (internal/attempt).
	if dir := os.Getenv(attempt.OwnerEnv); dir != "" {
		os.Exit(attempt.Own(dir))
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "grove: working directory: %s\n", err)
		os.Exit(1)
	}
	os.Exit(cli.Run(os.Args[1:], cwd, os.Stdout, os.Stderr))
}
