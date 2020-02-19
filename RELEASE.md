# Release

The kubernetes-base-addons repository is intended to support a minimal base set of Addons required to provide monitoring, logging, alerting, and backups for every version of Kubernetes actively supported by the upstream Kubernetes community.
It is not intended to be tied to any specific version of Kubernetes installer (ie. Konvoy or Kommander), or Kubeaddons.
Changes to Konvoy or Kommander to create a resource needed by an Addon should be avoided.

- [Schedule](#schedule)
- [Considerations](#considerations)
  - [Kubernetes](#kubernetes)
  - [Branches and Tags](#branches-and-tags)
- [Process](#process)
  - [Testing Release (Weekly, Thursday)](#testing-release-weekly-thursday)
  - [Stable Release (Biweekly, Wednesday)](#stable-release-biweekly-wednesday)

## Schedule

Releases should be bi-weekly<sup>[1](#footnote1)</sup> on Wednesdays, or as needed to address CVEs.

## Considerations

### Kubernetes

The Addons in this repo are the minimum base set of supported<sup>[2](#footnote2)</sup> Addons to be installed as a suite for any supported version of Kubernetes.

Kubernetes support must be maintained for the latest general release of Kubernetes and the two prior minor releases.
At the time of this writing, the latest Kubernetes release is 1.17.2.
Support from this repo, at this time, should cover 1.15.x, 1.16.x, and 1.17.x.

### Branches and Tags

As much as possible, we will try to maintain a `master` branch that is compatible with all supported versions of Kubernetes.
If there is a variation in the Kubernetes API which requires a _**breaking change**_ to the Addon, a _**branch**_ will be made for the prior Kubernetes versions.
All future changes adopted into master will need to be back-ported to those branches.

**NOTE**: No other changes may be breaking. Extreme changes, like moving from traefik 1.7 to 2.2, must be done in such a way that the transition is transparent to the user.

## Process

### Testing Release (Weekly, Thursday)

Every other _**Thursday**_, this repository should be tagged for SOAK testing as follows:

- Using automation, parse the PR logs for release notes and generate and commit a Changelog.md
- One _**tag**_ is made for the each supported version of Kubernetes with an incremented release counter
- If the current version of Kubernetes is `1.17.2`, and the last release was `stable-1.17-1.5.x`, the new SOAK tag will be `testing-1.17-1.6.0`.
- For the previous Kubernetes version, the last release may have been `stable-1.16-1.9.x`. The new tag for this Kubernetes version is `testing-1.16-1.10.0`.
- The oldest supported release similarly might be `testing-1.15-1.27.0` for a prior `stable-1.15-1.26.x`.
- These tag versions, in the form of `major.minor-semver`, only the `major.minor` refer to the kubernetes version. The api within a minor version should not be changing so there should never be a need to refer to the kubernetes patch version.

**NOTE:** If a breaking change causes a diversion from an older release of Kubernetes to a newer one, prior to tagging the older version must be branched, ie. `stable-1.16-1.9.0` would become `stable-1.16`, the changes since the `stable-1.16-1.9.0` tag would be merged into this branch, and the new _tag_ would still be `testing-1.16-1.10.0` but pointing to the last change on the `stable-1.16` branch.

- This set of Addons are installed into a SOAK cluster<sup>[3](#footnote3)</sup>.

### Stable Release (Biweekly, Wednesday)

Every other _**Wednesday**_, the day before the next testing release:

- As a standing agenda item in sig-ksphere-catalog, vote go/no-go on the release of the Addons that have been SOAK tested.
- Create `stable-` tags for the `testing-` tags that ran in SOAK.
- Announce the release.

<a name="footnote1">1</a>: Based on a two-week soak cycle. If we can have overlapping soak clusters, we can accelerate this.

<a name="footnote2">2</a>: A supported Addon is one which has been tested to work in concert with other Addons in the same release. This suite of Addons, as a whole, constitute a set for which D2iQ customers can get support with their software contract. Variations from the configurations and suite of Addons are not expected to be the responsibility of D2iQ support.

<a name="footnote3">3</a>: At the time of this writing, this process is as yet undetermined as there are no clusters in which to do this SOAK.
