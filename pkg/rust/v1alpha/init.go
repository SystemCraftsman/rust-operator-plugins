package rust

import (
	"fmt"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/rust/v1alpha/scaffolds"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation"

	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

// This file represents the CLI for this plugin.

const (
	groupFlag   = "group"
	versionFlag = "version"
	kindFlag    = "kind"
)

type initSubcommand struct {
	apiSubcommand createAPISubcommand

	config config.Config

	// For help text.
	commandName string

	// Flags
	group       string
	domain      string
	version     string
	kind        string
	projectName string
}

var (
	_ plugin.InitSubcommand = &initSubcommand{}
)

func (p *initSubcommand) UpdateMetadata(cliMeta plugin.CLIMetadata, subcmdMeta *plugin.SubcommandMetadata) {
	subcmdMeta.Description = `Initialize a new project based on the java-operator-sdk project.

Writes the following files:
- a basic, Quarkus-based operator set-up
- a pom.xml file to build the project with Maven
`
	p.commandName = cliMeta.CommandName
}

func (p *initSubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false
	fs.StringVar(&p.domain, "domain", "my.domain", "domain for groups")
	fs.StringVar(&p.projectName, "project-name", "", "name of this project, the default being directory name")

	fs.StringVar(&p.group, groupFlag, "", "resource Group")
	fs.StringVar(&p.version, versionFlag, "", "resource Version")
	fs.StringVar(&p.kind, kindFlag, "", "resource Kind")
	p.apiSubcommand.BindFlags(fs)
}

func (p *initSubcommand) InjectConfig(c config.Config) error {
	p.config = c

	if err := p.config.SetDomain(p.domain); err != nil {
		return err
	}

	// Assign a default project name
	if p.projectName == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %v", err)
		}
		p.projectName = strings.ToLower(filepath.Base(dir))
	}
	// Check if the project name is a valid k8s namespace (DNS 1123 label).
	if err := validation.IsDNS1123Label(p.projectName); err != nil {
		return fmt.Errorf("project name (%s) is invalid: %v", p.projectName, err)
	}
	if err := p.config.SetProjectName(p.projectName); err != nil {
		return err
	}

	return nil
}

func (p *initSubcommand) Validate() error {
	// TODO: validate the conditions you expect before running the plugin
	return nil
}

func (p *initSubcommand) PostScaffold() error {
	// print follow on instructions to better guide the user
	fmt.Printf("Next: define a resource with:\n$ %s create api\n", p.commandName)
	return nil
}

func (p *initSubcommand) Scaffold(fs machinery.Filesystem) error {
	scaffolder := scaffolds.NewInitScaffolder(p.config)
	scaffolder.InjectFS(fs)
	return scaffolder.Scaffold()
}
