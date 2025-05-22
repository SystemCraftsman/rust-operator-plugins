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

import "sigs.k8s.io/kubebuilder/v4/pkg/machinery"

var _ machinery.Template = &CargoToml{}

type CargoToml struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
}

func (f *CargoToml) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = "Cargo.toml"
	}

	f.TemplateBody = cargoTomlTemplate

	return nil
}

const cargoTomlTemplate = `[package]
name = "{{ .ProjectName }}"
version = "0.1.0"
edition = "2024"
rust-version = "1.87.0"

[[bin]]
name = "crdgen"
path = "src/crd_generator.rs"

[dependencies]
futures = "0.3.31"
k8s-openapi = { version = "0.24.0", features = ["latest"] }
kube = { version = "0.99.0", features = ["runtime", "client", "derive"] }
thiserror = "2.0.8"
tokio = { version = "1.42.0", features = ["macros", "rt-multi-thread", "rt"] }
schemars = "0.8.21"
serde = "1.0.216"
serde_json = "1.0.134"
serde_yaml = "0.9.34"
async-trait = "0.1.83"
`
