#!/bin/bash -eu

echo; echo "Creating ServiceAccount..."
kubectl create sa test-user --dry-run=client -o yaml | kubectl apply -f -

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
EOF

echo; echo "Binding Role..."
kubectl create rolebinding test \
        --role=test-role \
        --serviceaccount=default:test-user --dry-run=client -o yaml | kubectl apply -f -

echo; echo "Binding ClusterRole..."
kubectl create clusterrolebinding test --clusterrole edit --serviceaccount default:test-user --dry-run=client -o yaml | kubectl apply -f -

echo; echo "Binding Role[Group]..."
kubectl create rolebinding test-group \
        --role=test-role \
        --group developer --dry-run=client -o yaml | kubectl apply -f -

echo; echo "Binding ClusterRole[Group]..."
kubectl create clusterrolebinding test-group --clusterrole edit --group developer --dry-run=client -o yaml | kubectl apply -f -

echo; echo "Test..."
./_output/kubectl-rolesum test-user

echo; echo "Test[Group]..."
./_output/kubectl-rolesum -k Group developer

./clean.sh
