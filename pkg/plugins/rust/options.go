package rust

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
)

// Options contains the information required to build a new resource.Resource.
type Options struct {
	CRDVersion   string
	Namespaced   bool
	DoAPI        bool
	DoController bool
}

// UpdateResource updates the provided resource with the options
func (opts Options) UpdateResource(res *resource.Resource) {
	if opts.DoAPI {
		res.Path = ""

		res.API = &resource.API{
			CRDVersion: opts.CRDVersion,
			Namespaced: opts.Namespaced,
		}

	}

	if opts.DoController {
		res.Controller = true
	}
}
