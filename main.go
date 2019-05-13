package main

import (
	"fmt"
	"os"

	"github.com/Ladicle/kubectl-bindrole/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
	}
}
