module github.com/mesosphere/kubeaddons-base-addons/test

go 1.13

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/docker v1.4.2-0.20190916154449-92cc603036dd
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/mesosphere/ksphere-testing-framework v0.0.0-20200530001136-9d1c380ca073
	github.com/mesosphere/kubeaddons v0.15.1
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/helm v2.16.8+incompatible
	sigs.k8s.io/kind v0.7.0
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.4
	k8s.io/apiserver => k8s.io/apiserver v0.17.4
	k8s.io/client-go => k8s.io/client-go v0.17.4
	k8s.io/kubectl => k8s.io/kubectl v0.17.4
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.2
)
