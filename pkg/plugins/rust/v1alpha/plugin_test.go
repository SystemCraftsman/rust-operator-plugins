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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/stage"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

var _ = Describe("Plugin Test", func() {
	testPlugin := &Plugin{}

	Describe("Name", func() {
		It("should return the correct plugin name", func() {
			Expect(testPlugin.Name()).To(Equal("rust.sdk.operatorframework.io"))
		})
	})

	Describe("Version", func() {
		It("should return the correct plugin version", func() {
			Expect(testPlugin.Version()).To(Equal(plugin.Version{Number: 1, Stage: stage.Alpha}))
		})
	})

	Describe("GetInitSubcommand", func() {
		It("should return the correct plugin initSubcommand", func() {
			Expect(testPlugin.GetInitSubcommand()).To(Equal(&testPlugin.initSubcommand))
		})
	})

	Describe("GetCreateAPISubcommand", func() {
		It("should return the correct plugin createAPISubcommand", func() {
			Expect(testPlugin.GetCreateAPISubcommand()).To(Equal(&testPlugin.createAPISubcommand))
		})
	})
})
