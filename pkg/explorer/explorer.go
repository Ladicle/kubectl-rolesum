package explorer

import (
	"sort"

	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const psp = "podsecuritypolicies"

type PolicyExplorer struct {
	client *kubernetes.Clientset
}

func NewPolicyExplorer(client *kubernetes.Clientset) *PolicyExplorer {
	return &PolicyExplorer{client: client}
}

type SubjectRole struct {
	Name       string
	Namespace  string
	PolicyList *SubjectPolicyList
}

type SubjectPolicyList struct {
	APIPolicies []*ResourceAPIPolicy
	PSPs        []*policyv1beta1.PodSecurityPolicy
}

// NamespacedSbjRoles explores bound namespaced roles to the specified subject.
func (e *PolicyExplorer) NamespacedSbjRoles(sbj *rbacv1.Subject) ([]*SubjectRole, error) {
	sbjrbs, err := subjectRoleBindings(e.client, sbj)
	if err != nil {
		return nil, err
	}
	var sbjrs []*SubjectRole
	for _, b := range sbjrbs {
		role, err := e.client.RbacV1().Roles(sbj.Namespace).
			Get(b.RoleRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		sbjpl, err := rule2sbjpl(e.client, role.Rules)
		if err != nil {
			return nil, err
		}
		sbjrs = append(sbjrs, &SubjectRole{
			Name:       role.Name,
			Namespace:  role.Namespace,
			PolicyList: sbjpl,
		})
	}
	return sbjrs, nil
}

// ClusterSbjRoles explores bound cluster roles to the specified subject.
func (e *PolicyExplorer) ClusterSbjRoles(sbj *rbacv1.Subject) ([]*SubjectRole, error) {
	sbjcrbs, err := subjectClusterRoleBindings(e.client, sbj)
	if err != nil {
		return nil, err
	}
	var sbjrs []*SubjectRole
	for _, b := range sbjcrbs {
		role, err := e.client.RbacV1().ClusterRoles().
			Get(b.RoleRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		sbjpl, err := rule2sbjpl(e.client, role.Rules)
		if err != nil {
			return nil, err
		}
		sbjrs = append(sbjrs, &SubjectRole{
			Name:       role.Name,
			Namespace:  role.Namespace,
			PolicyList: sbjpl,
		})
	}
	return sbjrs, nil
}

func rule2sbjpl(client *kubernetes.Clientset, rules []rbacv1.PolicyRule) (*SubjectPolicyList, error) {
	sbjpl := &SubjectPolicyList{
		APIPolicies: []*ResourceAPIPolicy{},
		PSPs:        []*policyv1beta1.PodSecurityPolicy{},
	}
	rapipMap := make(map[string]*ResourceAPIPolicy)

	for _, rule := range rules {
		// Set Pod-Security-Policy
		if len(rule.Resources) == 1 && rule.Resources[0] == psp {
			for _, name := range rule.ResourceNames {
				psp, err := client.PolicyV1beta1().PodSecurityPolicies().
					Get(name, metav1.GetOptions{})
				if err != nil {
					return nil, err
				}
				sbjpl.PSPs = append(sbjpl.PSPs, psp)
			}
			continue
		}

		// Set API policies
		ress := rule2res(&rule)
		for _, res := range ress {
			respath := res.String()
			v, ok := rapipMap[respath]
			if ok {
				if equalStrings(v.ResourceName, rule.ResourceNames) &&
					equalStrings(v.NonResourceURL, rule.NonResourceURLs) {
					v.SetVerbs(rule.Verbs)
					rapipMap[respath] = v
					continue
				}
			}
			rapipMap[respath] = NewResourceAPIPolicy(res, rule)
		}
	}

	var keyarr []string
	for k := range rapipMap {
		keyarr = append(keyarr, k)

	}
	sort.Strings(keyarr)
	for _, key := range keyarr {
		sbjpl.APIPolicies = append(sbjpl.APIPolicies, rapipMap[key])
	}
	return sbjpl, nil
}
