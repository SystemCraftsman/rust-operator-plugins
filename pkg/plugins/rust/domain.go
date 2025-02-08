package rust

import plugins "github.com/SystemCraftsman/rust-operator-plugins/pkg"

// DefaultNameQualifier is the suffix appended to all kubebuilder plugin names for Rust operators.
const DefaultNameQualifier = "rust." + plugins.DefaultNameQualifier
