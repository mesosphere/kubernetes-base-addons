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

Releases should be at the end of every sprint.

## Considerations

### Kubernetes

**The Kubernetes versions supported by this repo are 1.17, 1.18, and 1.19.**
Support should be maintained for the latest general release of Kubernetes and the two prior minor releases.

The Addons in this repo are the minimum base set of supported<sup>[2](#footnote2)</sup> Addons to be installed as a suite for any supported version of Kubernetes.

### Branches and Tags

This project is a dependency of DKP and releases will be branched and maintained to facilitate the releases of Konvoy.

This repo will be branched, as necessary, to support the release cycle of Konvoy.
There will be no attempt to provide a single repo that can be used across all supported versions of Kubernetes, instead we will tie branching to the versions of the addons within the repo.
If any addon has a major version change, this repo will have a major version change.
If minor, minor, if only patch versions change, only the patch version will be bumped.
There will be no effort to tie the versions to the versions of any one addon.

All changes adopted into master will need to determine if backports to prior branches are necessary to support Konvoy releases.

**NOTE**: No other changes may be breaking. Extreme changes, like moving from traefik 1.7 to 2.2, must be done in such a way that the transition is transparent to the user.

## Process

### Testing Release (The beginning of every sprint)

At the beginning of every sprint, this repository should be branched for SOAK testing. Each branch tied to each konvoy release will be tested on the appropriate SOAK cluster.

### Stable Release (Second and Forth Wednesday)

At the end of each sprint:

- Create a PR branch for the `release/*` branch being released.
- Run `make release.pr KBA_MILESTONE=release/3.2 KBA_TAGS=v3.2.0,stable-1.19-3.2.0,stable-1.16-3.2.0`
  - Use the appropriate `milestone` and `tags` using semver, and konvoy-specific tags as demonstrated above.
  - One _**tag**_ is made for the semver of this release, and for each supported version of Kubernetes with a consistent semver suffix, eg: `v2.3.0`, `stable-1.16-2.3.0`, `stable-1.17-2.3.0`, etc.
  - These tag versions, in the form of `release-<major>.<minor>-<semver>`, only the `<major>.<minor>` refer to the kubernetes version.
    The api within a minor version should not be changing so there should never be a need to refer to the kubernetes patch version.
- Commit and open a PR to the correct release branch
- Once merged, `make release KBA_MILESTONE=release/3.2 KBA_TAGS=v3.2.0,stable-1.19-3.2.0,stable-1.16-3.2.0`
- Announce the release.

<a name="footnote1">1</a>: Based on a two-week soak cycle. If we can have overlapping soak clusters, we can accelerate this.

<a name="footnote2">2</a>: A supported Addon is one which has been tested to work in concert with other Addons in the same release. This suite of Addons, as a whole, constitute a set for which D2iQ customers can get support with their software contract. Variations from the configurations and suite of Addons are not expected to be the responsibility of D2iQ support.

**NOTE:** This document is governed by kep sig-ksphere-cluster/20200218-kubernetes-base-addon-release-process.md
