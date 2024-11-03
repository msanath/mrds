# Variables
GO := go
BUILD_DIR := bin
PROTO_GEN_DIR := gen/api
PROTO_SCRIPT := scripts/generate_proto.sh
KIND_CONFIG := test/manifests/kind_config.yaml
NODE_CONFIG := test/manifests/node.yaml
DEPLOYMENT_PLAN := test/manifests/deploymentplan.yaml
DEPLOYMENT_CONFIG := test/manifests/deployment-1.yaml
KIND_CLUSTER_NAME := kind
DB_USER := root
DB_NAME := mrds

# Define nodes to add to the kind cluster
KIND_NODES := kind-worker kind-worker2 kind-worker3 kind-worker4

# Default target
.DEFAULT_GOAL := build

# Prepare build directory
.PHONY: prepare
prepare: ## Create the build directory if it doesn't exist
	mkdir -p $(BUILD_DIR)

# Build all binaries
.PHONY: build
build: prepare proto tidy gofmt ## Compile all services and binaries
	$(GO) build -o $(BUILD_DIR)/mrds-apiserver ./cmd/apiserver
	$(GO) build -o $(BUILD_DIR)/mrds-controlplane ./cmd/controlplane
	$(GO) build -o $(BUILD_DIR)/mrds-ctl ./cmd/ctl

# Go module cleanup
.PHONY: tidy
tidy: ## Clean up go.mod and go.sum files
	$(GO) mod tidy

# Format Go code
.PHONY: gofmt
gofmt: ## Format Go code in all packages
	$(GO) fmt ./...

# Clean up generated files and binaries
.PHONY: clean
clean: ## Remove generated files and binaries
	rm -rf $(BUILD_DIR) $(PROTO_GEN_DIR) go.sum

# Generate Go bindings from proto files
.PHONY: proto
proto: ## Generate Go bindings from protobuf definitions
	mkdir -p $(PROTO_GEN_DIR)
	bash $(PROTO_SCRIPT)

# Run control plane locally
.PHONY: run-controlplane
run-controlplane: ## Run the control plane locally
	./bin/mrds-controlplane

.PHONY: run-apiserver
run-apiserver: ## Run the API server locally
	./bin/mrds-apiserver

.PHONY: run-temporal
run-temporal: ## Run the Temporal server locally
	temporal server start-dev --namespace mrds

# Run tests
.PHONY: test
test: ## Run unit tests
	$(GO) test ./...

# Clean, build, generate proto files, and test
.PHONY: localbuild
localbuild: clean build proto test ## Clean, build, generate proto files, and test

# Prepare a Kind cluster and add nodes (for testing)
.PHONY: test-kind-prep
test-kind-prep: test-kind-delete ## Create a kind cluster and add specified nodes (for testing)
	kind create cluster --config $(KIND_CONFIG)
	./bin/mrds-ctl node create -m "$(NODE_CONFIG)"
	$(foreach node,$(KIND_NODES),./bin/mrds-ctl node add-to-cluster $(node) --cluster-id $(KIND_CLUSTER_NAME);)

# Delete Kind cluster (for testing)
.PHONY: test-kind-delete
test-kind-delete: ## Delete the Kind cluster (for testing)
	kind delete cluster --name $(KIND_CLUSTER_NAME)

# Add a deployment to the cluster
.PHONY: add-deployment
add-deployment: ## Create a deployment in the cluster
	./bin/mrds-ctl deployment create -m $(DEPLOYMENT_PLAN)
	./bin/mrds-ctl deployment add-deployment -m $(DEPLOYMENT_CONFIG)

# Display help for each target
.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Reset the database
.PHONY: local-db-reset
local-db-reset:
	@echo "Dropping database..."
	mysql -u $(DB_USER) -e "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "Creating database..."
	mysql -u $(DB_USER) -e "CREATE DATABASE $(DB_NAME);"
