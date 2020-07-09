# Release Notes

## stable-1.15-2.0.0, stable-1.16-2.0.0, stable-1.17-2.0.0

* \[awsebscsiprovisioner\] The manual steps to upgrade the snapshot APIs from v1alpha1 to v1beta1 is no longer required. It has been automated in the chart CRD install hook by default. If you do not want that default behavior of cleaning up v1alpha1 snapshot CRDs, you can set `cleanupVolumeSnapshotCRDV1alpha1` to `false` and follow the instructions for upgrading to Kubernetes `1.17`. ([#273](https://github.com/mesosphere/kubernetes-base-addons/pull/273), [@sebbrandt87](https://github.com/sebbrandt87))
* \[gcpdisk-csi-driver\] The manual steps to upgrade the snapshot APIs from v1alpha1 to v1beta1 is no longer required. It has been automated in the chart CRD install hook by default. If you do not want that default behavior of cleaning up v1alpha1 snapshot CRDs, you can set `cleanupVolumeSnapshotCRDV1alpha1` to `false` and follow the instructions for upgrading to Kubernetes `1.17`.
  \[azuredisk-csi-driver\] The manual steps to upgrade the snapshot APIs from v1alpha1 to v1beta1 is no longer required. It has been automated in the chart CRD install hook by default. If you do not want that default behavior of cleaning up v1alpha1 snapshot CRDs, you can set `snapshot.cleanupVolumeSnapshotCRDV1alpha1` to `false` and follow the instructions for upgrading to Kubernetes `1.17`. ([#279](https://github.com/mesosphere/kubernetes-base-addons/pull/279), [@jieyu](https://github.com/jieyu))
* \[prometheus-operator\] Upgrade to version [0.38.1](https://github.com/coreos/prometheus-operator/releases/tag/v0.38.1)
    - \[prometheus\] Upgrade to version [2.17.2](https://github.com/prometheus/prometheus/releases/tag/v2.17.2)
    - \[grafana\] Upgrade to version [6.7.3](https://github.com/grafana/grafana/releases/tag/v6.7.3) ([#281](https://github.com/mesosphere/kubernetes-base-addons/pull/281), [@branden](https://github.com/branden))
* \[traefik\] fix an issue where `clusterhostname` can now be an ipaddress as well ([#286](https://github.com/mesosphere/kubernetes-base-addons/pull/286), [@GoelDeepak](https://github.com/GoelDeepak))
* [dex-k8s-authenticator] Fix bug in init container that could remove custom CA certificate from main cluster login instructions ([#291](https://github.com/mesosphere/kubernetes-base-addons/pull/291), [@mhrabovcin](https://github.com/mhrabovcin))
* \[traefik\] Distribute pods across nodes and zones when possible.
  \[traefik\] Set a PodDisruptionBudget to ensure at least 1 pod is running at all times. ([#292](https://github.com/mesosphere/kubernetes-base-addons/pull/292), [@branden](https://github.com/branden))
* Prometheus-alert-manager: increase memory and cpu limits due to OOM errors ([#298](https://github.com/mesosphere/kubernetes-base-addons/pull/298), [@hectorj2f](https://github.com/hectorj2f))
* Traefik is now upgradeable again when the `initCertJobImage` field is modified. ([#302](https://github.com/mesosphere/kubernetes-base-addons/pull/302), [@makkes](https://github.com/makkes))
* \[traefik\]:
  - upgrade to 1.7.24 
  - mTLS available
  - accessLogs.filters setable
  - caServer setable for acme challenge ([#304](https://github.com/mesosphere/kubernetes-base-addons/pull/304), [@sebbrandt87](https://github.com/sebbrandt87))
* Traefik: access log is enabled by default ([#305](https://github.com/mesosphere/kubernetes-base-addons/pull/305), [@mhrabovcin](https://github.com/mhrabovcin))
* Opsportal: fix a typo in 'lables' that caused issues during upgrades. ([#307](https://github.com/mesosphere/kubernetes-base-addons/pull/307), [@dkoshkin](https://github.com/dkoshkin))
* \[prometheus\]: Update prometheus-operator chart, which adds a grafana dashboard for monitoring autoscaler ([#308](https://github.com/mesosphere/kubernetes-base-addons/pull/308), [@gracedo](https://github.com/gracedo))
* \[dex-k8s-authenticator\]:
  - fix: render configure kubectl instructions with the cluster hostname. 
  - fix: add clippy js for clipboard support ([#309](https://github.com/mesosphere/kubernetes-base-addons/pull/309), [@samvantran](https://github.com/samvantran))
* \[prometheus\] Increases default Prometheus server resources. ([#310](https://github.com/mesosphere/kubernetes-base-addons/pull/310), [@branden](https://github.com/branden))
* ValuesRemap has been added for rewriting the forward authentication url in multiple addons. ([#315](https://github.com/mesosphere/kubernetes-base-addons/pull/315), [@jr0d](https://github.com/jr0d))
* Konvoyconfig has a new field `caCertificate` to support custom certificate in managed cluster ([#316](https://github.com/mesosphere/kubernetes-base-addons/pull/316), [@GoelDeepak](https://github.com/GoelDeepak))
* Istio addon upgraded to 1.6.3 ([#317](https://github.com/mesosphere/kubernetes-base-addons/pull/317), [@GoelDeepak](https://github.com/GoelDeepak))
* Opsportal: allow landing page deployment replica count to be configured ([#319](https://github.com/mesosphere/kubernetes-base-addons/pull/319), [@jieyu](https://github.com/jieyu))
* \[dashboard\] Upgrades the Kubernetes dashboard to 2.0.3.
  \[dashboard\] Adds metrics visualizations to the Kubernetes dashboard UI. ([#320](https://github.com/mesosphere/kubernetes-base-addons/pull/320), [@branden](https://github.com/branden))
* Traefik: revert changes to the service ports that broke Velero functionality. ([#328](https://github.com/mesosphere/kubernetes-base-addons/pull/328), [@dkoshkin](https://github.com/dkoshkin))
* Traefik-foward-auth: fix a bug that might cause /_oauth callback to be redirected to other services ([#334](https://github.com/mesosphere/kubernetes-base-addons/pull/334), [@jieyu](https://github.com/jieyu))
* Adds the Conductor service card to the cluster detail page of the UI. ([#344](https://github.com/mesosphere/kubernetes-base-addons/pull/344), [@natmegs](https://github.com/natmegs))

## stable-1.15-1.8.0, stable-1.16-1.8.0

* \[kibana\]: Fixes an issue causing an outdated version of Kibana to be deployed to GCP. ([#249](https://github.com/mesosphere/kubernetes-base-addons/pull/249), [@branden](https://github.com/branden))

## stable-1.15-1.7.0, stable-1.16-1.7.0

*  \[prometheus\]
   * \[CHANGE\] Restrict api extension RBAC rules
   * \[BUGFIX\] Fix statefulset crash loop on kubernetes ([#219](https://github.com/mesosphere/kubernetes-base-addons/pull/219), [@shaneutt](https://github.com/shaneutt))
* \[dex\]: support specifying root CA for LDAP connectors in Dex controller. ([#224](https://github.com/mesosphere/kubernetes-base-addons/pull/224), [@jieyu](https://github.com/jieyu))
* \[velero\]: bump velero to chart version 3.0.3, which includes velero-minio RELEASE.2020-04-10T03-34-42Z ([#215](https://github.com/mesosphere/kubernetes-base-addons/pull/215), [@jieyu](https://github.com/jieyu))
* \[dex-k8s-authenticator\] added support for the konvoy credentials plugin ([#193](https://github.com/mesosphere/kubernetes-base-addons/pull/193), [@jr0d](https://github.com/jr0d))
* \[velero\]: switch minio backend logging from plaintext to json ([#216](https://github.com/mesosphere/kubernetes-base-addons/pull/216), [@vespian](https://github.com/vespian))

## stable-1.15-1.6.0, stable-1.16-1.6.0

* \[dex-k8s-authenticator\]: Now supports a kubectl credentials plugin for automatically managing identity tokens. Instructions for downloading the plugin and configuring kubectl can be found at `https://<cluster-ip>/token/plugin`. ([#212](https://github.com/mesosphere/kubernetes-base-addons/pull/212), [@jr0d](https://github.com/jr0d))
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
