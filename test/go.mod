module github.com/mesosphere/kubeaddons-base-addons/test

go 1.13

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/docker v1.4.2-0.20190916154449-92cc603036dd
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/mesosphere/ksphere-testing-framework v0.0.0-20200624200651-6b661edc6888
	github.com/mesosphere/kubeaddons v0.16.2
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.6
	k8s.io/helm v2.16.9+incompatible
	sigs.k8s.io/kind v0.7.0
)

replace k8s.io/client-go => k8s.io/client-go v0.18.4 // this is needed as long as kubeaddons uses an pre-1.18 version
