package testutils

import (
	"bytes"
	"os"
)

func CaptureOutput(f func()) string {
	// Backup original stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the function
	f()

	// Restore stdout and capture the output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}
