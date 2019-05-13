package explorer

import (
	"fmt"
	"strings"

	"github.com/Ladicle/kubectl-bindrole/pkg/util/subject"
	"github.com/logrusorgru/aurora"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const psp = "podsecuritypolicies"

type PolicyExplorer struct {
	Subject *rbacv1.Subject

	client   *kubernetes.Clientset
	targetNS string
}

func NewPolicyExplorer(sub *rbacv1.Subject, client *kubernetes.Clientset) *PolicyExplorer {
	var ns string
	if sub.Kind == subject.KindSA {
		ns = sub.Namespace
	} else {
		ns = metav1.NamespaceAll
	}

	return &PolicyExplorer{
		Subject: sub,

		client:   client,
		targetNS: ns,
	}
}

type Rule struct {
	Resource   Resource
	APIPolicy  APIPolicy
	OtherVerbs []string
}

type Resource struct {
	Name        string
	Group       string
	Subresource string
}

func (r Resource) String() string {
	res := r.Name
	if r.Group != "" {
		res += "." + r.Group
	}
	if r.Subresource != "" {
		res += "/" + r.Subresource
	}
	return res
}

type APIPolicy struct {
	Get              bool
	List             bool
	Watch            bool
	Create           bool
	Update           bool
	Patch            bool
	Delete           bool
	DeleteCollection bool
}

type NamespacedPolicy struct {
	Name                string
	Namespace           string
	Rules               []rbacv1.PolicyRule
	PodSecurityPolicies []policyv1beta1.PodSecurityPolicy
}

func (e *PolicyExplorer) NamespacedPolicy() ([]NamespacedPolicy, error) {
	list, err := e.client.RbacV1().RoleBindings(e.targetNS).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var machiedBindings []rbacv1.RoleBinding
	for _, b := range list.Items {
		if e.isBind(b.Subjects) {
			machiedBindings = append(machiedBindings, b)
		}
	}

	var policies []NamespacedPolicy
	for _, b := range machiedBindings {
		role, err := e.client.RbacV1().Roles(e.targetNS).Get(b.RoleRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		policy := NamespacedPolicy{Name: role.Name, Namespace: role.Namespace}

		fmt.Printf("Role: %v/%v\n", role.Namespace, aurora.Green(role.Name).Bold())

		fmt.Printf("  %25v\t", "Resource")
		fmt.Printf("[%v]\t", "Verbs")
		fmt.Printf("[%v]\t", "ResourceName")
		fmt.Printf("[%v]\t", "NonResourceURL")
		fmt.Println()
		for _, rule := range role.Rules {
			ress := rule2resources(&rule)
			for _, res := range ress {
				fmt.Printf("  %25v\t", res)
				fmt.Printf("[%v]\t", strings.Join(rule.Verbs, ","))
				fmt.Printf("[%v]\t", strings.Join(rule.ResourceNames, ","))
				fmt.Printf("[%v]\t", strings.Join(rule.NonResourceURLs, ","))
				fmt.Println()
			}
		}
		policies = append(policies, policy)
	}
	return policies, nil
}

func rule2resources(rule *rbacv1.PolicyRule) []Resource {
	var resources []Resource
	for _, res := range rule.Resources {
		ss := strings.Split(res, "/")
		name := ss[0]

		var sub string
		if len(ss) == 2 {
			sub = ss[1]
		}

		for _, group := range rule.APIGroups {
			resources = append(resources, Resource{
				Name:        name,
				Group:       group,
				Subresource: sub,
			})
		}
	}
	return resources
}

type ClusterPolicy struct {
	Name                string
	Rules               []rbacv1.PolicyRule
	PodSecurityPolicies []policyv1beta1.PodSecurityPolicy
}

func (e *PolicyExplorer) ClusterPolicy() ([]ClusterPolicy, error) {
	list, err := e.client.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var machiedBindings []rbacv1.ClusterRoleBinding
	for _, b := range list.Items {
		if e.isBind(b.Subjects) {
			machiedBindings = append(machiedBindings, b)
		}
	}

	var policies []ClusterPolicy
	for _, b := range machiedBindings {
		role, err := e.client.RbacV1().ClusterRoles().Get(b.RoleRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		policy := ClusterPolicy{Name: role.Name}

		fmt.Printf("ClusterRole: %v\n", aurora.Green(role.Name).Bold())
		fmt.Printf("  %25v\t", "Resource")
		fmt.Printf("[%v]\t", "Verbs")
		fmt.Printf("[%v]\t", "ResourceName")
		fmt.Printf("[%v]\t", "NonResourceURL")
		fmt.Println()
		for _, rule := range role.Rules {
			ress := rule2resources(&rule)
			for _, res := range ress {
				fmt.Printf("  %25v\t", res)
				fmt.Printf("[%v]\t", strings.Join(rule.Verbs, ","))
				fmt.Printf("[%v]\t", strings.Join(rule.ResourceNames, ","))
				fmt.Printf("[%v]\t", strings.Join(rule.NonResourceURLs, ","))
				fmt.Println()
			}
		}

		policies = append(policies, policy)
	}
	return policies, nil
}

func (e *PolicyExplorer) isBind(subjects []rbacv1.Subject) bool {
	for _, sub := range subjects {
		if sub.Kind == e.Subject.Kind && sub.Name == e.Subject.Name {
			if sub.Kind == subject.KindSA &&
				sub.Namespace != e.Subject.Namespace {
				continue
			}
			return true
		}
	}
	return false
}

// func (e *PolicyExplorer) ClusterRoles()  {
// 		fmt.Println("# Cluster Role")
// 	crbList, err := client.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// TODO: Kind=Group,User
// 	var cBindedList []rbacv1.ClusterRoleBinding
// 	for _, crb := range crbList.Items {
// 		for _, sub := range crb.Subjects {
// 			if sub.Kind != "ServiceAccount" {
// 				continue
// 			}
// 			if sub.Namespace == ns && sub.Name == subName {
// 				cBindedList = append(cBindedList, crb)
// 			}
// 		}
// 	}

// 	pspList = []string{}
// 	for _, crb := range cBindedList {
// 		crole, err := client.RbacV1().ClusterRoles().
// 			Get(crb.RoleRef.Name, metav1.GetOptions{})
// 		if err != nil {
// 			log.Println(crb)
// 			log.Println(err)
// 			continue
// 		}
// 		fmt.Printf("Name: %v\n", crole.Name)
// 		fmt.Println("Rules:")
// 		for _, rule := range crole.Rules {
// 			for _, r := range rule.Resources {
// 				if r == psp {
// 					pspList = append(pspList, rule.ResourceNames...)
// 					continue
// 				}
// 				fmt.Printf("%v\t", r)
// 				fmt.Printf("%v\t", rule.Verbs)
// 				fmt.Printf("%v\t", strings.Join(rule.APIGroups, ","))
// 				fmt.Printf("%v\t", rule.ResourceNames)
// 				fmt.Printf("%v\t", rule.NonResourceURLs)
// 				fmt.Println()
// 			}
// 		}
// 		// TODO
// 		fmt.Printf("AggregationRules: %v\n\n", crole.AggregationRule)
// 	}

// 	fmt.Println("## Pod Security Policy")
// 	for _, name := range pspList {
// 		psp, err := client.PolicyV1beta1().PodSecurityPolicies().Get(name, metav1.GetOptions{})
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		fmt.Printf("%v\t", psp.Name)
// 		fmt.Printf("%v\t", psp.Spec.AllowPrivilegeEscalation)
// 		fmt.Printf("%v\t", psp.Spec.AllowedCapabilities)
// 		fmt.Printf("%v\t", psp.Spec.SELinux)
// 		fmt.Printf("%v\t", psp.Spec.RunAsUser)
// 		fmt.Printf("%v\t", psp.Spec.FSGroup)
// 		fmt.Printf("%v\t", psp.Spec.SupplementalGroups)
// 		fmt.Printf("%v\t", psp.Spec.ReadOnlyRootFilesystem)
// 		fmt.Printf("%v\t", psp.Spec.SupplementalGroups)
// 		fmt.Printf("%v\t", psp.Spec.ReadOnlyRootFilesystem)
// 		fmt.Printf("%v\t", psp.Spec.Volumes)
// 	}
// }
