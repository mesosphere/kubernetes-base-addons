# Kubernetes Base Addons

![build status label](https://teamcity.mesosphere.io/app/rest/builds/buildType:(id:kubeaddons_KubernetesBaseAddons),branch:(default:true)/statusIcon)

This repository contains the default sets of addons needed to bootstrap [D2iQ Kubernetes](https://d2iq.com/solutions/ksphere).

## Overview

The structure of this repository follows the [Kubeaddons Catalog Documentation](https://github.com/mesosphere/kommander-catalog-api/blob/master/README.md) in reference, and uses the [Addon Revision concept covered therein](https://github.com/mesosphere/kommander-catalog-api/blob/master/README.md#special-addonrepository-options---addon-revisions).

You will find the following directories here:

* `addons/` - containing the actual manifests for addon resources
* `metadata/` - containing static metadata for the addons in `addons/`
* `test/` - containing integration tests for the addons in `addons/`

## Support

There is **no community support for this software at this time**, and this repository (while publicly visible) is only intended for consumption by D2iQ personnel.

If you are here looking for addon support as an end-user, support is given as per the support provided for your [D2iQ Ksphere Solution](https://d2iq.com/solutions/ksphere)'s provided [enterprise support](https://d2iq.com/services-and-support) (solutions such as [Konvoy](https://d2iq.com/solutions/ksphere/konvoy)), but **support is not given *directly* via this repository**.

Please [contact D2iQ](https://d2iq.com/contact) for questions and more information.

## Releases

Releases signify a single supported instance of `kubernetes-base-addons`. They undergo significant integration and soak testing and must be used together to be supported. Any variation from these sets will not be supported.

While you may see other tags in our [releases page](https://github.com/mesosphere/kubernetes-base-addons/releases) the only releases which are official and supported releases are named with the prefix `stable` and not marked as prerelease.

**NOTE**: Do not use `master` for production. Instead, pick a supported release version.

For the release process, see the [release](RELEASE.md) document.

### version 3

The `master` branch of `kubernetes-base-addons` has made a hard transition to cert-manager v1. Prior to this change, this repository was using the v1alpha1 resource definition which cannot be upgraded to v1alpha2 or beyond.
Releases from this branch will be tagged with the major version 3.

### version 2

`release/2` is going to have backports of changes made to master in order to continue to support cert-manager v1alpha1.
The releases cut from this branch will be major version 2.

## Testing

The test suite can be exercised locally by running

    make test

Pull Requests against this repo is tested by [Teamcity](https://teamcity.mesosphere.io/viewType.html?buildTypeId=kubeaddons_KubernetesBaseAddons) and [Dispatch](https://konvoy-staging.production.d2iq.cloud/dispatch/tekton/#/pipelineruns).

**NOTE**: E2E tests for the UI are _only_ run in Dispatch.

[Dispatchfile](Dispatchfile) defines the config and exercises the test in the Makefile.

## Contributing

See our [Contributing Documentation](CONTRIBUTING.md).
