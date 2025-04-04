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

package rust

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
)

// Options contains the information required to build a new resource.Resource.
type Options struct {
	Namespaced   bool
	DoAPI        bool
	DoController bool
}

// UpdateResource updates the provided resource with the options
func (opts Options) UpdateResource(res *resource.Resource) {
	if opts.DoAPI {
		res.Path = ""

		res.API = &resource.API{
			CRDVersion: "v1",
			Namespaced: opts.Namespaced,
		}

	}

	if opts.DoController {
		res.Controller = true
	}
}
