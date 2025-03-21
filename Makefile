# GO_BUILD_ARGS should be set when running 'go build' or 'go install'.
VERSION_PKG = "$(shell go list -m)/internal/version"
export GIT_VERSION = $(shell git describe --dirty --tags --always)
export GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_DIR = $(PWD)/bin
GO_BUILD_ARGS = \
  -gcflags "all=-trimpath=$(shell dirname $(shell pwd))" \
  -asmflags "all=-trimpath=$(shell dirname $(shell pwd))" \
  -ldflags " \
    -s \
    -w \
    -X '$(VERSION_PKG).GitVersion=$(GIT_VERSION)' \
    -X '$(VERSION_PKG).GitCommit=$(GIT_COMMIT)' \
  " \

# Always use Go modules
export GO111MODULE = on

# This is to allow for building and testing on Apple Silicon.
# These values default to the host's GOOS and GOARCH, but should
# be overridden when running builds and tests on Apple Silicon unless
# you are only building the binary
BUILD_GOOS ?= $(shell go env GOOS)
BUILD_GOARCH ?= $(shell go env GOARCH)

OPERATOR_SDK_REPO_PATH ?= https://github.com/mabulgu/operator-sdk
OPERATOR_SDK_BRANCH ?= rust-operator
OPERATOR_SDK_DIR_NAME ?= operator-sdk
OPERATOR_SDK_BIN_NAME ?= operator-sdk

##@ Development

.PHONY: lint
lint:
	@./hack/check-license.sh
	@go fmt ./...

##@ Test

.PHONY: test
test:
	@go test -coverprofile=coverage.out -covermode=count -short ./...

##@ Build

.PHONY: download-sdk
download-sdk: ## Download Operator SDK
	rm -rf $(OPERATOR_SDK_DIR_NAME)
	git clone -b $(OPERATOR_SDK_BRANCH) $(OPERATOR_SDK_REPO_PATH)

.PHONY: configure-sdk
configure-sdk: download-sdk ## Configure Operator SDK
	cd $(OPERATOR_SDK_DIR_NAME) && go mod edit -replace github.com/SystemCraftsman/rust-operator-plugins=..
	cd $(OPERATOR_SDK_DIR_NAME) && go mod edit -replace sigs.k8s.io/kubebuilder/v4=github.com/mabulgu/kubebuilder/v4@rust-lang
	cd $(OPERATOR_SDK_DIR_NAME) && go mod tidy

.PHONY: build
build: ## Build plugin with Operator SDK
	cd $(OPERATOR_SDK_DIR_NAME) && make $@
	cp $(OPERATOR_SDK_DIR_NAME)/build/$(OPERATOR_SDK_BIN_NAME) $(BUILD_DIR)

.PHONY: install
install: ## Install Operator SDK CLI
	cd $(OPERATOR_SDK_DIR_NAME) && make $@
