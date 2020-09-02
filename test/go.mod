module github.com/mesosphere/kubeaddons-base-addons/test

go 1.15

replace k8s.io/client-go => k8s.io/client-go v0.19.0

require (
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/evanphx/json-patch/v5 v5.1.0 // indirect
	github.com/go-git/go-git/v5 v5.1.0 // indirect
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6 // indirect
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mesosphere/ksphere-testing-framework v0.0.0-20200824140305-1e889a7c7868
	github.com/mesosphere/kubeaddons v0.19.0
	github.com/nightlyone/lockfile v1.0.0 // indirect
	go.uber.org/goleak v1.1.10 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200828194041-157a740278f4 // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	helm.sh/helm/v3 v3.3.0 // indirect
	k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver v0.19.0 // indirect
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.16.10+incompatible
	k8s.io/klog/v2 v2.3.0 // indirect
	k8s.io/kubectl v0.18.6 // indirect
	k8s.io/utils v0.0.0-20200821003339-5e75c0163111 // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/kind v0.8.1
	sigs.k8s.io/structured-merge-diff/v3 v3.0.1-0.20200706213357-43c19bbb7fba // indirect
)
