# Contributing

Contributions to this repository are restricted to D2iQ personnel.

See the [Kubeaddons Contributing Documentation](https://github.com/mesosphere/kubeaddons/blob/master/CONTRIBUTING.md) which provides the baseline information about how to contribute to this repository.

Additionally see the sections below for notes about other rules and considerations for contributions.

## Testing

There are two types of tests associated with this repository.
* [ksphere-testing-framework](https://github.com/mesosphere/ksphere-testing-framework): An integration test built around a golang framework.
* [kubeaddons-tests](https://github.com/mesosphere/kubeaddons-tests): Regression tests built using kuttl.

## Deprecation

Sometimes you may want to **deprecate an older version (or revision) of an Addon**, for instance perhaps you have several revisions of the `v1.x` release of your addon, but now `v2.x` is out and you no longer wish to maintain the older major version.

[Releases](/README.md#Releases) for this repository are responsible for snapshotting collections of revisions and will historically keep your deprecated and unsupported versions, so when you are ready to deprecate any major/minor version of your addon follow these steps:

1. Create a PR that removes all files that are no longer supported
2. Once the PR merges, a new minor release of this repository should be created indicating removal of support for the older versions

Once this is complete end users who are using older releases of `kubernetes-base-addons` will be unaffected.

## Addon Revisions

You will find that any particular addon directory (e.g. `addons/prometheus`) may have several directories and several manifests nested in them each with variants of that addon. These are what we refered to above as "revisions".

The **intention of revisions is to maintain a flat history of addon changes**. If you are making changes to any particular addon you should be making a revision of that addon as a copy of the original file with the changes made therein and the `addon-revision` version updated to reflect the new version appropriately.

New directories can be created for minor versions of a release (e.g. a directory named `v0.9.x`) and contain any revisions matching that minor release. You create new files which should follow the patter `<addon_name>-<version>-<revision>` (e.g. `helloworld-v0.9.1-1.yaml`) keeping in mind revisions start over from 1 again if the patch version updates.

Right now adding revisions is a manual process (see [DCOS-62943](https://jira.mesosphere.com/browse/DCOS-62943) related to automating this in future iterations). To make it easier for reviewers to review PRs where revisions are [manually] added make the first commit in your branch a commit to ONLY copy the previous revision to the new revision file. That way the following commits can actually include your changes and will be easier to historically follow without needing to manually diff the files.
