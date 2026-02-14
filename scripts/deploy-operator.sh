#!/bin/bash
# Script to deploy the first-operator to the Kind cluster

set -e

OPERATOR_DIR="first-operator"

echo "Deploying first-operator to Kind cluster..."

cd $OPERATOR_DIR

# Install CRDs
echo "Installing CRDs..."
make install

# Build and load Docker image into Kind
echo "Building operator image..."
make docker-build IMG=first-operator:latest

echo "Loading image into Kind cluster..."
kind load docker-image first-operator:latest --name operator-dev

# Deploy the operator
echo "Deploying operator..."
make deploy IMG=first-operator:latest

echo "Waiting for operator to be ready..."
kubectl wait --for=condition=Available deployment/first-operator-controller-manager -n first-operator-system --timeout=120s

echo "Operator deployed successfully!"
echo "To check the status, run: kubectl get pods -n first-operator-system"
