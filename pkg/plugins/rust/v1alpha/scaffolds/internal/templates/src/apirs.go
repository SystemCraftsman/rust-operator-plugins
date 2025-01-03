package src

import (
	"fmt"
	localmachinery "github.com/SystemCraftsman/rust-operator-plugins/pkg/machinery"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/constants"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

const (
	defaultApiPath = "src/api.rs"
)

var _ machinery.Template = &Api{}

type Api struct {
	machinery.TemplateMixin
}

// SetTemplateDefaults implements file.Template
func (f *Api) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join(defaultApiPath)
	}

	f.TemplateBody = fmt.Sprintf(apiTemplate,
		localmachinery.NewMarkerFor(f.Path, constants.ModuleMarker),
	)

	return nil
}

var _ machinery.Inserter = &ApiUpdater{}

type ApiUpdater struct { //nolint:maligned
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
func (*ApiUpdater) GetPath() string {
	return defaultMainPath
}

// GetIfExistsAction implements file.Builder
func (*ApiUpdater) GetIfExistsAction() machinery.IfExistsAction {
	return machinery.OverwriteFile
}

// GetMarkers implements file.Inserter
func (f *ApiUpdater) GetMarkers() []machinery.Marker {
	return []machinery.Marker{
		machinery.NewMarkerFor(defaultApiPath, constants.ModuleMarker),
	}
}

const (
	moduleImportCodeFragment = `pub mod %s_types;
`
)

// GetCodeFragments implements file.Inserter
func (f *ApiUpdater) GetCodeFragments() machinery.CodeFragmentsMap {
	fragments := make(machinery.CodeFragmentsMap, 3)

	// If resource is not being provided we are creating the file, not updating it
	if f.Resource == nil {
		return fragments
	}

	// Generate module import code fragments
	imports := make([]string, 0)
	if f.WireController {
		imports = append(imports, fmt.Sprintf(moduleImportCodeFragment, f.Resource.Kind))
	}

	// Only store code fragments in the map if the slices are non-empty
	if len(imports) != 0 {
		fragments[machinery.NewMarkerFor(defaultApiPath, constants.ModuleMarker)] = imports
	}

	return fragments
}

// nolint:lll
var apiTemplate = `%s
`
