# Development Guide - first-operator

This guide provides detailed information for developers working on the first-operator.

## Architecture

The first-operator follows the Kubernetes Operator pattern using the Kubebuilder framework. It consists of:

1. **Custom Resource Definition (CRD)**: Defines the `FirstApp` resource
2. **Controller**: Watches FirstApp resources and reconciles the desired state
3. **Reconciler**: Contains the business logic for managing resources

### Project Structure

```
first-operator/
├── api/v1alpha1/              # API definitions
│   ├── firstapp_types.go      # FirstApp CRD definition
│   └── groupversion_info.go   # API group metadata
├── internal/controller/       # Controller implementation
│   ├── firstapp_controller.go # Reconciliation logic
│   └── suite_test.go          # Test suite setup
├── config/                    # Kubernetes manifests
│   ├── crd/                   # CRD manifests
│   ├── rbac/                  # RBAC rules
│   ├── manager/               # Operator deployment
│   └── samples/               # Example resources
├── cmd/main.go                # Operator entry point
└── Dockerfile                 # Container image
```

## Custom Resource Definition

The `FirstApp` CRD is defined in `api/v1alpha1/firstapp_types.go`:

```go
type FirstAppSpec struct {
    Size  int32  `json:"size,omitempty"`   // Number of replicas (1-10)
    Image string `json:"image"`             // Container image (required)
    Port  int32  `json:"port,omitempty"`    // Container port (1-65535)
}

type FirstAppStatus struct {
    Conditions        []metav1.Condition `json:"conditions,omitempty"`
    AvailableReplicas int32              `json:"availableReplicas,omitempty"`
}
```

### Field Validation

The CRD includes kubebuilder markers for validation:
- `Size`: Must be between 1 and 10
- `Image`: Required field
- `Port`: Must be between 1 and 65535

## Controller Logic

The controller implements the Reconcile loop:

1. **Fetch**: Get the FirstApp resource
2. **Check**: Verify if a Deployment exists
3. **Create**: If not, create a new Deployment
4. **Update**: If exists, ensure it matches the desired state
5. **Status**: Update the status with current replica count

### Reconciliation Flow

```
FirstApp Created/Updated
         ↓
    Reconcile()
         ↓
  Deployment Exists?
    ↓           ↓
   No          Yes
    ↓           ↓
 Create      Update if needed
    ↓           ↓
 Update Status
```

## Development Workflow

### Making Changes

1. **Modify the API**: Edit `api/v1alpha1/firstapp_types.go`
   ```bash
   # After changes, regenerate code
   make manifests generate
   ```

2. **Update Controller**: Edit `internal/controller/firstapp_controller.go`
   ```bash
   # Test your changes
   make test
   
   # Run the operator
   make run
   ```

3. **Test Changes**: Apply sample resources
   ```bash
   kubectl apply -f config/samples/
   kubectl get firstapps
   kubectl describe firstapp firstapp-sample
   ```

### Running Tests

```bash
# Unit tests
make test

# End-to-end tests (requires cluster)
make test-e2e

# Code coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Debugging

#### Local Debugging with VSCode

1. Set breakpoints in the code
2. Press F5 or select "Debug First Operator"
3. Apply a FirstApp resource to trigger reconciliation

#### Remote Debugging

Deploy the operator with debug flags:
```bash
make docker-build IMG=first-operator:debug
kind load docker-image first-operator:debug
make deploy IMG=first-operator:debug
```

#### View Logs

```bash
# Local run - see terminal output
make run

# Deployed to cluster
kubectl logs -n first-operator-system -l control-plane=controller-manager -f
```

## RBAC Permissions

The controller requires specific permissions defined via kubebuilder markers:

```go
// +kubebuilder:rbac:groups=apps.example.com,resources=firstapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=firstapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
```

Regenerate RBAC manifests after changes:
```bash
make manifests
```

## Adding New Features

### Example: Add Service Creation

1. **Update Types** (if needed):
   ```go
   type FirstAppSpec struct {
       // ... existing fields
       ServicePort int32 `json:"servicePort,omitempty"`
   }
   ```

2. **Update Controller**:
   ```go
   func (r *FirstAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
       // ... existing code
       
       // Create or update Service
       service := r.serviceForFirstApp(firstApp)
       if err := r.createOrUpdateService(ctx, service); err != nil {
           return ctrl.Result{}, err
       }
       
       return ctrl.Result{}, nil
   }
   ```

3. **Add RBAC**:
   ```go
   // +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
   ```

4. **Regenerate and Test**:
   ```bash
   make manifests generate
   make test
   make run
   ```

## Best Practices

### 1. Idempotency
Reconciliation should be idempotent - running it multiple times should produce the same result.

### 2. Error Handling
- Return errors for retryable failures
- Use exponential backoff for rate limiting
- Log errors with context

### 3. Status Updates
- Always update status in a separate call
- Use conditions to track resource state
- Include useful diagnostic information

### 4. Testing
- Write unit tests for business logic
- Use envtest for integration testing
- Test edge cases and error conditions

### 5. Resource Cleanup
- Set owner references for garbage collection
- Implement finalizers for custom cleanup
- Handle deletion gracefully

## Common Tasks

### Update CRD Schema

```bash
# Edit api/v1alpha1/firstapp_types.go
# Regenerate manifests
make manifests

# Reinstall CRDs
make install
```

### Add Webhook Validation

```bash
kubebuilder create webhook --group apps --version v1alpha1 --kind FirstApp --defaulting --programmatic-validation
```

### Release New Version

```bash
# Update version in Makefile
VERSION=v0.1.0

# Build and push image
make docker-build docker-push IMG=registry/first-operator:${VERSION}

# Deploy
make deploy IMG=registry/first-operator:${VERSION}
```

## Troubleshooting

### Controller Not Reconciling

Check:
1. Controller is running: `kubectl get pods -n first-operator-system`
2. RBAC permissions are correct: `kubectl auth can-i create deployments`
3. Events: `kubectl get events -n <namespace>`

### CRD Changes Not Applied

```bash
# Uninstall old CRD
make uninstall

# Reinstall new CRD
make install
```

### Stale Resources

```bash
# Force reconciliation
kubectl annotate firstapp firstapp-sample reconcile=true --overwrite

# Or delete and recreate
kubectl delete firstapp firstapp-sample
kubectl apply -f config/samples/apps_v1alpha1_firstapp.yaml
```

## References

- [Kubebuilder Book](https://book.kubebuilder.io/)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
- [Operator Best Practices](https://sdk.operatorframework.io/docs/best-practices/)
