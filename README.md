# Kubernetes Base Addons

This repository contains the default sets of addons needed to bootstrap [D2iQ Kubernetes](https://d2iq.com/solutions/ksphere).

# Overview

The structure of this repository follows the [Kubeaddons Catalog Documentation](https://github.com/mesosphere/kubeaddons/blob/master/tools/catalog/README.md) in reference, and uses the [Addon Revision concept covered therein](https://github.com/mesosphere/kubeaddons/blob/master/tools/catalog/README.md#special-addonrepository-options---addon-revisions).

You will find the following directories here:

* `addons/` - containing the actual manifests for addon resources
* `metadata/` - containing static metadata for the addons in `addons/`
* `test/` - containing integration tests for the addons in `addons/`

# Contributing

See the [Kubeaddons Contributing Documentation](https://github.com/mesosphere/kubeaddons/blob/master/CONTRIBUTING.md).

## Addon Revisions

You will find that any particular addon directory (e.g. `addons/prometheus`) may have several directories and several manifests nested in them each with variants of that addon. These are what we refered to above as "revisions".

The **intention of revisions is to maintain a flat history of addon changes**. If you are making changes to any particular addon you should be making a revision of that addon as a copy of the original file with the changes made therein and the `addon-revision` version updated to reflect the new version appropriately.

### notes & caveats

In the `v1beta1` version our of API, there are some elements of our addon specifications that act as metadata and [may be removed in future](https://jira.mesosphere.com/browse/DCOS-62438). Here are some conditions under which *a new revision will not be needed* when updating an addon:

1. If you're *only adding* a new `CloudProvider`
2. If you're *only adding* labels or annotations to the object metdata

These cases are safe to update in place.
