package src

import (
	"fmt"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"strings"
)

const (
	writerMarker = "writers"

	defaultCRDGeneratorPath = "src/crd_generator.rs"
)

var _ machinery.Template = &CRDGenerator{}

type CRDGenerator struct {
	machinery.TemplateMixin
}

// SetTemplateDefaults implements file.Template
func (f *CRDGenerator) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join(defaultCRDGeneratorPath)
	}

	f.TemplateBody = fmt.Sprintf(crdGeneratorTemplate,
		machinery.NewMarkerFor(f.Path, writerMarker),
	)

	return nil
}

var _ machinery.Inserter = &CRDGeneratorUpdater{}

type CRDGeneratorUpdater struct { //nolint:maligned
	machinery.ResourceMixin

	// Flags to indicate which parts need to be included when updating the file
	WireResource, WireController bool
}

// GetPath implements file.Builder
func (*CRDGeneratorUpdater) GetPath() string {
	return defaultCRDGeneratorPath
}

// GetIfExistsAction implements file.Builder
func (*CRDGeneratorUpdater) GetIfExistsAction() machinery.IfExistsAction {
	return machinery.OverwriteFile
}

// GetMarkers implements file.Inserter
func (f *CRDGeneratorUpdater) GetMarkers() []machinery.Marker {
	return []machinery.Marker{
		machinery.NewMarkerFor(defaultCRDGeneratorPath, writerMarker),
	}
}

const (
	writerCodeFragment = `write_crd_to_yaml(&api::%s_types::%s::crd());
`
)

// GetCodeFragments implements file.Inserter
func (f *CRDGeneratorUpdater) GetCodeFragments() machinery.CodeFragmentsMap {
	fragments := make(machinery.CodeFragmentsMap, 3)

	// If resource is not being provided we are creating the file, not updating it
	if f.Resource == nil {
		return fragments
	}

	// Generate writer code fragments
	writers := make([]string, 0)
	if f.WireController {
		writers = append(writers, fmt.Sprintf(writerCodeFragment, strings.ToLower(f.Resource.Kind), f.Resource.Kind))
	}

	// Only store code fragments in the map if the slices are non-empty
	if len(writers) != 0 {
		fragments[machinery.NewMarkerFor(defaultCRDGeneratorPath, writerMarker)] = writers
	}

	return fragments
}

// nolint:lll
var crdGeneratorTemplate = `mod api;

use k8s_openapi::apiextensions_apiserver::pkg::apis::apiextensions::v1::CustomResourceDefinition;
use kube::CustomResourceExt;
use std::fs;
use std::fs::File;

fn main() {
    fs::create_dir_all("target/kubernetes").expect("Error creating directory 'target/kubernetes'");
    %s
}

fn write_crd_to_yaml(crd: &CustomResourceDefinition) {
    let file_path = format!(
        "target/kubernetes/{name}-{version}.yaml",
        name = crd.metadata.name.clone().unwrap(),
        version = crd.spec.versions.first().unwrap().name
    );
    let file = File::create(file_path).expect("Error creating YAML file");
    serde_yaml::to_writer(file, crd).expect("Error writing to YAML file");
}
`
