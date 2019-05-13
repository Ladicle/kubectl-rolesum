package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/Ladicle/kubectl-bindrole/pkg/util/subject"
	"github.com/logrusorgru/aurora"
	rbacv1 "k8s.io/api/rbac/v1"
)

type PrettyPrinter struct {
	out io.Writer
}

func DefaultPrettyPrinter() *PrettyPrinter {
	return &PrettyPrinter{
		out: os.Stdout,
	}
}

func (p *PrettyPrinter) PrintSubject(sub *rbacv1.Subject) {
	var name string
	if sub.Kind == subject.KindSA {
		name = sub.Namespace + "/" + sub.Name
	} else {
		name = sub.Name
	}
	fmt.Fprintf(p.out, "%v %v\n", aurora.Green(sub.Kind).Bold(), name)
}

func (p *PrettyPrinter) BlankLine() {
	fmt.Fprintln(p.out)
}
