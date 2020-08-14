#!/bin/bash -eu

echo; echo "Clean up..."
kubectl delete sa test-user
kubectl delete psp test-psp
kubectl delete role test-role
kubectl delete rolebinding test
kubectl delete clusterrolebinding test
kubectl delete rolebinding test-group
kubectl delete clusterrolebinding test-group
