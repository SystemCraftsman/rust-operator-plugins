package scaffolds

import (
	"fmt"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates/src"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates/src/api"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates/src/controller"
	"log"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins"
)

var _ plugins.Scaffolder = &apiScaffolder{}

// apiScaffolder contains configuration for generating scaffolding for Rust type
// representing the API and controller that implements the behavior for the API.
type apiScaffolder struct {
	config   config.Config
	resource resource.Resource

	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem

	// force indicates whether to scaffold controller files even if it exists or not
	force bool
}

// NewAPIScaffolder returns a new Scaffolder for API/controller creation operations
func NewAPIScaffolder(config config.Config, res resource.Resource, force bool) plugins.Scaffolder {
	return &apiScaffolder{
		config:   config,
		resource: res,
		force:    force,
	}
}

func (s *apiScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

// Scaffold implements cmdutil.Scaffolder
func (s *apiScaffolder) Scaffold() error {
	log.Println("Writing scaffold for you to edit...")

	// Initialize the machinery.Scaffold that will write the files to disk
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
		machinery.WithResource(&s.resource),
	)

	// Keep track of these values before the update
	doAPI := s.resource.HasAPI()
	doController := s.resource.HasController()

	if err := s.config.UpdateResource(s.resource); err != nil {
		return fmt.Errorf("error updating resource: %w", err)
	}

	if doAPI {
		if err := scaffold.Execute(
			&api.Types{Force: s.force},
		); err != nil {
			return fmt.Errorf("error scaffolding APIs: %v", err)
		}

		if err := scaffold.Execute(
			&src.ApiUpdater{WireResource: doAPI, WireController: doController},
		); err != nil {
			return fmt.Errorf("error updating src/api.rs: %v", err)
		}

		if err := scaffold.Execute(
			&src.CRDGeneratorUpdater{WireResource: doAPI, WireController: doController},
		); err != nil {
			return fmt.Errorf("error updating src/crd_generator.rs: %v", err)
		}
	}

	if doController {
		if err := scaffold.Execute(
			&controller.Controllers{Force: s.force},
		); err != nil {
			return fmt.Errorf("error scaffolding controller: %v", err)
		}

		if err := scaffold.Execute(
			&src.ControllerUpdater{WireResource: doAPI, WireController: doController},
		); err != nil {
			return fmt.Errorf("error updating src/controller.rs: %v", err)
		}

		if err := scaffold.Execute(
			&src.MainUpdater{WireResource: doAPI, WireController: doController},
		); err != nil {
			return fmt.Errorf("error updating src/main.rs: %v", err)
		}
	}

	return nil
}
