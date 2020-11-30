module github.com/mesosphere/kubeaddons-base-addons/newtest

go 1.15

replace (
	github.com/mesosphere/dkp-test-framework => ../../dkp-test-framework
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	k8s.io/client-go => k8s.io/client-go v0.19.4
)

require (
	github.com/mesosphere/dkp-test-framework v0.0.0-00010101000000-000000000000
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)
