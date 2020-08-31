module github.com/mesosphere/kubeaddons-base-addons/test

go 1.15

replace k8s.io/client-go => k8s.io/client-go v0.19.0

require (
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/google/uuid v1.1.1
	github.com/mesosphere/ksphere-testing-framework v0.0.0-20200824140305-1e889a7c7868
	github.com/mesosphere/kubeaddons v0.19.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200828194041-157a740278f4 // indirect
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.16.10+incompatible
	k8s.io/utils v0.0.0-20200821003339-5e75c0163111 // indirect
	sigs.k8s.io/kind v0.8.1
)
