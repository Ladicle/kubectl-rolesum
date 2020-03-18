package main

import (
	"github.com/spf13/pflag"

	"github.com/Ladicle/kubectl-bindrole/cmd"
	cmdutil "github.com/Ladicle/kubectl-bindrole/pkg/util/cmd"
)

func init() {
	flags := pflag.NewFlagSet("kubectl-bindrole", pflag.ExitOnError)
	pflag.CommandLine = flags
}

func main() {
	command := cmd.NewBindroleCmd()
	cmdutil.CheckErr(command.Execute())
}
