# my-operator

A Kubernetes operator development workspace that contains multiple operator projects for developing and testing various Kubernetes Operators.

## Overview

This repository serves as a workspace for developing Kubernetes operators using [Kubebuilder](https://kubebuilder.io/). It includes example operators and utilities for local development and testing.

## Prerequisites

- [Go](https://golang.org/doc/install) (v1.21+)
- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) (Kubernetes in Docker)
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html#installation) (v4.0+)
- [VSCode](https://code.visualstudio.com/) with [GitHub Copilot](https://github.com/features/copilot) (recommended)

## Project Structure

```
.
├── .vscode/              # VSCode configuration
├── first-operator/       # Example operator project
├── scripts/              # Utility scripts
├── kind-config.yaml      # Kind cluster configuration
└── README.md
```

## Getting Started

### 1. Setup Local Kubernetes Cluster

Create a Kind cluster for local development:

```bash
./scripts/create-kind-cluster.sh
```

This will create a Kind cluster named `operator-dev` with one control-plane node and one worker node.

### 2. Build and Run the first-operator

Navigate to the first-operator directory:

```bash
cd first-operator
```

#### Option A: Run locally (outside cluster)

```bash
# Install CRDs into the cluster
make install

# Run the operator locally (connects to your current kubeconfig context)
make run
```

#### Option B: Deploy to cluster

```bash
# Build and deploy to the Kind cluster
cd ..
./scripts/deploy-operator.sh
```

### 3. Test the Operator

Create a sample FirstApp custom resource:

```bash
cd first-operator
kubectl apply -f config/samples/apps_v1alpha1_firstapp.yaml
```

Check the created resources:

```bash
# View the FirstApp resource
kubectl get firstapps

# View the deployment created by the operator
kubectl get deployments

# View the pods
kubectl get pods
```

## first-operator

The `first-operator` is an example Kubernetes operator that manages a custom resource called `FirstApp`. It automatically creates and manages Kubernetes Deployments based on the FirstApp specifications.

### FirstApp Custom Resource

The FirstApp CRD allows you to define:

- **size**: Number of replicas (1-10)
- **image**: Container image to deploy (required)
- **port**: Port to expose (1-65535)

Example:

```yaml
apiVersion: apps.example.com/v1alpha1
kind: FirstApp
metadata:
  name: firstapp-sample
spec:
  size: 2
  image: nginx:latest
  port: 80
```

### Features

- Automatically creates Deployments for FirstApp resources
- Manages replica count
- Updates deployment when FirstApp spec changes
- Tracks running replicas in status

## Development

### VSCode Setup

This repository includes VSCode configuration with recommended extensions:

- Go extension with language server
- GitHub Copilot
- Kubernetes Tools
- YAML support

Open the workspace in VSCode and install the recommended extensions when prompted.

### Debugging

Use the VSCode debugger configuration to debug the operator:

1. Set breakpoints in the code
2. Press F5 or use "Debug First Operator" launch configuration
3. The operator will run locally with the debugger attached

### Creating a New Operator

To create a new operator in this workspace:

```bash
# Create a new directory
mkdir my-new-operator
cd my-new-operator

# Initialize with kubebuilder
kubebuilder init --domain example.com --repo github.com/yubo-yue/my-operator/my-new-operator

# Create API
kubebuilder create api --group <group> --version <version> --kind <Kind>
```

## Cleanup

To delete the Kind cluster:

```bash
./scripts/delete-kind-cluster.sh
```

To uninstall the operator from the cluster:

```bash
cd first-operator
make undeploy
```

## Useful Commands

### first-operator Commands

```bash
cd first-operator

# Install CRDs
make install

# Uninstall CRDs
make uninstall

# Build operator binary
make build

# Run tests
make test

# Generate manifests (CRDs, RBAC, etc.)
make manifests

# Build and push Docker image
make docker-build docker-push IMG=<registry>/first-operator:tag

# Deploy to cluster
make deploy IMG=<registry>/first-operator:tag

# Undeploy from cluster
make undeploy

# Run linters
make lint
```

### Kubernetes Commands

```bash
# View all FirstApp resources
kubectl get firstapps

# Describe a FirstApp
kubectl describe firstapp firstapp-sample

# View operator logs (when deployed to cluster)
kubectl logs -n first-operator-system deployment/first-operator-controller-manager -f

# View events
kubectl get events --sort-by='.lastTimestamp'
```

## Resources

- [Kubebuilder Book](https://book.kubebuilder.io/)
- [Kubernetes Operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Kind Documentation](https://kind.sigs.k8s.io/)

## License

See [LICENSE](LICENSE) file for details.
