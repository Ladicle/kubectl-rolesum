package cmd

import (
	"errors"
	"fmt"
	"os"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/Ladicle/kubectl-rolesum/pkg/explorer"
	brcmdutil "github.com/Ladicle/kubectl-rolesum/pkg/util/cmd"
	"github.com/Ladicle/kubectl-rolesum/pkg/util/printer"
	"github.com/Ladicle/kubectl-rolesum/pkg/util/subject"
)

var (
	// set values via build flags
	command string
	version string
	commit  string

	rolesumExample = fmt.Sprintf(`%v
kubectl rolesum ci-bot

%v
kubectl rolesum -k Group developer`,
		aurora.BrightBlack("# Summarize roles bound to the \"ci-bot\" ServiceAccount."),
		aurora.BrightBlack("# Summarize roles bound to the \"developer\" Group."))
)

type Option struct {
	SubjectKind    string
	SubjectName    string
	ShowOptionFlag bool

	f cmdutil.Factory
}

func NewRolesumCmd() *cobra.Command {
	opt := Option{}
	cmd := &cobra.Command{
		Use: fmt.Sprintf("rolesum [options] <%v>",
			aurora.Yellow("SubjectName")),
		Version:               version,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Short:                 "Summarize RBAC roles for the specified subject",
		Long:                  "Summarize RBAC roles for the specified subject",
		Example:               templates.Examples(rolesumExample),
		Run: func(cmd *cobra.Command, args []string) {
			brcmdutil.CheckErr(opt.Validate(cmd, args))
			brcmdutil.CheckErr(opt.Run())
		},
	}

	cmd.SetVersionTemplate(brcmdutil.VersionTemplate)
	cmd.SetUsageTemplate(brcmdutil.UsageTemplate)

	fsets := cmd.PersistentFlags()
	cfgFlags := genericclioptions.NewConfigFlags(true)
	cfgFlags.AddFlags(fsets)
	matchVersionFlags := cmdutil.NewMatchVersionFlags(cfgFlags)
	matchVersionFlags.AddFlags(fsets)

	opt.f = cmdutil.NewFactory(matchVersionFlags)

	cmd.Flags().StringVarP(&opt.SubjectKind, "subject-kind", "k", subject.KindSA, "Set SubjectKind to summarize")
	cmd.Flags().BoolVarP(&opt.ShowOptionFlag, "options", "o", false, "List of all options for this command")

	return cmd
}

func (o *Option) Validate(cmd *cobra.Command, args []string) error {
	if o.ShowOptionFlag {
		cmd.SetUsageTemplate(brcmdutil.OptionTemplate)
		fmt.Println(cmd.UsageString())
		os.Exit(0)
	}

	if len(args) == 0 {
		return errors.New(fmt.Sprintf("<%v> is required argument",
			aurora.Cyan("SubjectName").Bold()))
	}
	o.SubjectName = args[0]

	switch o.SubjectKind {
	case subject.KindGroup, subject.KindSA, subject.KindUser:
	default:
		return errors.New(fmt.Sprintf("\"%v\" is unknown SubjectKind",
			aurora.Cyan(o.SubjectKind).Bold()))
	}
	return nil
}

func (o *Option) Run() error {
	sub := &rbacv1.Subject{
		Name: o.SubjectName,
		Kind: o.SubjectKind,
	}
	if sub.Kind == subject.KindSA {
		k8sCfg := o.f.ToRawKubeConfigLoader()
		ns, _, err := k8sCfg.Namespace()
		if err != nil {
			return err
		}
		sub.Namespace = ns
	}

	client, err := o.f.KubernetesClientSet()
	if err != nil {
		return err
	}

	exp := explorer.NewPolicyExplorer(client)
	nsp, err := exp.NamespacedSbjRoles(sub)
	if err != nil {
		return err
	}
	clusterp, err := exp.ClusterSbjRoles(sub)
	if err != nil {
		return err
	}

	pp := printer.DefaultPrettyPrinter()

	if sub.Kind == subject.KindSA {
		sa, err := client.CoreV1().ServiceAccounts(sub.Namespace).
			Get(sub.Name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return errors.New(fmt.Sprintf("ServiceAccount \"%v\" not found",
				aurora.Cyan(fmt.Sprintf("%v/%v", sub.Namespace, sub.Name)).Bold()))
		} else if err != nil {
			return err
		}
		pp.PrintSubject(sub)
		pp.PrintSA(sa)
	} else {
		pp.PrintSubject(sub)
	}

	pp.BlankLine()
	pp.PrintHeader("Policies")
	pp.PrintPolicies(nsp)
	pp.BlankLine()
	pp.PrintPolicies(clusterp)

	return nil
}
