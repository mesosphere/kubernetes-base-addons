module github.com/mesosphere/kubeaddons-base-addons/test

go 1.13

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/emicklei/go-restful v2.11.1+incompatible // indirect
	github.com/go-openapi/spec v0.19.4 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/mesosphere/kubeaddons v0.6.1-0.20191205191456-60012764fb31
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/net v0.0.0-20191119073136-fc4aabc6c914 // indirect
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/yaml.v2 v2.2.7
	k8s.io/apiextensions-apiserver v0.0.0-20191121021419-88daf26ec3b8 // indirect
	k8s.io/utils v0.0.0-20191114200735-6ca3b61696b6 // indirect
	sigs.k8s.io/controller-runtime v0.4.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
