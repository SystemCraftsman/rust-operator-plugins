package rust

import (
	"github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"os"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin/util"
)

var _ = Describe("API test", func() {
	var (
		testAPISubcommand createAPISubcommand
	)

	BeforeEach(func() {
		testAPISubcommand = createAPISubcommand{
			resourceFlag:   &pflag.Flag{Changed: true},
			controllerFlag: &pflag.Flag{Changed: true},
			options: &rust.Options{
				DoAPI:        true,
				DoController: true,
				Namespaced:   true,
			},
		}
	})

	Describe("UpdateResource", func() {
		It("verify that resource fields were set", func() {
			testAPIOptions := &rust.Options{
				Namespaced:   true,
				DoAPI:        true,
				DoController: true,
			}
			updateTestResource := resource.Resource{}
			testAPIOptions.UpdateResource(&updateTestResource)
			Expect(updateTestResource.API.Namespaced).To(Equal(testAPIOptions.Namespaced))
			Expect(updateTestResource.API.CRDVersion).To(Equal("v1"))
			Expect(updateTestResource.Controller).To(Equal(testAPIOptions.DoController))
			Expect(updateTestResource.Path).To(Equal(""))
		})
	})

	Describe("BindFlags", func() {
		It("should set SortFlags to false", func() {
			flagTest := pflag.NewFlagSet("testFlag", -1)
			testAPISubcommand.BindFlags(flagTest)
			Expect(flagTest.SortFlags).To(BeTrue())
			Expect(testAPISubcommand.options.DoController).To(BeTrue())
			Expect(testAPISubcommand.options.DoAPI).To(BeTrue())
			Expect(testAPISubcommand.options.Namespaced).To(BeTrue())
		})
	})

	Describe("InjectConfig", func() {
		It("should set config", func() {
			testConfig, _ := config.New(config.Version{Number: 3})
			err := testAPISubcommand.InjectConfig(testConfig)
			Expect(testAPISubcommand.config).To(Equal(testConfig))
			Expect(err).To(BeNil())
		})
	})
	
	Describe("PostScaffold", func() {
		It("should return nil", func() {
			// Create a temporary directory for testing
			tmpDir, _ := os.MkdirTemp("", "test-dir-")
			defer os.RemoveAll(tmpDir) // Clean up after the test

			// Change the working directory to the temporary one
			_ = os.Chdir(tmpDir)

			err := util.RunCmd("Format code", "cargo", "init")
			Expect(err).To(BeNil())

			testResource := resource.Resource{
				GVK: resource.GVK{
					Group:   "test-group",
					Version: "v1",
					Kind:    "Test-Kind",
				},
				Plural: "test-plural",
			}

			testConfig, _ := config.New(config.Version{Number: 3})
			testAPISubcommand.InjectConfig(testConfig)
			err = testAPISubcommand.InjectResource(&testResource)
			Expect(err).To(BeNil())

			Expect(testAPISubcommand.PostScaffold()).To(BeNil())
		})
	})

	Describe("InjectResource", func() {
		It("verify that wrong GVKs fail", func() {
			failResource := resource.Resource{
				GVK: resource.GVK{
					Group:   "Fail-Test-Group",
					Version: "test-version",
					Kind:    "test-kind",
				},
				Plural: "test-plural",
			}

			testConfig, _ := config.New(config.Version{Number: 3})
			testAPISubcommand.InjectConfig(testConfig)
			groupErr := testAPISubcommand.InjectResource(&failResource)
			Expect(testAPISubcommand.resource, failResource)
			Expect(groupErr).To(HaveOccurred())

			failResource.GVK.Group = "test-group"
			versionErr := testAPISubcommand.InjectResource(&failResource)
			Expect(versionErr).To(HaveOccurred())

			failResource.GVK.Version = "v1"
			kindError := testAPISubcommand.InjectResource(&failResource)
			Expect(kindError).To(HaveOccurred())
		})

		It("verify that a correct GVK succeeds", func() {
			testResource := resource.Resource{
				GVK: resource.GVK{
					Group:   "test-group",
					Version: "v1",
					Kind:    "Test-Kind",
				},
				Plural: "test-plural",
			}

			testConfig, _ := config.New(config.Version{Number: 3})
			testAPISubcommand.InjectConfig(testConfig)
			noErr := testAPISubcommand.InjectResource(&testResource)
			Expect(testAPISubcommand.resource, testResource)
			Expect(noErr).To(BeNil())
		})
	})
})
