package rust

import (
	"errors"
	"fmt"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/rust/v1alpha/scaffolds"
	"github.com/spf13/pflag"

	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

const (
	crdVersionFlag = "crd-version"
	namespacedFlag = "namespaced"

	defaultCrdVersion      = "v1"
	defaultNamespacedValue = true
)

type createAPIOptions struct {
	CRDVersion string
	Namespaced bool
}

type createAPISubcommand struct {
	config   config.Config
	resource *resource.Resource
	options  createAPIOptions
}

func (opts createAPIOptions) UpdateResource(res *resource.Resource) {

	res.API = &resource.API{
		CRDVersion: opts.CRDVersion,
		Namespaced: opts.Namespaced,
	}

	res.Path = ""
	res.Controller = false
}

var (
	_ plugin.CreateAPISubcommand = &createAPISubcommand{}
)

func (p *createAPISubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false
	fs.StringVar(&p.options.CRDVersion, crdVersionFlag, defaultCrdVersion, "crd version to generate")
	fs.BoolVar(&p.options.Namespaced, namespacedFlag, defaultNamespacedValue, "resource is namespaced")
}

func (p *createAPISubcommand) InjectConfig(c config.Config) error {
	p.config = c

	return nil
}

func (p *createAPISubcommand) Run(fs machinery.Filesystem) error {
	return nil
}

func (p *createAPISubcommand) Validate() error {
	return nil
}

func (p *createAPISubcommand) PostScaffold() error {
	return nil
}

func (p *createAPISubcommand) Scaffold(fs machinery.Filesystem) error {
	scaffolder := scaffolds.NewCreateAPIScaffolder(p.config, *p.resource)
	scaffolder.InjectFS(fs)

	if err := scaffolder.Scaffold(); err != nil {
		return err
	}

	return nil
}

func (p *createAPISubcommand) InjectResource(res *resource.Resource) error {
	p.resource = res

	p.options.UpdateResource(p.resource)

	if err := p.resource.Validate(); err != nil {
		return err
	}

	// Check that resource doesn't have the API scaffolded
	if res, err := p.config.GetResource(p.resource.GVK); err == nil && res.HasAPI() {
		return errors.New("the API resource already exists")
	}

	// Check that the provided group can be added to the project
	if !p.config.IsMultiGroup() && p.config.ResourcesLength() != 0 && !p.config.HasGroup(p.resource.Group) {
		return fmt.Errorf("multiple groups are not allowed by default, to enable multi-group set 'multigroup: true' in your PROJECT file")
	}

	// Selected CRD version must match existing CRD versions.
	if p.hasDifferentCRDVersion(p.config, p.resource.API.CRDVersion) {
		return fmt.Errorf("only one CRD version can be used for all resources, cannot add %q", p.resource.API.CRDVersion)
	}

	return nil
}

// hasDifferentCRDVersion returns true if any other CRD version is tracked in the project configuration.
func (p *createAPISubcommand) hasDifferentCRDVersion(config config.Config, crdVersion string) bool {
	return hasDifferentAPIVersion(config.ListCRDVersions(), crdVersion)
}

func hasDifferentAPIVersion(versions []string, version string) bool {
	return !(len(versions) == 0 || (len(versions) == 1 && versions[0] == version))
}
