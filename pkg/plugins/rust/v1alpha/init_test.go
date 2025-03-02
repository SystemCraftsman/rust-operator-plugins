package rust

import (
	"fmt"
	testutils "github.com/SystemCraftsman/rust-operator-plugins/pkg/plugins/rust/v1alpha/utils/test"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

func TestInitSubcommand_UpdateMetadata(t *testing.T) {
	successInitSubcommand := initSubcommand{
		domain: "testDomain",
	}

	cliMetadata := plugin.CLIMetadata{CommandName: "TestCommand"}
	subcommandMetadata := plugin.SubcommandMetadata{}

	// Test that commandName is updated correctly
	successInitSubcommand.UpdateMetadata(cliMetadata, &subcommandMetadata)
	assert.Equal(t, "TestCommand", successInitSubcommand.commandName)
}

func TestInitSubcommand_BindFlags(t *testing.T) {
	successInitSubcommand := initSubcommand{}
	flagTest := pflag.NewFlagSet("testFlag", pflag.ContinueOnError)

	successInitSubcommand.BindFlags(flagTest)

	// Ensure the flags were set correctly
	assert.False(t, flagTest.SortFlags)
	assert.Equal(t, "my.domain", successInitSubcommand.domain)
	assert.Equal(t, "", successInitSubcommand.projectName)
}

func TestInitSubcommand_InjectConfig(t *testing.T) {
	successInitSubcommand := initSubcommand{}
	testConfig, _ := config.New(config.Version{Number: 3})

	// Inject valid config and check that values are set
	err := successInitSubcommand.InjectConfig(testConfig)
	assert.NoError(t, err)

	dir, _ := os.Getwd()
	assert.Equal(t, testConfig.GetDomain(), successInitSubcommand.domain)
	assert.Equal(t, strings.ToLower(filepath.Base(dir)), successInitSubcommand.projectName)
	assert.Equal(t, testConfig.GetProjectName(), successInitSubcommand.projectName)

	// Inject invalid config to test failure
	failureInitSubcommand := initSubcommand{projectName: "?&fail&?"}
	err = failureInitSubcommand.InjectConfig(testConfig)
	assert.Error(t, err)
}

func TestInitSubcommand_PreScaffold(t *testing.T) {
	successInitSubcommand := initSubcommand{}
	err := successInitSubcommand.PreScaffold(machinery.Filesystem{})
	assert.Error(t, err)

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test-dir-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) // Clean up after the test

	// Change the working directory to the temporary one
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = successInitSubcommand.PreScaffold(machinery.Filesystem{})
	assert.NoError(t, err)
}

func TestInitSubcommand_PostScaffold(t *testing.T) {
	successInitSubcommand := initSubcommand{
		commandName: "init",
	}

	capturedOutput := testutils.CaptureOutput(func() {
		successInitSubcommand.PostScaffold()
	})

	expectedOutput := fmt.Sprintf("Next: define a resource with:\n$ %s create api\n",
		successInitSubcommand.commandName)

	assert.Equal(t, expectedOutput, capturedOutput)
}
