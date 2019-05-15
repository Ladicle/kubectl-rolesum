# kubectl-bindrole

Finding Kubernetes Roles bound to a specified ServiceAccount, Group or User.


## Design

```bash
$ kubectl bindrole test-user

[ServiceAccount] default/test-user
Secrets:
* default/test-user-token
BindedRoles:
* */edit
* default/test-role

Policies:
- Name: default/test-role
  APIPolicies: |-

  PodSecurityPolicies:  |-

- Name: edit
  APIPolicies: |-
  PodSecurityPolicies:  |-

```
