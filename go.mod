module github.com/SystemCraftsman/rust-operator-plugins

go 1.23.4

require (
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.36.2
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/afero v1.12.0
	github.com/spf13/pflag v1.0.6
	k8s.io/apimachinery v0.32.2
	sigs.k8s.io/kubebuilder/v4 v4.2.0
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gobuffalo/flect v1.0.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/utils v0.0.0-20241210054802-24370beab758 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace sigs.k8s.io/kubebuilder/v4 => github.com/mabulgu/kubebuilder/v4 v4.2.1-rust
