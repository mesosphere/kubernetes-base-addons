# Kubernetes Base Addons

This repository contains the default sets of addons needed to bootstrap [D2iQ Kubernetes](https://d2iq.com/solutions/ksphere).

# Overview

The structure of this repository follows the [Kubeaddons Catalog Documentation](https://github.com/mesosphere/kubeaddons/blob/master/tools/catalog/README.md) in reference, and uses the [Addon Revision concept covered therein](https://github.com/mesosphere/kubeaddons/blob/master/tools/catalog/README.md#special-addonrepository-options---addon-revisions).

You will find the following directories here:

* `addons/` - containing the actual manifests for addon resources
* `deployments/` - containing the default addons depending on the Kubernetes version
* `metadata/` - containing static metadata for the addons in `addons/`
* `test/` - containing integration tests for the addons in `addons/`

# Support

There is **no community support for this software at this time**, and this repository (while publicly visible) is only intended for consumption by D2iQ personnel.

If you are here looking for addon support as an end-user, support is given as per the support provided for your [D2iQ Ksphere Solution](https://d2iq.com/solutions/ksphere)'s provided [enterprise support](https://d2iq.com/services-and-support) (solutions such as [Konvoy](https://d2iq.com/solutions/ksphere/konvoy)), but **support is not given *directly* via this repository**.

Please [contact D2iQ](https://d2iq.com/contact) for questions and more information.

# Releases

Releases are tags which are cut from `master` and other branches and are intended to signify a single supported instance of `kubernetes-base-addons`.

While you may see other tags in our [releases page](https://github.com/mesosphere/kubernetes-base-addons/releases) the only releases which are official and supported releases are those designated as non-prerelease and specifically mentioned to be official releases.

**WARNING**: Do not use `master` for production, instead pick a supported release version.

## Creating a Release

Creating a releases is as simple as cutting a tag and making a [Github release](https://help.github.com/en/github/administering-a-repository/creating-releases).

If the release is meant to be an official KSphere supported release, ensure it's not marked as a pre-release, follows the `stable-{MAJOR}.{MINOR}.{PATCH}-{REVISION}` pattern where the major, minor and patch version are that of the corresponding Kubernetes version. Ensure the language in the title and description indicate "Official Release".

Some supported releases are supported via the terms of support for some other KSphere entity using them, particularly releases will be connected with a [Konvoy Release](https://github.com/mesosphere/konvoy). Make sure that you mention and link to any externally related sources for your release in the release description (e.g. if the release is specifically intended to support a specific Konvoy release, say so in the description).

For all other non-official releases, make sure your tag and description are distinctly different from the official release pattern, explain the purpose of your release, and mark is as a `pre-release`.

# Contributing

Contributions to this repository are restricted to D2iQ personnel.

See the [Kubeaddons Contributing Documentation](https://github.com/mesosphere/kubeaddons/blob/master/CONTRIBUTING.md) which provides the baseline information about how to contribute to this repository.

Additionally see the sections below for notes about other rules and considerations for contributions.

## Deprecation

Sometimes you may want to **deprecate an older version (or revision) of an Addon**, for instance perhaps you have several revisions of the `v1.x` release of your addon, but now `v2.x` is out and you no longer wish to maintain the older major version.

[Releases](/README.md#Releases) for this repository are responsible for snapshotting collections of revisions and will historically keep your deprecated and unsupported versions, so when you are ready to deprecate any major/minor version of your addon follow these steps:

1. Create a PR that removes all files that are no longer supported
2. Once the PR merges, a new minor release of this repository should be created indicating removal of support for the older versions

Once this is complete end users who are using older releases of `kubernetes-base-addons` will be unaffected.

## Addon Revisions

You will find that any particular addon directory (e.g. `addons/prometheus`) may have several directories and several manifests nested in them each with variants of that addon. These are what we refered to above as "revisions".

The **intention of revisions is to maintain a flat history of addon changes**. If you are making changes to any particular addon you should be making a revision of that addon as a copy of the original file with the changes made therein and the `addon-revision` version updated to reflect the new version appropriately.

