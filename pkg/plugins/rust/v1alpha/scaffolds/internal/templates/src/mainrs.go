package src

import (
	"fmt"
	localmachinery "github.com/SystemCraftsman/rust-operator-plugins/pkg/machinery"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
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
}

// SetTemplateDefaults implements file.Template
func (f *Main) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join(defaultMainPath)
	}

	f.TemplateBody = fmt.Sprintf(mainTemplate,
		localmachinery.NewMarkerFor(f.Path, importMarker),
		localmachinery.NewMarkerFor(f.Path, runnerMarker),
	)

	return nil
}

var _ machinery.Inserter = &MainUpdater{}

// MainUpdater updates cmd/main.go to run Controllers
type MainUpdater struct { //nolint:maligned
	machinery.ResourceMixin

	// Flags to indicate which parts need to be included when updating the file
	WireResource, WireController, WireWebhook bool

	// Deprecated - The flag should be removed from go/v5
	// IsLegacyPath indicates if webhooks should be scaffolded under the API.
	// Webhooks are now decoupled from APIs based on controller-runtime updates and community feedback.
	// This flag ensures backward compatibility by allowing scaffolding in the legacy/deprecated path.
	IsLegacyPath bool
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
		machinery.NewMarkerFor(defaultMainPath, importMarker),
		machinery.NewMarkerFor(defaultMainPath, runnerMarker),
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
		imports = append(imports, fmt.Sprintf(reconcilerImportCodeFragment, f.Resource.Kind, f.Resource.Kind))
	}

	// Generate setup code fragments
	setup := make([]string, 0)
	if f.WireController {
		setup = append(setup, fmt.Sprintf(reconcilerSetupCodeFragment, f.Resource.Kind))
	}

	// Only store code fragments in the map if the slices are non-empty
	if len(imports) != 0 {
		fragments[machinery.NewMarkerFor(defaultMainPath, importMarker)] = imports
	}
	if len(setup) != 0 {
		fragments[machinery.NewMarkerFor(defaultMainPath, runnerMarker)] = setup
	}

	return fragments
}

// nolint:lll
var mainTemplate = `mod api;
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