package main

import (
	"fmt"
	"os"

	"github.com/Ladicle/kubectl-bindrole/cmd"
)

func main() {
	command := cmd.NewBindroleCmd()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
	}
}
