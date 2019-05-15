package explorer

import (
	"github.com/Ladicle/kubectl-bindrole/pkg/util/subject"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// subjectRoles retrieve RoleBindings for the specified subject.
func subjectRoleBindings(client *kubernetes.Clientset, sbj *rbacv1.Subject) ([]rbacv1.RoleBinding, error) {
	list, err := client.RbacV1().RoleBindings(sbj.Namespace).
		List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var sbjrs []rbacv1.RoleBinding
	for _, b := range list.Items {
		if containSubject(sbj, b.Subjects) {
			sbjrs = append(sbjrs, b)
		}
	}
	return sbjrs, nil
}

// subjectClusterRoles retrieve ClusterRoleBindings for the specified subject.
func subjectClusterRoleBindings(client *kubernetes.Clientset, sbj *rbacv1.Subject) ([]rbacv1.ClusterRoleBinding, error) {
	list, err := client.RbacV1().ClusterRoleBindings().
		List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var sbjcrs []rbacv1.ClusterRoleBinding
	for _, b := range list.Items {
		if containSubject(sbj, b.Subjects) {
			sbjcrs = append(sbjcrs, b)
		}
	}
	return sbjcrs, nil
}

// containSubject returns true if the specified subject is on the list.
func containSubject(s *rbacv1.Subject, list []rbacv1.Subject) bool {
	for _, sub := range list {
		if sameSubject(&sub, s) {
			return true
		}
	}
	return false
}

// sameSubject returns true if s1 equals s2.
func sameSubject(s1, s2 *rbacv1.Subject) bool {
	if s1.Kind != s2.Kind || s1.Name != s2.Name {
		return false
	}
	if s1.Kind == subject.KindSA && s1.Namespace != s2.Namespace {
		return false
	}
	return true
}
