package rust

import (
	"fmt"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"strings"
)

const kbPrefix = "+kubebuilder:scaffold:"

var commentsByExt = map[string]string{
	".rs": "//",
}

func NewMarkerFor(path string, value string) machinery.Marker {
	ext := filepath.Ext(path)
	if comment, ok := commentsByExt[ext]; ok {
		return machinery.Marker{
			Prefix:  markerPrefix(kbPrefix),
			Comment: comment,
			Value:   value,
		}

	}
	panic(fmt.Errorf("unsupported file extension: %s", ext))
}

func markerPrefix(prefix string) string {
	trimmed := strings.TrimSpace(prefix)
	var builder strings.Builder
	if !strings.HasPrefix(trimmed, "+") {
		builder.WriteString("+")
	}
	builder.WriteString(trimmed)
	if !strings.HasSuffix(trimmed, ":") {
		builder.WriteString(":")
	}
	return builder.String()
}
