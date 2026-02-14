.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Cluster Management

.PHONY: cluster-create
cluster-create: ## Create a Kind cluster for development
	./scripts/create-kind-cluster.sh

.PHONY: cluster-delete
cluster-delete: ## Delete the Kind cluster
	./scripts/delete-kind-cluster.sh

##@ First Operator

.PHONY: first-operator-build
first-operator-build: ## Build first-operator
	cd first-operator && make build

.PHONY: first-operator-test
first-operator-test: ## Test first-operator
	cd first-operator && make test

.PHONY: first-operator-run
first-operator-run: ## Run first-operator locally
	cd first-operator && make install && make run

.PHONY: first-operator-deploy
first-operator-deploy: ## Deploy first-operator to cluster
	./scripts/deploy-operator.sh

.PHONY: first-operator-undeploy
first-operator-undeploy: ## Undeploy first-operator from cluster
	cd first-operator && make undeploy

##@ Development

.PHONY: lint
lint: ## Run linters for all operators
	cd first-operator && make lint

.PHONY: fmt
fmt: ## Format code for all operators
	cd first-operator && go fmt ./...

.PHONY: clean
clean: ## Clean build artifacts
	cd first-operator && make clean || rm -rf bin/

##@ Quick Start

.PHONY: setup
setup: cluster-create first-operator-build ## Complete setup: create cluster and build operator
	@echo "Setup complete! Next steps:"
	@echo "  1. Run 'make first-operator-run' to run the operator locally, or"
	@echo "  2. Run 'make first-operator-deploy' to deploy to the cluster"

.PHONY: demo
demo: ## Run a complete demo (create cluster, deploy operator, create sample)
	@echo "Starting demo..."
	@make cluster-create
	@make first-operator-deploy
	@echo "Applying sample FirstApp..."
	kubectl apply -f first-operator/config/samples/apps_v1alpha1_firstapp.yaml
	@echo "Waiting for deployment..."
	@sleep 5
	@echo "\nFirstApp resources:"
	@kubectl get firstapps
	@echo "\nDeployments:"
	@kubectl get deployments
	@echo "\nPods:"
	@kubectl get pods
