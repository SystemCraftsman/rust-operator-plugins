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

package rust

import (
	"os"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"

	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

var _ = Describe("Init test", func() {
	var (
		successInitSubcommand initSubcommand
	)

	BeforeEach(func() {
		successInitSubcommand = initSubcommand{
			domain:      "testDomain",
			commandName: "testCommand",
		}
	})

	Describe("UpdateMetadata", func() {
		It("Check that function call sets data correctly", func() {
			testCliMetadata := plugin.CLIMetadata{CommandName: "TestCommand"}
			testSubcommandMetadata := plugin.SubcommandMetadata{}
			Expect(successInitSubcommand.commandName).NotTo(Equal(testCliMetadata.CommandName))

			successInitSubcommand.UpdateMetadata(testCliMetadata, &testSubcommandMetadata)
			Expect(successInitSubcommand.commandName).To(Equal(testCliMetadata.CommandName))
		})
	})

	Describe("BindFlags", func() {
		It("verify all fields were set correctly", func() {
			flagTest := pflag.NewFlagSet("testFlag", -1)
			successInitSubcommand.BindFlags(flagTest)
			Expect(flagTest.SortFlags).To(BeFalse())
			Expect(successInitSubcommand.domain).To(Equal("my.domain"))
			Expect(successInitSubcommand.projectName).To(Equal(""))
			Expect(successInitSubcommand.version).To(Equal(""))
		})
	})

	Describe("InjectConfig", func() {
		It("verify all fields were set correctly", func() {
			testConfig, _ := config.New(config.Version{Number: 3})
			dir, _ := os.Getwd()
			Expect(successInitSubcommand.InjectConfig(testConfig)).To(Succeed())
			Expect(successInitSubcommand.config).To(Equal(testConfig))
			Expect(successInitSubcommand.domain).To(Equal(testConfig.GetDomain()))
			Expect(successInitSubcommand.projectName).To(Equal(strings.ToLower(filepath.Base(dir))))
			Expect(successInitSubcommand.projectName).To(Equal(testConfig.GetProjectName()))
			Expect(successInitSubcommand.InjectConfig(testConfig)).To(BeNil())
		})
	})

	Describe("PreScaffold", func() {
		It("should return nil", func() {
			// Create a temporary directory for testing
			tmpDir, _ := os.MkdirTemp("", "test-dir-")
			defer os.RemoveAll(tmpDir) // Clean up after the test

			// Change the working directory to the temporary one
			_ = os.Chdir(tmpDir)

			Expect(successInitSubcommand.PreScaffold(machinery.Filesystem{})).To(BeNil())
		})
	})

	Describe("PostScaffold", func() {
		It("should return nil", func() {
			Expect(successInitSubcommand.PostScaffold()).To(BeNil())
		})
	})
})
