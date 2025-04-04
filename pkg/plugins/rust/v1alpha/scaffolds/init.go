/*
Copyright 2025 System Craftsman LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scaffolds

import (
	"errors"
	"fmt"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates/hack"
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/scaffolds/internal/templates/src"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugins"
)

var _ plugins.Scaffolder = &initScaffolder{}

type initScaffolder struct {
	config          config.Config
	boilerplatePath string
	license         string
	owner           string
	commandName     string

	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem
}

// NewInitScaffolder returns a new plugins.Scaffolder for project initialization operations
func NewInitScaffolder(config config.Config, license, owner, commandName string) plugins.Scaffolder {
	return &initScaffolder{
		config:          config,
		boilerplatePath: hack.DefaultBoilerplatePath,
		license:         license,
		owner:           owner,
		commandName:     commandName,
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

	if s.license != "none" {
		bpFile := &hack.Boilerplate{
			License: s.license,
			Owner:   s.owner,
		}
		bpFile.Path = s.boilerplatePath
		if err := scaffold.Execute(bpFile); err != nil {
			return err
		}

		boilerplate, err := afero.ReadFile(s.fs.FS, s.boilerplatePath)
		if err != nil {
			if errors.Is(err, afero.ErrFileNotFound) {
				log.Warnf("Unable to find %s: %s.\n"+"This file is used to generate the license header in the project.\n"+
					"Note that controller-gen will also use this. Therefore, ensure that you "+
					"add the license file or configure your project accordingly.",
					s.boilerplatePath, err)
				boilerplate = []byte("")
			} else {
				return fmt.Errorf("unable to load boilerplate: %w", err)
			}
		}
		// Initialize the machinery.Scaffold that will write the files to disk
		scaffold = machinery.NewScaffold(s.fs,
			machinery.WithConfig(s.config),
			machinery.WithBoilerplate(string(boilerplate)),
		)
	} else {
		s.boilerplatePath = ""
		// Initialize the machinery.Scaffold without boilerplate
		scaffold = machinery.NewScaffold(s.fs,
			machinery.WithConfig(s.config),
		)
	}

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
