package version

import "fmt"

const (
	Unknown    = "unknown"
	modulePath = "github.com/SystemCraftsman/rust-operator-plugins"
)

var (
	GitVersion = Unknown
	GitCommit  = Unknown
)

var Version = Context{
	Name:    modulePath,
	Version: GitVersion,
	Commit:  GitCommit,
}

type Context struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func (vc *Context) String() string {
	return fmt.Sprintf("%s <commit: %s>", vc.Version, vc.Commit)
}
