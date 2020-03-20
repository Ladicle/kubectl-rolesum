package printer

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"

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
)

var bullet = aurora.Magenta(uBullet)

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

func colorBool(flag bool) string {
	if flag {
		return aurora.Green("True").String()
	}
	return aurora.Red("False").String()
}

func mark(flag bool) aurora.Value {
	if flag {
		return aurora.Green(uCheck2)
	}
	return aurora.Red(uNG2)
}

func joinCap(caps []v1.Capability) string {
	var scaps []string
	for _, cap := range caps {
		scaps = append(scaps, string(cap))
	}
	return join(scaps)
}

func joinFsType(fsts []v1beta1.FSType) string {
	var sfsts []string
	for _, fstype := range fsts {
		sfsts = append(sfsts, string(fstype))
	}
	return join(sfsts)
}

func tabHead(header string) string {
	return aurora.Yellow(header).String()
}
