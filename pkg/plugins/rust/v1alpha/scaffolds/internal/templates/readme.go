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
	"strings"
)

var _ machinery.Template = &Readme{}

// Readme scaffolds a README.md file
type Readme struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
	machinery.BoilerplateMixin

	License string
}

// SetTemplateDefaults implements file.Template
func (f *Readme) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = "README.md"
	}

	f.License = strings.Replace(
		strings.Replace(f.Boilerplate, "/*", "", 1),
		"*/", "", 1)

	f.TemplateBody = fmt.Sprintf(readmeFileTemplate,
		codeFence("make build"),
		codeFence("make run"),
		codeFence("make image-build image-push IMG=<some-registry>/{{ .ProjectName }}:tag"),
		codeFence("make generate-crds"),
		codeFence("make install"),
		codeFence("make deploy IMG=<some-registry>/{{ .ProjectName }}:tag"),
		codeFence("kubectl apply -k path/to/your/samples/"),
		codeFence("kubectl delete -k path/to/your/samples/"),
		codeFence("make uninstall"),
		codeFence("make undeploy"),
	)

	return nil
}

//nolint:lll
const readmeFileTemplate = `# {{ .ProjectName }}

// TODO(user): Add simple overview of use/purpose

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

### Prerequisites

- cargo version 1.86.0
- docker version 27.5.0+
- kubectl version v1.32.1+.
- Access to a Kubernetes v1.25.3+ cluster.

### To Run locally

**Build your operator:**

%s

**Run your operator:**

%s

### To Deploy on the cluster

**Build and push your image to the location specified by ` + "`IMG`" + `:**

%s

> **NOTE:** This image ought to be published in the personal registry you specified.
> And it is required to have access to pull the image from the working environment.
> Make sure you have the proper permission to the registry if the above commands don’t work.

**Generate the CRDs:**

%s

**Install the CRDs into the cluster:**

%s

**Deploy the operator to the cluster with the image specified by ` + "`IMG`" + `:**

%s

> **IMPORTANT**: You will face API access errors, as this script only creates the deployment.
> You will need to create the required role bindings for your deployment for now.

**Create instances of your solution**
You can apply your example CRs:

%s

> **IMPORTANT**: Ensure that the samples has default values to test it out.

### To Uninstall

**Delete the instances (CRs) from the cluster:**

%s

**Delete the APIs(CRDs) from the cluster:**

%s

**UnDeploy the controller from the cluster:**

%s

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run ` + "`make help`" + ` for more information on all potential ` + "`make`" + ` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

// TODO(user): Add a license
`

func codeFence(code string) string {
	return "```sh" + "\n" + code + "\n" + "```"
}
