GO=go

BUILD_DIR=bin

.PHONY: prepare
prepare:
	mkdir -p $(BUILD_DIR)

.PHONY: build
build: prepare proto tidy gofmt
	go build -o $(BUILD_DIR)/mrds-apiserver ./cmd/apiserver && \
	go build -o $(BUILD_DIR)/mrds-controlplane ./cmd/controlplane && \
	go build -o $(BUILD_DIR)/mrds-ctl ./cmd/ctl

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: gomodtidy
gomodtidy:
	go mod tidy

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -rf bin 2>/dev/null || true
	rm -rf gen 2>/dev/null || true
	rm go.sum 2>/dev/null || true

# Generate Go bindings from proto files
.PHONY: proto
proto:
	mkdir -p gen/api && \
	bash scripts/generate_proto.sh

.PHONE: run-controlplane
run-controlplane:
	go run cmd/controlplane/main.go

.PHONY: test
test:
	go test ./...

localbuild: clean build generate test

.PHONY: kind-prep
kind-prep: ## Create a kind cluster and add the nodes to MRDS
	kind create cluster --config test_manifests/kind_config.yaml && \
	./bin/mrds-ctl node create -m "test_manifests/node.yaml" && \
	./bin/mrds-ctl node add-to-cluster kind-worker --cluster-id kind && \
	./bin/mrds-ctl node add-to-cluster kind-worker2 --cluster-id kind && \
	./bin/mrds-ctl node add-to-cluster kind-worker3 --cluster-id kind && \
	./bin/mrds-ctl node add-to-cluster kind-worker4 --cluster-id kind

.PHONY: kind-delete
kind-delete: ## Delete the kind cluster
	kind delete cluster

.PHONY: add-deployment
add-deployment: ## Add a deployment to the cluster
	./bin/mrds-ctl deployment create -m test_manifests/deploymentplan.yaml && \
	./bin/mrds-ctl deployment add-deployment -m test_manifests/deployment.yaml

