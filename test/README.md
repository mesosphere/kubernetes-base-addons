# Addon Tests

In this directory you will find [go tests](https://golang.org/pkg/testing/) which cover integration testing of the addons in this repo.

This uses the [Kubeaddons](https://github.com/mesosphere/kubeaddons) addon testing framework, which is [documented here](https://github.com/mesosphere/kubeaddons/blob/master/docs/test/framework.md).

## New Addon Tests

When addons are added to the repository, CI will fail on validation if tests (that  pass) are not provided for them.

In order to add tests for a new addon you'll need to add it to an existing (or new) testing group in the [groups.yaml](/test/groups.yaml) configuration file.

If you create a new testing group, you must update the test functions in [addons_test.go](/test/addons_test.go) to cover the new testing group.

At it's simplest, a test of a testing group may look like the following:

```golang
func TestGeneralGroup(t *testing.T) {
	if err := testgroup(t, "general"); err != nil {
		t.Fatal(err)
	}
}
```

Where `"general"` is a testing group containing several addons to test.

From here, you can expand your tests within this function.

