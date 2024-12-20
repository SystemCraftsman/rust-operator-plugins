package rust

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	cfgv3 "sigs.k8s.io/kubebuilder/v4/pkg/config/v3"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

const pluginName = "rust.sdk.operatorframework.io"

var (
	pluginVersion            = plugin.Version{Number: 1}
	supportedProjectVersions = []config.Version{cfgv3.Version}
	pluginKey                = plugin.KeyFor(Plugin{})
)

var (
	_ plugin.Plugin    = Plugin{}
	_ plugin.Init      = Plugin{}
	_ plugin.CreateAPI = Plugin{}
)

type Plugin struct {
	initSubcommand
	createAPISubcommand
}

func (Plugin) Name() string                                         { return pluginName }
func (Plugin) Version() plugin.Version                              { return pluginVersion }
func (Plugin) SupportedProjectVersions() []config.Version           { return supportedProjectVersions }
func (p Plugin) GetInitSubcommand() plugin.InitSubcommand           { return &p.initSubcommand }
func (p Plugin) GetCreateAPISubcommand() plugin.CreateAPISubcommand { return &p.createAPISubcommand }
