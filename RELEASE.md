# Release

The kubernetes-base-addons repository is intended to support a minimal base set of Addons required to provide monitoring, logging, alerting, and backups for every version of Kubernetes actively supported by the upstream Kubernetes community.
It is not intended to be tied to any specific version of Kubernetes installer (ie. Konvoy or Kommander), or Kubeaddons.
Changes to Konvoy or Kommander to create a resource needed by an Addon should be avoided.

- [Schedule](#schedule)
- [Considerations](#considerations)
  - [Kubernetes](#kubernetes)
  - [Branches and Tags](#branches-and-tags)
- [Process](#process)
  - [Testing Release (Second and Forth Thursday)](#testing-release-second-and-forth-thursday)
  - [Stable Release (Second and Forth Wednesday)](#stable-release-second-and-forth-wednesday)

## Schedule

Releases should be twice monthly<sup>[1](#footnote1)</sup> on the second and forth Wednesdays, or as needed to address CVEs.

## Considerations

### Kubernetes

**The Kubernetes versions supported by this repo are 1.15, 1.16, and 1.17.**
Support should be maintained for the latest general release of Kubernetes and the two prior minor releases.

The Addons in this repo are the minimum base set of supported<sup>[2](#footnote2)</sup> Addons to be installed as a suite for any supported version of Kubernetes.

### Branches and Tags

As much as possible, we will try to maintain a `master` branch that is compatible with all supported versions of Kubernetes.
If there is a variation in the Kubernetes API which requires a _**breaking change**_ to the Addon, a _**branch**_ will be made for the prior Kubernetes versions.
All future changes adopted into master will need to be back-ported to those branches.

**NOTE**: No other changes may be breaking. Extreme changes, like moving from traefik 1.7 to 2.2, must be done in such a way that the transition is transparent to the user.

## Process

### Testing Release (Second and Forth Thursday)

On the second and forth _**Thursday**_, this repository should be branched for SOAK testing by setting the `testing` branch to the head of master and force-pushing.

- This set of Addons are installed into the SOAK cluster.

### Stable Release (Second and Forth Wednesday)

On the second and forth _**Wednesday**_:

- Using automation, parse the PR logs in the `testing` branch for release notes and generate and commit a Changelog.md
- One _**tag**_ is made for the each supported version of Kubernetes with a consistent semver suffix
- If the current version of Kubernetes is `1.17.2`, and the last release was `stable-1.17-1.2.x`, the new SOAK tag will be `release-1.17-1.3.0`
- The same semver portion of the tag `1.3.0` is used for each supported kubernetes version, the new tags being `release-1.16-1.3.0`, and `release-1.15-1.3.0`
- These tag versions, in the form of `release-<major>.<minor>-<semver>`, only the `<major>.<minor>` refer to the kubernetes version.
  The api within a minor version should not be changing so there should never be a need to refer to the kubernetes patch version.
- As a standing agenda item in sig-ksphere-catalog, vote go/no-go on the release of the Addons that have been SOAK tested.
- Merge the `testing` branch into the `stable` branch<sup>[3](#footnote3)</sup>
- Announce the release.

<a name="footnote1">1</a>: Based on a two-week soak cycle. If we can have overlapping soak clusters, we can accelerate this.

<a name="footnote2">2</a>: A supported Addon is one which has been tested to work in concert with other Addons in the same release. This suite of Addons, as a whole, constitute a set for which D2iQ customers can get support with their software contract. Variations from the configurations and suite of Addons are not expected to be the responsibility of D2iQ support.

<a name="footnote3">3</a>: In the future, there may need to be multiple `stable` branches as needed to maintain our support commitment.

**NOTE:** This document is governed by kep sig-ksphere-cluster/20200218-kubernetes-base-addon-release-process.md