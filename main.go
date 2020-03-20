package main

import (
	"github.com/spf13/pflag"

	"github.com/Ladicle/kubectl-rolesum/cmd"
	cmdutil "github.com/Ladicle/kubectl-rolesum/pkg/util/cmd"
)

func init() {
	flags := pflag.NewFlagSet("kubectl-rolesum", pflag.ExitOnError)
	pflag.CommandLine = flags
}

func main() {
	command := cmd.NewRolesumCmd()
	cmdutil.CheckErr(command.Execute())
}
