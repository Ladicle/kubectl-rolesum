package cmd

import (
	"errors"
	"fmt"

	"github.com/Ladicle/kubectl-bindrole/pkg/explorer"
	"github.com/Ladicle/kubectl-bindrole/pkg/util/printer"
	"github.com/Ladicle/kubectl-bindrole/pkg/util/subject"
	"github.com/spf13/pflag"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	// set values via build flags
	command string
	version string
	commit  string

	subkind string
	verflag bool
)

func Execute() error {
	fsets := pflag.CommandLine
	fsets.StringVarP(&subkind, "subject-kind", "k", subject.KindSA,
		"The Kind of subject which is bound Roles.")
	fsets.BoolVarP(&verflag, "version", "v", false, "Print command version")
	cfgFlags := genericclioptions.NewConfigFlags(true)
	cfgFlags.AddFlags(fsets)

	pflag.Parse()

	if verflag {
		fmt.Printf("%v - %v - %v", command, version, commit)
		return nil
	}
	if pflag.NArg() != 1 {
		return errors.New("subject name is required")
	}

	sub := &rbacv1.Subject{
		Name: pflag.Arg(0),
		Kind: subkind,
	}
	if sub.Kind == subject.KindSA {
		k8sCfg := cfgFlags.ToRawKubeConfigLoader()
		ns, _, err := k8sCfg.Namespace()
		if err != nil {
			return err
		}
		sub.Namespace = ns
	}

	cfg, err := cfgFlags.ToRESTConfig()
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(cfg)
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
