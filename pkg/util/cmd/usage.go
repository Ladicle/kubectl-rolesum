package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

var UsageTemplate = fmt.Sprintf(`%v:{{if .Runnable}}
  kubectl {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

%v:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%v:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}

%v:
  - ServiceAccount (default)
  - User
  - Group

%v:
  -h, --help                   Display this help message
  -n, --namespace string       Change the namespace scope for this CLI request
  -k, --subject-kind string    Set SubjectKind to summarize (default: ServiceAccount)
  -o, --options                List of all options for this command
      --version                Show version for this command

Use "kubectl bindrole --options" for a list of all options (applies to this command).
`,
	aurora.Cyan("Usage"),
	aurora.Cyan("Aliases"),
	aurora.Cyan("Examples"),
	aurora.Cyan("SubjectKinds"),
	aurora.Cyan("Options"))

var OptionTemplate = `The following options can be passed to this command:

{{if .HasAvailableLocalFlags}}{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}`

var VersionTemplate = fmt.Sprintf(` %v {{.Name}} version {{.Version}}
`, aurora.Yellow(">"))
