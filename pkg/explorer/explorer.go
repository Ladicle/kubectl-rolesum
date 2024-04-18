package explorer

import (
	"context"
	"sort"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PolicyExplorer struct {
	client *kubernetes.Clientset
}

func NewPolicyExplorer(client *kubernetes.Clientset) *PolicyExplorer {
	return &PolicyExplorer{client: client}
}

type SubjectRole struct {
	BindingName      string
	BindingNamespace string
	BindingType      string
	Name             string
	Namespace        string
	Type             string
	PolicyList       *SubjectPolicyList
}

type SubjectPolicyList struct {
	APIPolicies []*ResourceAPIPolicy
}

// NamespacedSbjRoles explores bound namespaced roles to the specified subject.
func (e *PolicyExplorer) NamespacedSbjRoles(ctx context.Context, sbj *rbacv1.Subject) ([]*SubjectRole, error) {
	sbjrbs, err := subjectRoleBindings(ctx, e.client, sbj)

	if err != nil {
		return nil, err
	}

	var sbjrs []*SubjectRole

	for _, b := range sbjrbs {
		if b.RoleRef.Kind == "ClusterRole" {
			role, err := e.client.RbacV1().ClusterRoles().
				Get(ctx, b.RoleRef.Name, metav1.GetOptions{})

			if err != nil {
				return nil, err
			}

			sbjpl, err := rule2sbjpl(role.Rules)

			if err != nil {
				return nil, err
			}

			sbjrs = append(sbjrs, &SubjectRole{
				BindingName:      b.Name,
				BindingNamespace: b.Namespace,
				BindingType:      "RB",
				Name:             role.Name,
				Namespace:        role.Namespace,
				Type:             "CR",
				PolicyList:       sbjpl,
			})
		} else if b.RoleRef.Kind == "Role" {
			role, err := e.client.RbacV1().Roles(sbj.Namespace).
				Get(ctx, b.RoleRef.Name, metav1.GetOptions{})

			if err != nil {
				return nil, err
			}

			sbjpl, err := rule2sbjpl(role.Rules)

			if err != nil {
				return nil, err
			}

			sbjrs = append(sbjrs, &SubjectRole{
				BindingName:      b.Name,
				BindingNamespace: b.Namespace,
				BindingType:      "RB",
				Name:             role.Name,
				Namespace:        role.Namespace,
				Type:             "R",
				PolicyList:       sbjpl,
			})
		}
	}
	return sbjrs, nil
}

// ClusterSbjRoles explores bound cluster roles to the specified subject.
func (e *PolicyExplorer) ClusterSbjRoles(ctx context.Context, sbj *rbacv1.Subject) ([]*SubjectRole, error) {
	sbjcrbs, err := subjectClusterRoleBindings(ctx, e.client, sbj)
	if err != nil {
		return nil, err
	}
	var sbjrs []*SubjectRole
	for _, b := range sbjcrbs {
		role, err := e.client.RbacV1().ClusterRoles().
			Get(ctx, b.RoleRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		sbjpl, err := rule2sbjpl(role.Rules)
		if err != nil {
			return nil, err
		}
		sbjrs = append(sbjrs, &SubjectRole{
			BindingName:      b.Name,
			BindingNamespace: b.Namespace,
			BindingType:      "CRB",
			Name:             role.Name,
			Namespace:        role.Namespace,
			Type:             "CR",
			PolicyList:       sbjpl,
		})
	}
	return sbjrs, nil
}

func rule2sbjpl(rules []rbacv1.PolicyRule) (*SubjectPolicyList, error) {
	sbjpl := &SubjectPolicyList{
		APIPolicies: []*ResourceAPIPolicy{},
	}
	rapipMap := make(map[string]*ResourceAPIPolicy)

	for _, rule := range rules {
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
