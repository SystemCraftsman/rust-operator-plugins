package scaffolds

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins"
)

type apiScaffolder struct {
	fs machinery.Filesystem

	config   config.Config
	resource resource.Resource
}

func NewCreateAPIScaffolder(cfg config.Config, res resource.Resource) plugins.Scaffolder {
	return &apiScaffolder{
		config:   cfg,
		resource: res,
	}
}

func (s *apiScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

func (s *apiScaffolder) Scaffold() error {

	if err := s.config.UpdateResource(s.resource); err != nil {
		return err
	}

	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithDirectoryPermissions(0755),
		machinery.WithFilePermissions(0644),
		machinery.WithConfig(s.config),
		machinery.WithResource(&s.resource),
	)

	var createAPITemplates []machinery.Builder
	//TODO: Append API templates here
	createAPITemplates = append(createAPITemplates)

	return scaffold.Execute(createAPITemplates...)
}
