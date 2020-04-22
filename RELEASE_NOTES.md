# Release Notes

## stable-1.15-1.6.0, stable-1.16-1.6.0

* \[cert-manager\] `usages` is no longer definable as part of `issuerRef`, instead it is a key on its own ([#196](https://github.com/mesosphere/kubernetes-base-addons/pull/196), [@sebbrandt87](https://github.com/sebbrandt87))
* \[elasticsearch\] Fixes an issue that may cause the elasticsearch addon to fail to deploy. ([#206](https://github.com/mesosphere/kubernetes-base-addons/pull/206), [@branden](https://github.com/branden))

## stable-1.15-1.5.0, stable-1.16-1.5.0
* \[Elasticsearch\] revert the PVC size to default (30G) for data nodes ([#203](https://github.com/mesosphere/kubernetes-base-addons/pull/203), [@jieyu](https://github.com/jieyu))
* \[Prometheus\] Upgrade prometheus-operator chart to v8.8.4 ([#205](https://github.com/mesosphere/kubernetes-base-addons/pull/205), [@joejulian](https://github.com/joejulian))
* \[awsebscsiprovisioner\] Upgrade awsebscsiprovisioner chart to 0.3.5 and aws-ebs-csi-driver to 0.5.0. ([#186](https://github.com/mesosphere/kubernetes-base-addons/pull/186), [@sebbrandt87](https://github.com/sebbrandt87))
* \[kube-oidc-proxy\] allow using default system CA bundle. ([#191](https://github.com/mesosphere/kubernetes-base-addons/pull/191), [@jieyu](https://github.com/jieyu))
* \[Traefik\] Upgrade Traefik to 1.7.23. This change fixes the ability to access the Kubernetes API server when the connection needs to be upgraded to SPDY, among other bug fixes. For more details, see https://github.com/mesosphere/charts/pull/514. ([#190](https://github.com/mesosphere/kubernetes-base-addons/pull/190), [@joejulian](https://github.com/joejulian))
* \[dex-k8s-authenticator\] allow to use system default CA ([#189](https://github.com/mesosphere/kubernetes-base-addons/pull/189), [@jieyu](https://github.com/jieyu))
* \[Istio\] Disable Istio PodDisruptionBudget, the default settings and replica count of 1 prevents pods on nodes from being drained. ([#183](https://github.com/mesosphere/kubernetes-base-addons/pull/183), [@dkoshkin](https://github.com/dkoshkin))

## stable-1.15-1.4.1, stable-1.16-1.4.1

* \[Velero\] revert the velero refactor in stable-1.16-1.4.0 due to a data loss issue ([#197](https://github.com/mesosphere/kubernetes-base-addons/pull/197), [@jieyu](https://github.com/jieyu))
* \[Velero-minio\] fix a data loss issue after upgrade ([#200](https://github.com/mesosphere/kubernetes-base-addons/pull/200), [@jieyu](https://github.com/jieyu))

## stable-1.15-1.4.0, stable-1.16-1.4.0

* \[Dex\] Add SAML connector support in dex controller allowing users to add SAML IDP using Kubernetes API. ([#173](https://github.com/mesosphere/kubernetes-base-addons/pull/173), [@jieyu](https://github.com/jieyu))
* \[Velero\] switch to use minio helm chart (instead of operator) for backup storage. This allow users to install their own minio operator for general purpose object storage. ([#174](https://github.com/mesosphere/kubernetes-base-addons/pull/174), [@jieyu](https://github.com/jieyu))

## stable-1.15-1.3.0, stable-1.16-1.3.0

* \[ElasticSearch, fluentbit\] Create index template
  Create ElasticSearch Index Template. Require Fluentbit to deploy only after ElasticSearch deploys.

## stable-1.15-1.2.0, stable-1.16-1.2.0

* fluent-bit
  * Disable audit log  collection
    It's been observed in production clusters that the audit log bloats the number of fields in an index.
    This causes resource limits to be filled and throttling to occur.
    We are disabling this collection pending further investigation.
* dex:
  * improve the LDAP connector validation in Dex controller
  * fix an issue in dex addon which disallowed adding local users
  * use Dex controller v0.4.1, which includes the support for OIDC group claims
  * upgrade Dex to v2.22.0, which supports groups claims for OIDC connectors
* dex-k8s-authenticator: 
  * allow scopes to be configured, and drop the `offline_access` scope as it is not used
* kube-oidc-proxy:
  *  enable token passthrough
* opsportal:
  * set `opsportalRBAC.allowAllAuthenticated` to true
  * add RBAC support
* traefik-forward-auth:
  * enable RBAC and impersonation
  * remove whitelisting
* kibana:
  * upgrade to 6.8.2
* elasticsearch-curator:
  * added and enabled curator to remove old indexes from elasticsearch to free up storage


Add support for kubernetes clusters on GCP
Various chart bumps for stability, bug and security fixes.
