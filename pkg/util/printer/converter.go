package printer

import (
	"fmt"
	"strings"

	aurora "github.com/logrusorgru/aurora/v3"

	"github.com/Ladicle/kubectl-rolesum/pkg/explorer"
)

const (
	allNamespace = "*"

	uBullet        = "•"
	uCheckBoxBlank = "☐"
	uCheckBoxNG    = "☒"
	uCheckBoxOK    = "☑"
	uNG            = "✗"
	uNG2           = "✖"
	uCheck         = "✓"
	uCheck2        = "✔"
	uArrow         = "⟶"
)

var bullet = aurora.Magenta(uBullet)
var arrow = aurora.BrightCyan(uArrow)

func blank2Asterisk(s string) string {
	if strings.TrimSpace(s) == "" {
		return allNamespace
	}
	return s
}

func joinOrAsterisk(list []string) string {
	if len(list) == 0 {
		return "[*]"
	}
	return join(list)
}

func joinOrDash(list []string) string {
	if len(list) == 0 {
		return "[-]"
	}
	return join(list)
}

func join(list []string) string {
	return "[" + strings.Join(list, ", ") + "]"
}

func apiVerb2CheckTable(flag uint) string {
	return fmt.Sprintf("%v %v %v %v %v %v %v %v",
		mark(flag&explorer.VerbGet == explorer.VerbGet),
		mark(flag&explorer.VerbList == explorer.VerbList),
		mark(flag&explorer.VerbWatch == explorer.VerbWatch),
		mark(flag&explorer.VerbCreate == explorer.VerbCreate),
		mark(flag&explorer.VerbUpdate == explorer.VerbUpdate),
		mark(flag&explorer.VerbPatch == explorer.VerbPatch),
		mark(flag&explorer.VerbDelete == explorer.VerbDelete),
		mark(flag&explorer.VerbDeletionC == explorer.VerbDeletionC))
}

func mark(flag bool) aurora.Value {
	if flag {
		return aurora.Green(uCheck2)
	}
	return aurora.Red(uNG2)
}

func tabHead(header string) string {
	return aurora.Yellow(header).String()
}
