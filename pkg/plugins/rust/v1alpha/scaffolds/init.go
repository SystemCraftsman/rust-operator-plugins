package scaffolds

import (
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates/src"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins"
)

var _ plugins.Scaffolder = &initScaffolder{}

type initScaffolder struct {
	config config.Config

	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem
}

// NewInitScaffolder returns a new plugins.Scaffolder for project initialization operations
func NewInitScaffolder(config config.Config) plugins.Scaffolder {
	return &initScaffolder{
		config: config,
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
		&src.CRDGenerator{},
		&templates.CargoToml{},
		&templates.GitIgnore{},
		&templates.Makefile{},
		&templates.Dockerfile{},
		&templates.DockerIgnore{},
		&templates.Readme{},
	)
}
