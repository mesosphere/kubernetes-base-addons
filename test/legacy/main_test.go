package main_test

import (
	"testing"

	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
)

var testAddon = "cert-manager"

func TestSetup(t *testing.T) {
	r, err := local.NewRepository("kba", "../../")
	if err != nil {
		t.Fatalf("unable to read directory as an addon repository: %s", err)
	}

	a, err := r.GetAddon(testAddon)
	if err != nil {
		t.Fatalf("unable to retrieve %s from the addon repository: %s", testAddon, err)
	}

	if a[0].GetName() != testAddon {
		t.Fatalf("expected addon name %s, received %s", testAddon, a[0].GetName())
	}

	t.Logf("SUCCESS: able to retrieve addon %s: addon repository was comprehensible with older kubeaddons catalog, all set!", testAddon)
}
