package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
	core "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	rbac "k8s.io/api/rbac/v1"

	"github.com/Ladicle/kubectl-rolesum/pkg/explorer"
	"github.com/Ladicle/kubectl-rolesum/pkg/util/subject"
)

type PrettyPrinter struct {
	out io.Writer
}

func DefaultPrettyPrinter() *PrettyPrinter {
	return &PrettyPrinter{
		out: os.Stdout,
	}
}

func (p *PrettyPrinter) PrintSubject(sub *rbac.Subject) {
	var name string
	if sub.Kind == subject.KindSA {
		name = sub.Namespace + "/" + sub.Name
	} else {
		name = sub.Name
	}
	fmt.Fprintf(p.out, "%v: %v\n", aurora.Yellow(sub.Kind), name)
}

func (p *PrettyPrinter) PrintSA(sa *core.ServiceAccount) {
	p.PrintHeader("Secrets")
	for _, s := range sa.Secrets {
		fmt.Fprintf(p.out, "%v %v/%v\n", bullet, blank2Asterisk(s.Namespace), s.Name)
	}
}

func (p *PrettyPrinter) PrintRolesums(sbjrs []*explorer.SubjectRole) {
	for _, r := range sbjrs {
		fmt.Fprintf(p.out, "%v %v/%v\n", bullet, blank2Asterisk(r.Namespace), r.Name)
	}
}

func (p *PrettyPrinter) PrintPolicies(sbjrs []*explorer.SubjectRole) {
	for i, r := range sbjrs {
		if i != 0 {
			p.BlankLine()
		}
		fmt.Fprintf(p.out, "%v [%v] %v/%v %v  [%v] %v/%v\n",
			bullet,
			aurora.BrightCyan(r.BindingType),
			blank2Asterisk(r.BindingNamespace), r.BindingName,
			arrow,
			aurora.BrightCyan(r.Type),
			blank2Asterisk(r.Namespace), r.Name)

		if len(r.PolicyList.APIPolicies) != 0 {
			p.printAPIPolicy(r.PolicyList.APIPolicies)
		}

		p.BlankLine()
		if len(r.PolicyList.PSPs) != 0 {
			p.printPSP(r.PolicyList.PSPs)
		}
	}
}

func (p *PrettyPrinter) printAPIPolicy(apips []*explorer.ResourceAPIPolicy) {
	tw := p.newTabwriter()
	defer tw.Render()

	tw.Append([]string{
		tabHead("Resource"),
		tabHead("Name"),
		tabHead("Exclude"),
		tabHead("Verbs"),
		tabHead("G L W C U P D DC")})

	tw.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_LEFT})

	for _, policy := range apips {
		tw.Append([]string{
			policy.Resource.String(),
			joinOrAsterisk(policy.ResourceName),
			joinOrDash(policy.NonResourceURL),
			joinOrDash(policy.OtherVerbs),
			apiVerb2CheckTable(policy.APIVerbFlag),
		})
	}
}

func (p *PrettyPrinter) printPSP(psps []*policy.PodSecurityPolicy) {
	tw := p.newTabwriter()
	defer tw.Render()

	tw.Append([]string{
		tabHead("Name"),
		tabHead("PRIV"),
		tabHead("RO-RootFS"),
		tabHead("Volumes"),
		tabHead("Caps"),
		tabHead("SELinux"),
		tabHead("RunAsUser"),
		tabHead("FSgroup"),
		tabHead("SUPgroup")})

	tw.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_CENTER})

	for _, policy := range psps {
		tw.Append([]string{
			policy.Name,
			colorBool(*policy.Spec.AllowPrivilegeEscalation),
			colorBool(policy.Spec.ReadOnlyRootFilesystem),
			joinFsType(policy.Spec.Volumes),
			joinCap(policy.Spec.AllowedCapabilities),
			string(policy.Spec.SELinux.Rule),
			string(policy.Spec.RunAsUser.Rule),
			string(policy.Spec.FSGroup.Rule),
			string(policy.Spec.SupplementalGroups.Rule),
		})
	}
}

func (p *PrettyPrinter) BlankLine() {
	fmt.Fprintln(p.out)
}

func (p *PrettyPrinter) newTabwriter() *tablewriter.Table {
	tw := tablewriter.NewWriter(p.out)
	tw.SetRowSeparator("")
	tw.SetCenterSeparator("")
	tw.SetColumnSeparator("")
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetHeaderLine(false)
	tw.SetAutoWrapText(false)
	return tw
}

func (p *PrettyPrinter) PrintHeader(header string) {
	fmt.Fprintln(p.out, aurora.BrightCyan(header+":"))
}

func (p *PrettyPrinter) printHeaderL2(header string) {
	fmt.Fprintln(p.out, aurora.BrightCyan("  "+header+":"))
}
