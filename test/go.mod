module github.com/mesosphere/kubeaddons-base-addons/test

go 1.15

replace (
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce // locked to this version to avoid upgrades from kind that would drop our volume provisioning code. TODO: we should update/remove our volume provisioning code.
	k8s.io/client-go => k8s.io/client-go v0.19.2 // locked for general sanity in k8s lib resolution
	sigs.k8s.io/kind => sigs.k8s.io/kind v0.9.0 // locked to avoid changes, as we've historically had issues with kind being changed as a side effect of other changes
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/datawire/ambassador v1.7.3
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/go-logr/logr v0.2.1 // indirect
	github.com/google/go-cmp v0.5.2
	github.com/google/uuid v1.1.2
	github.com/mesosphere/ksphere-testing-framework v0.2.3
	github.com/mesosphere/kubeaddons v0.22.3
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.17.0+incompatible
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	sigs.k8s.io/controller-runtime v0.6.3
	sigs.k8s.io/kind v0.9.0
)
