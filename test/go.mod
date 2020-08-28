module github.com/mesosphere/kubeaddons-base-addons/test

go 1.14

require (
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/google/uuid v1.1.1
	github.com/mesosphere/ksphere-testing-framework v0.0.0-20200814171113-1a98809a8734
	github.com/mesosphere/kubeaddons v0.19.0
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.16.10+incompatible
	sigs.k8s.io/kind v0.8.1
)

replace (
	github.com/mesosphere/ksphere-testing-framework => github.com/mesosphere/ksphere-testing-framework v0.0.0-20200806132303-8e10596082d3
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/client-go => k8s.io/client-go v0.18.6
)
