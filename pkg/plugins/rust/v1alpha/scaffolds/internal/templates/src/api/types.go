package api

import (
	"log"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

var _ machinery.Template = &Types{}

// Types scaffolds the file that defines the schema for a CRD
// nolint:maligned
type Types struct {
	machinery.TemplateMixin
	machinery.ResourceMixin
	machinery.BoilerplateMixin

	Force bool
}

func (f *Types) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("src", "api", "%[kind]_types.rs")
	}

	f.Path = f.Resource.Replacer().Replace(f.Path)
	log.Println(f.Path)

	f.TemplateBody = typesTemplate

	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		f.IfExistsAction = machinery.Error
	}

	return nil
}

const typesTemplate = `{{ .Boilerplate }}

use k8s_openapi::serde::{Deserialize, Serialize};
use kube::CustomResource;
use schemars::JsonSchema;

#[derive(CustomResource, Deserialize, Serialize, Clone, Debug, JsonSchema)]
#[kube(
    kind = "{{ .Resource.Kind }}",
    group = "{{ .Resource.Group }}",
    version = "{{ .Resource.Version }}",
    namespaced
	status = "{{ .Resource.Kind }}Status"
)]
pub struct {{ .Resource.Kind }}Spec {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster

	// foo is an example field of {{ .Resource.Kind }}. Edit {{ lower .Resource.Kind }}_types.rs to remove/update
    foo: String,
}

#[derive(Deserialize, Serialize, Clone, Debug, JsonSchema)]
pub struct {{ .Resource.Kind }}Status {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
}
`
