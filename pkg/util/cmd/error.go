package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func CheckErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, " %v Error: %v.\n", aurora.Red(">"), err)
	fmt.Fprintf(os.Stderr, " %v Run %v command for the usage.\n", aurora.Red(">"), aurora.Yellow("kubectl bindrole -h"))
	os.Exit(1)
}
