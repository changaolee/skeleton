package main

import (
	"os"

	"github.com/changaolee/skeleton/internal/sktctl/cmd"
)

func main() {
	command := cmd.NewDefaultSKTCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
