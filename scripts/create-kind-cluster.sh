#!/bin/bash
# Script to create a Kind cluster for operator development

set -e

CLUSTER_NAME="operator-dev"
KIND_CONFIG="kind-config.yaml"

echo "Creating Kind cluster: $CLUSTER_NAME"

# Check if cluster already exists
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
    echo "Cluster $CLUSTER_NAME already exists."
    read -p "Do you want to delete and recreate it? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Deleting existing cluster..."
        kind delete cluster --name $CLUSTER_NAME
    else
        echo "Using existing cluster."
        exit 0
    fi
fi

# Create the cluster
echo "Creating new cluster..."
kind create cluster --config $KIND_CONFIG

# Wait for cluster to be ready
echo "Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=120s

echo "Cluster created successfully!"
echo "To use this cluster, run: kubectl cluster-info --context kind-$CLUSTER_NAME"
