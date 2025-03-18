#!/bin/bash

# Create a secret in the k8s cluster called demo-w-secret
kubectl create secret generic demo-w-secret --from-literal=STATELY_ACCESS_KEY="$STATELY_ACCESS_KEY"
