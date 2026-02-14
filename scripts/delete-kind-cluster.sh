#!/bin/bash
# Script to delete the Kind cluster

set -e

CLUSTER_NAME="operator-dev"

echo "Deleting Kind cluster: $CLUSTER_NAME"
kind delete cluster --name $CLUSTER_NAME

echo "Cluster deleted successfully!"
