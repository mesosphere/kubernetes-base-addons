module github.com/mesosphere/kubeaddons-base-addons/test

go 1.15

replace (
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce // locked to this version to avoid upgrades from kind that would drop our volume provisioning code. TODO: we should update/remove our volume provisioning code.
	k8s.io/client-go => k8s.io/client-go v0.19.2 // locked for general sanity in k8s lib resolution
	sigs.k8s.io/kind => sigs.k8s.io/kind v0.9.0 // locked to avoid changes, as we've historically had issues with kind being changed as a side effect of other changes
)

require (
	github.com/alecthomas/gometalinter v3.0.0+incompatible // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/coreos/etcd v3.3.15+incompatible // indirect
	github.com/datawire/ambassador v1.11.0
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/emicklei/go-restful v2.9.6+incompatible // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.2.0
	github.com/gophercloud/gophercloud v0.2.0 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc // indirect
	github.com/kisielk/errcheck v1.4.0 // indirect
	github.com/mesosphere/ksphere-testing-framework v0.2.6
	github.com/mesosphere/kubeaddons v0.24.1
	github.com/nicksnyder/go-i18n v1.10.1 // indirect
	github.com/nightlyone/lockfile v0.0.0-20180618180623-0ad87eef1443 // indirect
	github.com/sclevine/agouti v3.0.0+incompatible // indirect
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9 // indirect
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20171010053543-63abe20a23e2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.19.5
	k8s.io/apimachinery v0.19.5
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.17.0+incompatible
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009
	sigs.k8s.io/controller-runtime v0.6.4
	sigs.k8s.io/kind v0.10.0
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
	sigs.k8s.io/testing_frameworks v0.1.2 // indirect
)
