package cmd

import (
	"errors"
	"fmt"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/Ladicle/kubectl-bindrole/pkg/explorer"
	"github.com/Ladicle/kubectl-bindrole/pkg/util/printer"
	"github.com/Ladicle/kubectl-bindrole/pkg/util/subject"
)

var (
	// set values via build flags
	command string
	version string
	commit  string
)

type Option struct {
	SubjectKind string
	SubjectName string

	f cmdutil.Factory
}

func NewBindroleCmd() *cobra.Command {
	opt := Option{}
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("%s <SubjectName>", command),
		Version:               fmt.Sprintf("%v @%v", version, commit),
		DisableFlagsInUseLine: true,
		Short:                 "Summarize RBAC roles for the specified subject",
		Long:                  templates.LongDesc("Summarize RBAC roles for the specified subject"),
		Example: templates.Examples(fmt.Sprintf(`# Summarize roles tied to the "ci-bot" ServiceAccount.
%s ci-bot

# Summarize roles tied to the "developer" Group.
%s developer -k Group`, command, command)),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(opt.Validate(cmd, args))
			cmdutil.CheckErr(opt.Run())
		},
	}

	templates.ActsAsRootCommand(cmd, []string{"options"})

	fsets := cmd.PersistentFlags()
	cfgFlags := genericclioptions.NewConfigFlags(true)
	cfgFlags.AddFlags(fsets)
	matchVersionFlags := cmdutil.NewMatchVersionFlags(cfgFlags)
	matchVersionFlags.AddFlags(fsets)

	opt.f = cmdutil.NewFactory(matchVersionFlags)

	cmd.Flags().StringVarP(&opt.SubjectKind, "subject-kind", "k", subject.KindSA,
		fmt.Sprintf("subject kind (available: %s, %s or %s)", subject.KindSA, subject.KindGroup, subject.KindUser))

	return cmd
}

func (o *Option) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("subject name is required")
	}
	o.SubjectName = args[0]
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
	pp.PrintSubject(sub)
	if sub.Kind == subject.KindSA {
		sa, err := client.CoreV1().ServiceAccounts(sub.Namespace).
			Get(sub.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		pp.PrintSA(sa)
	}

	pp.BlankLine()
	pp.PrintHeader("Policies")
	pp.PrintPolicies(nsp)
	pp.BlankLine()
	pp.PrintPolicies(clusterp)

	return nil
}
