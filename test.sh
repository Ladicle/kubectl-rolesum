#!/bin/bash -eu

echo; echo "Creating ServiceAccount..."
kubectl create sa test-user

echo; echo "Creating PSP..."
cat <<EOF | kubectl apply -f -
apiVersion: extensions/v1beta1
kind: PodSecurityPolicy
metadata:
  name: test-psp
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'
spec:
  privileged: false  
  allowPrivilegeEscalation: false 
  allowedCapabilities:
  - '*'
  volumes:
  - '*'
  hostNetwork: false 
  hostPorts:
  - min: 0
    max: 65535
  hostIPC: true 
  hostPID: true 
  runAsUser:
    rule: 'MustRunAsNonRoot' 
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
EOF

echo; echo "Creating Role..."
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: test-role
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "edit", "exec"]
- apiGroups: ['policy']
  resources: ['podsecuritypolicies']
  verbs:     ['use']
  resourceNames:
  - "test-psp"
EOF

echo; echo "Binding Role..."
kubectl create rolebinding test \
        --role=test-role \
        --serviceaccount=default:test-user

echo; echo "Binding ClusterRole..."
kubectl create clusterrolebinding test --clusterrole edit --serviceaccount default:test-user

echo; echo "Test..."
./_output/kubectl-bindrole test-user
