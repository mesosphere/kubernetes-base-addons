# Contributing

Contributions to this repository are restricted to D2iQ personnel.

See the [Kubeaddons Contributing Documentation](https://github.com/mesosphere/kubeaddons/blob/master/CONTRIBUTING.md) which provides the baseline information about how to contribute to this repository.

Additionally see the sections below for notes about other rules and considerations for contributions.

## Testing

There are three types of tests associated with this repository.
* [ksphere-testing-framework](https://github.com/mesosphere/ksphere-testing-framework): An integration test built around a golang framework.
* [kubeaddons-tests](https://github.com/mesosphere/kubeaddons-tests): Regression tests built using kuttl.
* [E2E tests](https://github.com/mesosphere/kommander/tree/master/system-tests#system-tests): E2E tests that spin up AWS clusters and deploy projects with addons.

## Deprecation

Sometimes you may want to **deprecate an older version (or revision) of an Addon**, for instance perhaps you have several revisions of the `v1.x` release of your addon, but now `v2.x` is out and you no longer wish to maintain the older major version.

[Releases](/README.md#Releases) for this repository are responsible for snapshotting collections of revisions and will historically keep your deprecated and unsupported versions, so when you are ready to deprecate any major/minor version of your addon follow these steps:

1. Create a PR that removes all files that are no longer supported
2. Once the PR merges, a new minor release of this repository should be created indicating removal of support for the older versions

Once this is complete end users who are using older releases of `kubernetes-base-addons` will be unaffected.

## Addon Revisions

For Addons, "revisions" are a reference to the application version with an added revision count which indicates the latest iteration on that version. This enables multiple different versions of an Addon which ultimately utilize the same underlying application version so that configuration and other aspects of the Addon can change without overriding a previously released Addon.

The [Kubeaddons Catalog API](https://github.com/mesosphere/kubeaddons/tree/master/pkg/catalog) supports two different modes for addon repositories to host revisions:

* single file mode: a single file for the addon exists at `addons/<addon-name>/<addon-name>.yaml` and the revision for that addon must be adjusted forward when any changes are made
* multi file mode: each new revision is a separate file, and you can structure this like `addons/<addon-name>/<app-version>/<addon-name>-<revision>.yaml`. However, there's no obligation to follow this pattern as the Kubeaddons controller will consider every valid Addon manifest it finds in the directory hierarchy.

In this repository we use single file mode because that Catalog API is not used in such a way that we need to bother searching the historical revisions (this repository is predominantly used by Konvoy which does it's versioning for Addons based on Git).

You may find in other repositories (such as [mesosphere/kubeaddons-enterprise](https://github.com/mesosphere/kubeaddons-enterprise)) that multi-file mode is used to support flat searching of the Addon revision history.
