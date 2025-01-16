package templates

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

var _ machinery.Template = &DockerIgnore{}

// DockerIgnore scaffolds a file that defines which files should be ignored by the containerized build process
type DockerIgnore struct {
	machinery.TemplateMixin
}

// SetTemplateDefaults implements file.Template
func (f *DockerIgnore) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = ".dockerignore"
	}

	f.TemplateBody = dockerignorefileTemplate

	return nil
}

const dockerignorefileTemplate = `# Include any files or directories that you don't want to be copied to your
# container here (e.g., local build artifacts, temporary files, etc.).
#
# For more help, visit the .dockerignore file reference guide at
# https://docs.docker.com/engine/reference/builder/#dockerignore-file

**/.DS_Store
**/.classpath
**/.dockerignore
**/.env
**/.git
**/.gitignore
**/.project
**/.settings
**/.toolstarget
**/.vs
**/.vscode
**/*.*proj.user
**/*.dbmdl
**/*.jfm
**/charts
**/docker-compose*
**/compose*
**/Dockerfile*
**/node_modules
**/npm-debug.log
**/secrets.dev.yaml
**/values.dev.yaml
/bin
/target
LICENSE
README.md
`
