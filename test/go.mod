module github.com/mesosphere/kubeaddons-base-addons/test

go 1.13

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/google/uuid v1.1.1
	github.com/mesosphere/ksphere-testing-framework v0.0.0-20200624200651-6b661edc6888
	github.com/mesosphere/kubeaddons v0.18.2
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/helm v2.16.9+incompatible
	sigs.k8s.io/kind v0.8.1
)

replace k8s.io/client-go => k8s.io/client-go v0.18.4 // this is needed as long as kubeaddons uses an pre-1.18 version
