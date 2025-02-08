package rust

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds"
	"github.com/spf13/pflag"
	"log"
	"os"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin/util"
)

const (
	forceFlag      = "force"
	crdVersionFlag = "crd-version"
	namespacedFlag = "namespaced"
	resourceFlag   = "resource"
	controllerFlag = "controller"

	isForced              = false
	defaultCrdVersion     = "v1beta1"
	isNamespaced          = true
	isResourceAPICreation = true
	isControllerCreation  = true
)

// DefaultMainPath is default file path of main.go
const DefaultMainPath = "src/main.rs"

var _ plugin.CreateAPISubcommand = &createAPISubcommand{}

type createAPISubcommand struct {
	config   config.Config
	resource *resource.Resource
	options  *rust.Options

	// Check if we have to scaffold resource and/or controller
	resourceFlag   *pflag.Flag
	controllerFlag *pflag.Flag

	// force indicates that the resource should be created even if it already exists
	force bool
}

func (p *createAPISubcommand) UpdateMetadata(cliMeta plugin.CLIMetadata, subcmdMeta *plugin.SubcommandMetadata) {
	subcmdMeta.Description = `Scaffold a Kubernetes API by writing a Resource definition and/or a Controller.

If information about whether the resource and controller should be scaffolded
was not explicitly provided, it will prompt the user if they should be.

After the scaffold is written, the dependencies will be updated and
make generate will be run.
`
	subcmdMeta.Examples = fmt.Sprintf(`  # Create a frigates API with Group: ship, Version: v1 and Kind: Frigate
  %[1]s create api --group ship --version v1 --kind Frigate

  # Edit the API Scheme

  vim src/api/frigate_types.rs

  # Edit the Controller
  vim src/controller/frigate_controller.rs

  # Generate CRDs
  make generate-crds

  # Install CRDs into the Kubernetes cluster using kubectl apply
  make install

  # Regenerate code and run against the Kubernetes cluster configured by ~/.kube/config
  make run
`, cliMeta.CommandName)
}

func (p *createAPISubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&p.force, forceFlag, isForced,
		"attempt to create resource even if it already exists")

	p.options = &rust.Options{}

	fs.StringVar(&p.options.CRDVersion, crdVersionFlag, defaultCrdVersion, "crd version to generate")
	fs.BoolVar(&p.options.Namespaced, namespacedFlag, isNamespaced, "resource is namespaced")
	fs.BoolVar(&p.options.DoAPI, resourceFlag, isResourceAPICreation,
		"if set, generate the resource without prompting the user")
	p.resourceFlag = fs.Lookup(resourceFlag)
	fs.BoolVar(&p.options.DoController, controllerFlag, isControllerCreation,
		"if set, generate the controller without prompting the user")
	p.controllerFlag = fs.Lookup(controllerFlag)
}

func (p *createAPISubcommand) InjectConfig(c config.Config) error {
	p.config = c

	return nil
}

func (p *createAPISubcommand) InjectResource(res *resource.Resource) error {
	p.resource = res

	reader := bufio.NewReader(os.Stdin)
	if !p.resourceFlag.Changed {
		log.Println("Create Resource [y/n]")
		p.options.DoAPI = util.YesNo(reader)
	}
	if !p.controllerFlag.Changed {
		log.Println("Create Controller [y/n]")
		p.options.DoController = util.YesNo(reader)
	}

	p.options.UpdateResource(p.resource)

	if err := p.resource.Validate(); err != nil {
		return err
	}

	// In case we want to scaffold a resource API we need to do some checks
	if p.options.DoAPI {
		// Check that resource doesn't have the API scaffolded or flag force was set
		if r, err := p.config.GetResource(p.resource.GVK); err == nil && r.HasAPI() && !p.force {
			return errors.New("API resource already exists")
		}
	}

	return nil
}

func (p *createAPISubcommand) PreScaffold(machinery.Filesystem) error {
	// check if main.rs is present in the src/ directory
	if _, err := os.Stat(DefaultMainPath); os.IsNotExist(err) {
		return fmt.Errorf("%s file should present in the root directory", DefaultMainPath)
	}

	return nil
}

func (p *createAPISubcommand) Scaffold(fs machinery.Filesystem) error {
	scaffolder := scaffolds.NewAPIScaffolder(p.config, *p.resource, p.force)
	scaffolder.InjectFS(fs)
	return scaffolder.Scaffold()
}

func (p *createAPISubcommand) PostScaffold() error {
	err := util.RunCmd("Format code", "cargo", "fmt")
	if err != nil {
		return err
	}
	if p.resource.HasAPI() {
		// print follow on instructions to better guide the user
		fmt.Print("Next: implement your new API and generate the CRDs with:\n$ make generate-crds\n")
	}
	return nil
}
