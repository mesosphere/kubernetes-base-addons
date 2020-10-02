module github.com/mesosphere/kubeaddons-base-addons/test

go 1.15

replace (
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce // locked to this version to avoid upgrades from kind that would drop our volume provisioning code. TODO: we should update/remove our volume provisioning code.
	k8s.io/client-go => k8s.io/client-go v0.19.1 // locked for general sanity in k8s lib resolution
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.1-0.20200909023352-d6829e9c4db8 // pinned to SHA to allow 1.19.x libs at the time of writing, as controller-runtime was behind
	sigs.k8s.io/kind => sigs.k8s.io/kind v0.9.0 // locked to avoid changes, as we've historically had issues with kind being changed as a side effect of other changes
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/google/uuid v1.1.2
	github.com/mesosphere/ksphere-testing-framework v0.2.0
	github.com/mesosphere/kubeaddons v0.22.2
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.16.12+incompatible
	sigs.k8s.io/kind v0.9.0
)
