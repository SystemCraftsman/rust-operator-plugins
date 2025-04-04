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
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/constants"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"strings"
)

const (
	defaultApiPath = "src/api.rs"
)

var _ machinery.Template = &Api{}

type Api struct {
	machinery.TemplateMixin
	machinery.BoilerplateMixin
}

// SetTemplateDefaults implements file.Template
func (f *Api) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join(defaultApiPath)
	}

	f.TemplateBody = fmt.Sprintf(apiTemplate,
		machinery.NewMarkerFor(f.Path, constants.ModuleMarker),
	)

	return nil
}

var _ machinery.Inserter = &ApiUpdater{}

type ApiUpdater struct { //nolint:maligned
	machinery.ResourceMixin

	// Flags to indicate which parts need to be included when updating the file
	WireResource, WireController bool
}

// GetPath implements file.Builder
func (*ApiUpdater) GetPath() string {
	return defaultApiPath
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

	// Generate module code fragments
	modules := make([]string, 0)
	if f.WireResource {
		modules = append(modules, fmt.Sprintf(moduleImportCodeFragment, strings.ToLower(f.Resource.Kind)))
	}

	// Only store code fragments in the map if the slices are non-empty
	if len(modules) != 0 {
		fragments[machinery.NewMarkerFor(defaultApiPath, constants.ModuleMarker)] = modules
	}

	return fragments
}

// nolint:lll
var apiTemplate = `{{ .Boilerplate }}

%s
`
