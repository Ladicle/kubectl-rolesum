package main

import (
	"fmt"
	"os"

	"github.com/Ladicle/kubectl-bindrole/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
	}
}
