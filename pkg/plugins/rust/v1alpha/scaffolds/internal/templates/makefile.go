/*
Copyright 2025 System Craftsman LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package templates

import (
	"fmt"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

var _ machinery.Template = &Makefile{}

type Makefile struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin

	Image string
}

func (f *Makefile) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = "Makefile"
	}

	f.TemplateBody = makefileTemplate

	f.IfExistsAction = machinery.Error

	if f.Image == "" {
		f.Image = fmt.Sprintf("%s:latest", f.ProjectName)
	}

	return nil
}

const makefileTemplate = `# Image URL to use for all building/pushing image targets
IMG ?= {{ .Image }}

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

NOT-IMPLEMENTED:
	@echo
	@echo [WARN] This target is not yet implemented.
	@echo

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: generate-crds
generate-crds:
	 cargo run --bin crdgen

##@ Build

.PHONY: build
build: ## Build operator binary.
	cargo build

.PHONY: run
run:  ## Run operator from your host.
	cargo run --package {{ .ProjectName }} --bin {{ .ProjectName }}

.PHONY: image-build
image-build: ## Build docker image.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: image-push
image-push: ## Push container image.
	$(CONTAINER_TOOL) push ${IMG}

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	@$(foreach file, $(wildcard target/kubernetes/*-v1alpha1.yaml), kubectl apply -f $(file);)

.PHONY: uninstall
uninstall: ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	@$(foreach file, $(wildcard target/kubernete s/*-v1alpha1.yaml), kubectl delete -f $(file) --ignore-not-found=$(ignore-not-found);)

.PHONY: deploy
deploy: ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	kubectl create deployment {{ .ProjectName }} --image ${IMG}

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	kubectl delete deployment {{ .ProjectName }} --ignore-not-found=$(ignore-not-found)
`
