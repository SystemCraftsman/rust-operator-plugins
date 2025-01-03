package scaffolds

import (
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates/src"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins"
)

const (
	// ControllerRuntimeVersion is the kubernetes-sigs/controller-runtime version to be used in the project
	ControllerRuntimeVersion = "v0.19.1"

	// kustomizeVersion is the sigs.k8s.io/kustomize version to be used in the project
	kustomizeVersion = "v3.5.4"

	imageName = "controller:latest"
)

var _ plugins.Scaffolder = &initScaffolder{}

type initScaffolder struct {
	config          config.Config
	boilerplatePath string
	license         string
	owner           string

	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem
}

// NewInitScaffolder returns a new plugins.Scaffolder for project initialization operations
func NewInitScaffolder(config config.Config, license, owner string) plugins.Scaffolder {
	return &initScaffolder{
		config: config,
		//boilerplatePath: hack.DefaultBoilerplatePath,
		license: license,
		owner:   owner,
	}
}

// InjectFS implements Scaffolder
func (s *initScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

// Scaffold implements Scaffolder
func (s *initScaffolder) Scaffold() error {
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
	)

	return scaffold.Execute(
		&src.Main{},
		&src.Api{},
		&src.Controller{},
		&templates.CargoToml{},
		&templates.GitIgnore{},
		//TODO makefile
		//TODO readme
	)
}
