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

package src

import (
	"fmt"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"strings"
)

const (
	importMarker = "imports"
	runnerMarker = "runners"

	defaultMainPath = "src/main.rs"
)

var _ machinery.Template = &Main{}

// Main scaffolds a file that defines the controller manager entry point
type Main struct {
	machinery.TemplateMixin
	machinery.BoilerplateMixin
}

// SetTemplateDefaults implements file.Template
func (f *Main) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join(defaultMainPath)
	}

	f.TemplateBody = fmt.Sprintf(mainTemplate,
		rust.NewMarkerFor(f.Path, importMarker),
		rust.NewMarkerFor(f.Path, runnerMarker),
	)

	return nil
}

var _ machinery.Inserter = &MainUpdater{}

// MainUpdater updates src/main.rs to add reconcilers
type MainUpdater struct { //nolint:maligned
	machinery.ResourceMixin

	// Flags to indicate which parts need to be included when updating the file
	WireResource, WireController bool
}

// GetPath implements file.Builder
func (*MainUpdater) GetPath() string {
	return defaultMainPath
}

// GetIfExistsAction implements file.Builder
func (*MainUpdater) GetIfExistsAction() machinery.IfExistsAction {
	return machinery.OverwriteFile
}

// GetMarkers implements file.Inserter
func (f *MainUpdater) GetMarkers() []machinery.Marker {
	return []machinery.Marker{
		rust.NewMarkerFor(defaultMainPath, importMarker),
		rust.NewMarkerFor(defaultMainPath, runnerMarker),
	}
}

const (
	reconcilerImportCodeFragment = `use crate::controller::%s_controller::%sReconciler;
`
	reconcilerSetupCodeFragment = `tokio::spawn(async {
		ControllerRunner::run::<%sReconciler>().await;
	}),
`
)

// GetCodeFragments implements file.Inserter
func (f *MainUpdater) GetCodeFragments() machinery.CodeFragmentsMap {
	fragments := make(machinery.CodeFragmentsMap, 3)

	// If resource is not being provided we are creating the file, not updating it
	if f.Resource == nil {
		return fragments
	}

	// Generate import code fragments
	imports := make([]string, 0)
	if f.WireController {
		imports = append(imports, fmt.Sprintf(reconcilerImportCodeFragment, strings.ToLower(f.Resource.Kind), f.Resource.Kind))
	}

	// Generate setup code fragments
	setup := make([]string, 0)
	if f.WireController {
		setup = append(setup, fmt.Sprintf(reconcilerSetupCodeFragment, f.Resource.Kind))
	}

	// Only store code fragments in the map if the slices are non-empty
	if len(imports) != 0 {
		fragments[rust.NewMarkerFor(defaultMainPath, importMarker)] = imports
	}
	if len(setup) != 0 {
		fragments[rust.NewMarkerFor(defaultMainPath, runnerMarker)] = setup
	}

	return fragments
}

// nolint:lll
var mainTemplate = `{{ .Boilerplate }}

mod api;
mod controller;

use crate::controller::ControllerRunner;
%s

#[tokio::main]
async fn main() {
    let _ = tokio::join!(
        %s
    );
}
`
