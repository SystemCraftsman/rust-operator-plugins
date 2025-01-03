/*
This file is add as a temporary replacement for kubebuilder's marker
as it doesn't support .rs file extensions currently.

When contributed to kubebuilder, this file won't be necessary.
*/

package machinery

import (
	"fmt"
	"path/filepath"
	"strings"
)

const kbPrefix = "+kubebuilder:scaffold:"

var commentsByExt = map[string]string{
	".go":   "//",
	".yaml": "#",
	".yml":  "#",
	".rs":   "//",
	// When adding additional file extensions, update also the NewMarkerFor documentation and error
}

type CustomMarker struct {
	prefix  string
	comment string
	value   string
}

// NewMarkerFor creates a new marker customized for the specific file. The created marker
// is prefixed with `+kubebuilder:scaffold:` the default prefix for kubebuilder.
// Supported file extensions: .go, .yaml, .yml, .rs.
func NewMarkerFor(path string, value string) CustomMarker {
	return NewMarkerWithPrefixFor(kbPrefix, path, value)
}

// NewMarkerWithPrefixFor creates a new custom prefixed marker customized for the specific file
// Supported file extensions: .go, .yaml, .yml
func NewMarkerWithPrefixFor(prefix string, path string, value string) CustomMarker {

	ext := filepath.Ext(path)
	if comment, found := commentsByExt[ext]; found {
		return CustomMarker{
			prefix:  markerPrefix(prefix),
			comment: comment,
			value:   value,
		}
	}

	extensions := make([]string, 0, len(commentsByExt))
	for extension := range commentsByExt {
		extensions = append(extensions, fmt.Sprintf("%q", extension))
	}
	panic(fmt.Errorf("unknown file extension: '%s', expected one of: %s", ext, strings.Join(extensions, ", ")))
}

// String implements Stringer
func (m CustomMarker) String() string {
	return m.comment + " " + m.prefix + m.value
}

// EqualsLine compares a marker with a string representation to check if they are the same marker
func (m CustomMarker) EqualsLine(line string) bool {
	line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), m.comment))
	return line == m.prefix+m.value
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
