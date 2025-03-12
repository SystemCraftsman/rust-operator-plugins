package rust

import (
	"fmt"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation"

	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

// This file represents the CLI for this plugin.
type initSubcommand struct {
	config config.Config

	// For help text.
	commandName string

	// Flags
	domain      string
	version     string
	projectName string
}

var _ plugin.InitSubcommand = &initSubcommand{}

func (p *initSubcommand) UpdateMetadata(cliMeta plugin.CLIMetadata, subcmdMeta *plugin.SubcommandMetadata) {
	p.commandName = cliMeta.CommandName

	subcmdMeta.Description = `Initialize a new project including the following files:
  - a "PROJECT" file that stores project configuration
  - a "Makefile" with several useful make targets for the project
`
	subcmdMeta.Examples = fmt.Sprintf(`  # Initialize a new project with your domain name
  %[1]s init --plugins rust/v1alpha --domain example.org

  # Initialize a new project defining a specific project version
  %[1]s init --plugins rust/v1alpha --project-version 3
`, cliMeta.CommandName)
}

func (p *initSubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false
	fs.StringVar(&p.domain, "domain", "my.domain", "domain for groups")
	fs.StringVar(&p.projectName, "project-name", "", "name of this project, the default being directory name")
	fs.StringVar(&p.version, "version", "", "resource version")
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

func (p *initSubcommand) PreScaffold(machinery.Filesystem) error {
	// Check if the current directory has not files or directories which does not allow to init the project
	return checkDir()
}

func (p *initSubcommand) Scaffold(fs machinery.Filesystem) error {
	scaffolder := scaffolds.NewInitScaffolder(p.config)
	scaffolder.InjectFS(fs)
	err := scaffolder.Scaffold()
	if err != nil {
		return err
	}

	return nil
}

func (p *initSubcommand) PostScaffold() error {
	// print follow on instructions to better guide the user
	fmt.Printf("Next: define a resource with:\n$ %s create api\n", p.commandName)
	return nil
}

// checkDir will return error if the current directory has files which are not allowed.
// Note that, it is expected that the directory to scaffold the project is cleaned.
// Otherwise, it might face issues to do the scaffold.
func checkDir() error {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Allow directory trees starting with '.'
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			// Allow files starting with '.'
			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}
			// Allow files ending with '.md' extension
			if strings.HasSuffix(info.Name(), ".md") && !info.IsDir() {
				return nil
			}
			// Allow capitalized files except PROJECT
			isCapitalized := true
			for _, l := range info.Name() {
				if !unicode.IsUpper(l) {
					isCapitalized = false
					break
				}
			}
			if isCapitalized && info.Name() != "PROJECT" {
				return nil
			}
			// Allow files in the following list
			allowedFiles := []string{
				"Cargo.toml",
				"Cargo.lock",
			}
			for _, allowedFile := range allowedFiles {
				if info.Name() == allowedFile {
					return nil
				}
			}
			// Do not allow any other file
			return fmt.Errorf(
				"target directory is not empty (only %s, files and directories with the prefix \".\", "+
					"files with the suffix \".md\" or capitalized files name are allowed); "+
					"found existing file %q", strings.Join(allowedFiles, ", "), path)
		})
	if err != nil {
		return err
	}
	return nil
}
