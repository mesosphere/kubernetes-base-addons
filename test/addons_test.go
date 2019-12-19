package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/blang/semver"
	"gopkg.in/yaml.v2"

	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

const defaultKubernetesVersion = "1.15.6"

var addonTestingGroups = make(map[string][]string)

func init() {
	b, err := ioutil.ReadFile("groups.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, addonTestingGroups); err != nil {
		panic(err)
	}
}

func TestValidateUnhandledAddons(t *testing.T) {
	unhandled, err := findUnhandled()
	if err != nil {
		t.Fatal(err)
	}

	if len(unhandled) != 0 {
		names := make([]string, len(unhandled))
		for _, addon := range unhandled {
			names = append(names, addon.GetName())
		}
		t.Fatal(fmt.Errorf("the following addons are not handled as part of a testing group: %+v", names))
	}
}

func TestGeneralGroup(t *testing.T) {
	if err := testgroup(t, "general"); err != nil {
		t.Fatal(err)
	}
}

//func TestElasticSearchGroup(t *testing.T) {
//	if err := testgroup(t, "elasticsearch"); err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestPrometheusGroup(t *testing.T) {
//	if err := testgroup(t, "prometheus"); err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestKommanderGroup(t *testing.T) {
//	if err := testgroup(t, "kommander"); err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestIstioGroup(t *testing.T) {
//	if err := testgroup(t, "istio"); err != nil {
//		t.Fatal(err)
//	}
//}


// -----------------------------------------------------------------------------
// Private Functions
// -----------------------------------------------------------------------------

func testgroup(t *testing.T, groupname string) error {
	t.Logf("testing group %s", groupname)

	version, err := semver.Parse(defaultKubernetesVersion)
	if err != nil {
		return err
	}

	cluster, err := kind.NewCluster(version)
	if err != nil {
		return err
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster); err != nil {
		return err
	}

	if err := deployCertManagerCA(cluster); err != nil {
		return err
	}

	addons, err := addons(addonTestingGroups[groupname]...)
	if err != nil {
		return err
	}

	ph, err := test.NewBasicTestHarness(t, cluster, addons...)
	if err != nil {
		return err
	}
	defer ph.Cleanup()

	ph.Validate()
	ph.Deploy()

	return nil
}

func addons(names ...string) ([]v1beta1.AddonInterface, error) {
	var testAddons []v1beta1.AddonInterface

	addons, err := temp.Addons("../addons/")
	if err != nil {
		return testAddons, err
	}

	for _, addon := range addons {
		for _, name := range names {
			overrides(addon[0])
			if addon[0].GetName() == name {
				testAddons = append(testAddons, addon[0])
			}
		}
	}

	if len(testAddons) != len(names) {
		return testAddons, fmt.Errorf("got %d addons, expected %d", len(testAddons), len(names))
	}

	return testAddons, nil
}

func findUnhandled() ([]v1beta1.AddonInterface, error) {
	var unhandled []v1beta1.AddonInterface

	addons, err := temp.Addons("../addons/")
	if err != nil {
		return unhandled, err
	}

	for _, revisions := range addons {
		addon := revisions[0]
		found := false
		for _, v := range addonTestingGroups {
			for _, name := range v {
				if name == addon.GetName() {
					found = true
				}
			}
		}
		if !found {
			unhandled = append(unhandled, addon)
		}
	}

	return unhandled, nil
}

func deployCertManagerCA(cluster test.Cluster) error {
	if err := kubectl("--kubeconfig", cluster.ConfigPath(), "create", "namespace", "cert-manager"); err != nil {
		return err
	}

	// create secret
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization: []string{"d2iq"},
			Country:      []string{"US"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)

	certPath := path.Join(wd, "ca.crt")
	keyPath := path.Join(wd, "ca.key")

	certOut, err := os.Create(certPath)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
	certOut.Close()

	// Private key
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	// create kubernetes-root-ca secret
	if err := kubectl("--kubeconfig", cluster.ConfigPath(), "create", "secret", "tls", "kubernetes-root-ca", "--namespace=cert-manager", fmt.Sprintf("--cert=%s", certPath), fmt.Sprintf("--key=%s", keyPath)); err != nil {
		return err
	}

	// create konvoyconfig-kubeaddons configmap
	//if err := kubectl("--kubeconfig", cluster.ConfigPath(), "create", "configmap", "konvoyconfig-kubeaddons", "--namespace=kubeaddons", "--from-literal=clusterHostname=kommander.example.com"); err != nil {
	//	return err
	//}

	defer os.Remove(certPath)
	defer os.Remove(keyPath)
	return nil
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// -----------------------------------------------------------------------------
// Private - CI Values Overrides
// -----------------------------------------------------------------------------

// TODO: a temporary place to put configuration overrides for addons
// See: https://jira.mesosphere.com/browse/DCOS-62137
func overrides(addon v1beta1.AddonInterface) {
	if v, ok := addonOverrides[addon.GetName()]; ok {
		addon.GetAddonSpec().ChartReference.Values = &v
	}
}

var addonOverrides = map[string]string{
	"metallb": `
---
configInline:
  address-pools:
  - name: default
    protocol: layer2
    addresses:
    - "172.17.1.200-172.17.1.250"
`,
	"istio":
`---
      kiali:
       enabled: true
       contextPath: /ops/portal/kiali
       ingress:
         enabled: true
         kubernetes.io/ingress.class: traefik
         hosts:
           - ""
       dashboard:
         auth:
           strategy: anonymous
       prometheusAddr: http://prometheus-kubeaddons-prom-prometheus.kubeaddons:9090

      tracing:
        enabled: true
        contextPath: /ops/portal/jaeger
        ingress:
          enabled: true
          kubernetes.io/ingress.class: traefik
          hosts:
            - ""

      grafana:
        enabled: true

      prometheus:
        serviceName: prometheus-kubeaddons-prom-prometheus.kubeaddons

      istiocoredns:
        enabled: true

      security:
       selfSigned: true
       caCert: /etc/cacerts/tls.crt
       caKey: /etc/cacerts/tls.key
       rootCert: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
       certChain: /etc/cacerts/tls.crt
       enableNamespacesByDefault: false

      global:
       podDNSSearchNamespaces:
       - global
       - "{{ valueOrDefault .DeploymentMeta.Namespace \"default\" }}.global"

       mtls:
        enabled: true

       multiCluster:
        enabled: true

       controlPlaneSecurityEnabled: true
`,
}

