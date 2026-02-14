# Quick Start Guide - first-operator

This guide will help you quickly get started with the first-operator example.

## Prerequisites

Make sure you have the following installed:
- Go 1.21+
- Docker
- kubectl
- Kind

## Quick Start (5 minutes)

### 1. Create a Kind Cluster

From the root of the repository:

```bash
make cluster-create
```

This creates a Kubernetes cluster named `operator-dev` using Kind.

### 2. Run the Operator Locally

This is the fastest way to test the operator during development:

```bash
cd first-operator
make install  # Install the CRDs
make run      # Run the operator locally
```

The operator will connect to your Kind cluster and start watching for FirstApp resources.

### 3. Create a FirstApp Resource

Open a new terminal and apply the sample FirstApp:

```bash
kubectl apply -f first-operator/config/samples/apps_v1alpha1_firstapp.yaml
```

### 4. Verify the Deployment

Check that the operator created a deployment:

```bash
# View the FirstApp resource
kubectl get firstapps

# View the deployment
kubectl get deployments

# View the pods
kubectl get pods
```

Expected output:
```
NAME              READY   STATUS    RESTARTS   AGE
firstapp-sample   1/1     Running   0          30s
```

## Understanding the Example

The FirstApp resource you created looks like this:

```yaml
apiVersion: apps.example.com/v1alpha1
kind: FirstApp
metadata:
  name: firstapp-sample
spec:
  size: 2          # Number of replicas
  image: nginx:latest
  port: 80
```

The operator automatically:
1. Creates a Deployment with the specified number of replicas
2. Configures the pods to use the specified image and port
3. Updates the deployment when you change the FirstApp spec
4. Tracks the number of running replicas in the status

## Try It Out

### Scale the Application

Edit the FirstApp to change the number of replicas:

```bash
kubectl edit firstapp firstapp-sample
```

Change `size: 2` to `size: 3` and save. Watch the operator create a new pod:

```bash
kubectl get pods -w
```

### Update the Image

You can also change the container image:

```bash
kubectl patch firstapp firstapp-sample --type merge -p '{"spec":{"image":"nginx:alpine"}}'
```

### View Operator Logs

If you're running the operator locally (with `make run`), you'll see logs in your terminal showing the reconciliation process.

If deployed to the cluster:

```bash
kubectl logs -n first-operator-system deployment/first-operator-controller-manager -f
```

## Deploy to Cluster (Alternative)

Instead of running locally, you can deploy the operator to the cluster:

```bash
# From the root directory
make first-operator-deploy
```

This will:
1. Build a Docker image
2. Load it into Kind
3. Deploy the operator to the cluster

## Cleanup

### Delete the FirstApp Resource

```bash
kubectl delete firstapp firstapp-sample
```

The operator will automatically clean up the deployment and pods.

### Stop the Operator

If running locally, press `Ctrl+C` in the terminal where `make run` is running.

If deployed to cluster:

```bash
cd first-operator
make undeploy
```

### Delete the Kind Cluster

```bash
make cluster-delete
```

## Next Steps

- Explore the controller code in `first-operator/internal/controller/firstapp_controller.go`
- Modify the FirstApp CRD in `first-operator/api/v1alpha1/firstapp_types.go`
- Add new fields to the spec and implement handling in the controller
- Run tests with `make test` in the first-operator directory

## Troubleshooting

### Operator not starting

Make sure:
1. The Kind cluster is running: `kind get clusters`
2. kubectl context is set correctly: `kubectl cluster-info`
3. CRDs are installed: `kubectl get crds | grep firstapp`

### Pods not created

Check operator logs to see if there are errors:
```bash
# If running locally, check the terminal output
# If deployed, check pod logs
kubectl logs -n first-operator-system -l control-plane=controller-manager
```

### CRD already exists error

If you get an error about CRDs already existing:
```bash
cd first-operator
make uninstall
make install
```

## Learning Resources

- [Kubebuilder Book](https://book.kubebuilder.io/) - Complete guide to building operators
- [Controller Runtime Docs](https://pkg.go.dev/sigs.k8s.io/controller-runtime) - API documentation
- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) - Understanding operators
