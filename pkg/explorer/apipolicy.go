package explorer

import rbacv1 "k8s.io/api/rbac/v1"

const (
	VerbGet uint = 1 << iota
	VerbList
	VerbWatch
	VerbCreate
	VerbUpdate
	VerbPatch
	VerbDelete
	VerbDeletionC
)

type ResourceAPIPolicy struct {
	Resource       Resource
	APIVerbFlag    uint
	OtherVerbs     []string
	ResourceName   []string
	NonResourceURL []string
}

func NewResourceAPIPolicy(res Resource, rule rbacv1.PolicyRule) *ResourceAPIPolicy {
	rapip := &ResourceAPIPolicy{
		Resource:       res,
		OtherVerbs:     []string{},
		ResourceName:   rule.ResourceNames,
		NonResourceURL: rule.NonResourceURLs,
	}
	rapip.SetVerbs(rule.Verbs)
	return rapip
}

func (r *ResourceAPIPolicy) SetVerbs(verbs []string) {
	for _, v := range verbs {
		switch v {
		case "get":
			r.APIVerbFlag |= VerbGet
		case "list":
			r.APIVerbFlag |= VerbList
		case "update":
			r.APIVerbFlag |= VerbUpdate
		case "delete":
			r.APIVerbFlag |= VerbDelete
		case "deletecollection":
			r.APIVerbFlag |= VerbDeletionC
		case "patch":
			r.APIVerbFlag |= VerbPatch
		case "create":
			r.APIVerbFlag |= VerbCreate
		case "watch":
			r.APIVerbFlag |= VerbWatch
		default:
			r.OtherVerbs = append(r.OtherVerbs, v)
		}
	}
}
