# Kubernetes Base Addons

This repository contains the default sets of addons needed to bootstrap [D2iQ Kubernetes](https://d2iq.com/solutions/ksphere).

## Overview

The structure of this repository follows the [Kubeaddons Catalog Documentation](https://github.com/mesosphere/kubeaddons/blob/master/tools/catalog/README.md) in reference, and uses the [Addon Revision concept covered therein](https://github.com/mesosphere/kubeaddons/blob/master/tools/catalog/README.md#special-addonrepository-options---addon-revisions).

You will find the following directories here:

* `addons/` - containing the actual manifests for addon resources
* `deployments/` - containing the default addons depending on the Kubernetes version
* `metadata/` - containing static metadata for the addons in `addons/`
* `test/` - containing integration tests for the addons in `addons/`

## Support

There is **no community support for this software at this time**, and this repository (while publicly visible) is only intended for consumption by D2iQ personnel.

If you are here looking for addon support as an end-user, support is given as per the support provided for your [D2iQ Ksphere Solution](https://d2iq.com/solutions/ksphere)'s provided [enterprise support](https://d2iq.com/services-and-support) (solutions such as [Konvoy](https://d2iq.com/solutions/ksphere/konvoy)), but **support is not given *directly* via this repository**.

Please [contact D2iQ](https://d2iq.com/contact) for questions and more information.

## Releases

Releases are tags which are cut from `master` and other branches and are intended to signify a single supported instance of `kubernetes-base-addons`.

While you may see other tags in our [releases page](https://github.com/mesosphere/kubernetes-base-addons/releases) the only releases which are official and supported releases are those designated as non-prerelease and specifically mentioned to be official releases.

**NOTE**: Do not use `master` for production. Instead, pick a supported release version.

For the release process, see the [release](RELEASE.md) document.

### Creating a Release

Creating a releases is as simple as cutting a tag and making a [Github release](https://help.github.com/en/github/administering-a-repository/creating-releases).

If the release is meant to be an official KSphere supported release, ensure it's not marked as a pre-release, follows the `stable-{MAJOR}.{MINOR}.{PATCH}-{REVISION}` pattern where the major, minor and patch version are that of the corresponding Kubernetes version. Ensure the language in the title and description indicate "Official Release".

Some supported releases are supported via the terms of support for some other KSphere entity using them, particularly releases will be connected with a [Konvoy Release](https://github.com/mesosphere/konvoy). Make sure that you mention and link to any externally related sources for your release in the release description (e.g. if the release is specifically intended to support a specific Konvoy release, say so in the description).

For all other non-official releases, make sure your tag and description are distinctly different from the official release pattern, explain the purpose of your release, and mark is as a `pre-release`.

### Testing

The test suite can be exercised locally by running

    make test


Pull Requests against this repo is tested by Teamcity and Dispatch. 
Dispatchfile defines the config and exercises the test in the Makefile.

## Contributing

See our [Contributing Documentation](CONTRIBUTING.md).
