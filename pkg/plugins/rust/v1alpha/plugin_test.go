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
