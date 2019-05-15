package explorer

import (
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
)

type Resource struct {
	Name        string
	Groups      []string
	Subresource string
}

func (r Resource) String() string {
	res := r.Name
	if r.Groups != nil {
		if len(r.Groups) == 1 {
			if r.Groups[0] != "" {
				res += "." + r.Groups[0]
			}
		} else {
			res += ".[" + strings.Join(r.Groups, ",") + "]"
		}
	}
	if r.Subresource != "" {
		res += "/" + r.Subresource
	}
	return res
}

// rule2res converts Rule into the Resource list.
func rule2res(rule *rbacv1.PolicyRule) []Resource {
	var resources []Resource
	for _, res := range rule.Resources {
		ss := strings.Split(res, "/")
		name := ss[0]

		var sub string
		if len(ss) == 2 {
			sub = ss[1]
		}
		resources = append(resources, Resource{
			Name:        name,
			Groups:      rule.APIGroups,
			Subresource: sub,
		})
	}
	return resources
}
