# Release Notes

## stable-1.16-2.5.0, stable-1.17-2.5.0, stable-1.18-2.5.0

* Cert-manager
  - use `delete` upgrade-strategy.
* Fluent-bit:
  - bump the fluent-bit app version to 1.5.6
    - aws: utils: fix mem leak in flb_imds_request
    - fix double free when destroying connections if the endpoint in unavailable
    - remove noisy error introduced in v1.5.5
    - fix deletion of pending connections in the destroy_queue ([#538](https://github.com/mesosphere/kubernetes-base-addons/pull/538), [@d2iq-dispatch](https://github.com/d2iq-dispatch))
  - changes the update strategy to `delete`. ([#574](https://github.com/mesosphere/kubernetes-base-addons/pull/574), [@dkoshkin](https://github.com/dkoshkin))
  - Upgrades fluent-bit to v1.5.7. See https://fluentbit.io/announcements/v1.5.7.
  - Adds chart value `podLabels`.
  - Fix fluentbit configuration to unblock output buffer. ([#590](https://github.com/mesosphere/kubernetes-base-addons/pull/590), [@branden](https://github.com/branden))
* Kibana
  - Fixes an issue that causes Kibana to deploy without an audit log dashboard. ([#511](https://github.com/mesosphere/kubernetes-base-addons/pull/511), [@branden](https://github.com/branden))

### Preview
* Istio
  - Bug Fixes
    - Fixed HTTP match request without headers conflict
    - Fixed Istio operator to watch multiple namespaces (Istio &#35;26317)
    - Fixed EDS cache when an endpoint appears after its service resource (Istio &#35;26983)
    - Fixed istioctl remove-from-mesh not removing init containers on CNI installations.
    - Fixed istioctl add-to-mesh and remove-from-mesh commands from affecting OwnerReferences (Istio &#35;26720)
    - Fixed cleaning up of service information when the cluster secret is deleted
    - Fixed egress gateway ports binding to 80⁄443 due to user permissions
    - Fixed gateway listeners created with traffic direction outbound to be drained properly on exit
    - Fixed headless services not updating listeners (Istio &#35;26617)
    - Fixed inaccurate endpointsPendingPodUpdate metric
    - Fixed ingress SDS from not getting secret update (Istio &#35;18912)
    - Fixed ledger capacity size
    - Fixed operator to update service monitor due to invalid permissions (Istio &#35;26961)
    - Fixed regression in gateway name resolution (Istio 26264)
    - Fixed rotated certificates not being stored to /etc/istio-certs VolumeMount (Istio &#35;26821)
    - Fixed trust domain validation in transport socket level (Istio &#35;26435)
  - Improvements
    - Added istioctl analyzer to detect when Destination Rules do not specify caCertificates (Istio &#35;25652)
    - Added missing telemetry.loadshedding.- options to mixer container arguments
    - Improved specifying network for a cluster without meshNetworks also being configured
    - Improved the cache readiness state with TTL (Istio &#35;26418)
    - Updated SDS timeout to fetch workload certificates to 0s
    - Updated app_containers to use comma separated values for container specification
    - Updated default protocol sniffing timeout to 5s (Istio &#35;24379) ([#516](https://github.com/mesosphere/kubernetes-base-addons/pull/516), [@shaneutt](https://github.com/shaneutt))

## stable-1.15-2.4.0, stable-1.16-2.4.0, stable-1.17-2.4.0

* Istio:
  - The "kubernetes-service-monitor" service monitor has been removed. ([#481](https://github.com/mesosphere/kubernetes-base-addons/pull/481), [@gracedo](https://github.com/gracedo))

  - Bumped Istio to v1.6.8:
    - Fixed security issues:
      - CVE-2020-12603: By sending a specially crafted packet, an attacker could cause Envoy to consume excessive amounts of memory when proxying HTTP/2 requests or responses.
      - CVE-2020-12605: An attacker could cause Envoy to consume excessive amounts of memory when processing specially crafted HTTP/1.1 packets.
      - CVE-2020-8663: An attacker could cause Envoy to exhaust file descriptors when accepting too many connections.
      - CVE-2020-12604: An attacker could cause increased memory usage when processing specially crafted packets.
      - CVE-2020-15104: When validating TLS certificates, Envoy incorrectly allows a wildcard DNS Subject Alternative Name to apply to multiple subdomains. For example, with a SAN of   .example.com, Envoy incorrectly allows nested.subdomain.example.com, when it should only allow subdomain.example.com.
      - CVE-2020-16844: Callers to TCP services that have a defined Authorization Policies with DENY actions using wildcard suffixes (e.g. *-some-suffix) for source principals or namespace fields will never be denied access.
    - Other changes:
      - Fixed return the proper source name after Mixer does a lookup by IP if multiple pods have the same IP.
      - Improved the sidecar injection control based on revision at a per-pod level (Issue 24801)
      - Improved istioctl validate to disallow unknown fields not included in the Open API specification (Issue 24860)
      - Changed stsPort to sts_port in Envoy’s bootstrap file.
      - Preserved existing WASM state schema for state objects to reference it later as needed.
      - Added targetUri to stackdriver_grpc_service.
      - Updated WASM state to log for Access Log Service.
      - Increased default protocol detection timeout from 100 ms to 5 s (Issue 24379)
      - Removed UDP port 53 from Istiod.
      - Allowed setting status.sidecar.istio.io/port to zero (Issue 24722)
      - Fixed EDS endpoint selection for subsets with no or empty label selector. (Issue 24969)
      - Allowed k8s.overlays on BaseComponentSpec. (Issue 24476)
      - Fixed istio-agent to create elliptical curve CSRs when ECC_SIGNATURE_ALGORITHM is set.
      - Improved mapping of gRPC status codes into HTTP domain for telemetry.
      - Fixed scaleTargetRef naming in HorizontalPodAutoscaler for Istiod (Issue 24809)
      - Optimized performance in scenarios with large numbers of gateways. (Issue 25116)
      - Fixed an issue where out of order events may cause the Istiod update queue to get stuck. This resulted in proxies with stale configuration.
      - Fixed istioctl upgrade so that it no longer checks remote component versions when using --dry-run. (Issue 24865)
      - Fixed long log messages for clusters with many gateways.
      - Fixed outlier detection to only fire on user configured errors and not depend on success rate. (Issue 25220)
      - Fixed demo profile to use port 15021 as the status port. (Issue &#35;25626)
      - Fixed Galley to properly handle errors from Kubernetes tombstones.
      - Fixed an issue where manually enabling TLS/mTLS for communication between a sidecar and an egress gateway did not work. (Issue 23910)
      - Fixed Bookinfo demo application to verify if a specified namespace exists and if not, use the default namespace.
      - Added a label to the pilot_xds metric in order to give more information on data plane versions without scraping the data plane.
      - Added CA_ADDR field to allow configuring the certificate authority address on the egress gateway configuration and fixed the istio-certs mount secret name.
      - Updated Bookinfo demo application to latest versions of libraries.
      - Updated Istio to disable auto mTLS when sending traffic to headless services without a sidecar.
      - Fixed an issue which prevented endpoints not associated with pods from working. (Issue &#35;25974) ([#489](https://github.com/mesosphere/kubernetes-base-addons/pull/489), [@shaneutt](https://github.com/shaneutt))

* Traefik-forward-auth:
  - Update traefik-foward-auth to 0.2.14
  - Add an option to bypass tfa deployment ([#456](https://github.com/mesosphere/kubernetes-base-addons/pull/456), [@d2iq-dispatch](https://github.com/d2iq-dispatch))

* Fixed an upgrade issue for several addons which would cause them to not be properly targeted for upgrade ([#492](https://github.com/mesosphere/kubernetes-base-addons/pull/492), [@shaneutt](https://github.com/shaneutt))

## stable-1.15-2.3.0, stable-1.16-2.3.0, stable-1.17-2.3.0

- Azuredisk-csi-driver:
  - enable the Snapshot controller ([#443](https://github.com/mesosphere/kubernetes-base-addons/pull/443), [@dkoshkin](https://github.com/dkoshkin))
- Cert-manager:
  - `Issuer` namespace setable
    - `Certificate` namespace setable ([#378](https://github.com/mesosphere/kubernetes-base-addons/pull/378), [@sebbrandt87](https://github.com/sebbrandt87))
- Dex-k8s-authenticator:
  - Windows download support for the credentials plugin ([#377](https://github.com/mesosphere/kubernetes-base-addons/pull/377), [@jr0d](https://github.com/jr0d))
  - Fixed bug causing `certificate-authority=`  option to be added to token instructions on the windows tab when it should have been omitted. ([#436](https://github.com/mesosphere/kubernetes-base-addons/pull/436), [@jr0d](https://github.com/jr0d))
- Elasticsearch-curator:
  - version 5.8.1 ([#374](https://github.com/mesosphere/kubernetes-base-addons/pull/374), [@sebbrandt87](https://github.com/sebbrandt87))
  - Added value `cronjob.startingDeadlineSeconds`: Amount of time to try reschedule job if we can't run on time ([#447](https://github.com/mesosphere/kubernetes-base-addons/pull/447), [@d2iq-dispatch](https://github.com/d2iq-dispatch))
- Elasticsearch-exporter:
  - updated from 2.11 to 3.7.0
    - Add a parameter for the elasticsearch-exporter: es.indices_settings as it is supported since version 1.0.4 (the elasticsearch-exporter chart is supporting the version 1.1.0)
    - Update description for envFromSecret parameter in readme
    - Feature flap the flag es.uri to allow fallback to env var ES_URI
    - Allow setting environment variables with k8s secret information to support referencing already existing sensitive parameters.
    - Add es.ssl.client.enabled value for better functionality readability
    - Add option to disable client cert auth in Elasticsearch exporter
    - Add the serviceMonitor targetLabels key as documented in the Prometheus Operator API
    - Add log.level and log.format configs
    - Add the ServiceMonitor metricRelabelings key as documented in the Prometheus Operator API
    - Add sampleLimit configuration option ([#449](https://github.com/mesosphere/kubernetes-base-addons/pull/449), [@d2iq-dispatch](https://github.com/d2iq-dispatch))
- Fluent-bit:
  - Three different elasticsearch indicies created
    - kubernetes_cluster-- (for container logs)
    - kubernetes_audit-- (for audit logs from kube-apiserver)
    - kubernetes_host-- (for all systemd host logs)
  - version 1.5.2
    - Kernel messages forwarded ([#375](https://github.com/mesosphere/kubernetes-base-addons/pull/375), [@sebbrandt87](https://github.com/sebbrandt87))
  - apply meaningful aliases to plugins and their metrics. ([#432](https://github.com/mesosphere/kubernetes-base-addons/pull/432), [@branden](https://github.com/branden))
- Istio:
  - the "kubernetes-service-monitor" service monitor has been removed. ([#483](https://github.com/mesosphere/kubernetes-base-addons/pull/483), [@gracedo](https://github.com/gracedo))
- Traefik-foward-auth:
  - update to 0.2.14
    - Add an option to bypass tfa deployment ([#456](https://github.com/mesosphere/kubernetes-base-addons/pull/456), [@d2iq-dispatch](https://github.com/d2iq-dispatch))
- Kibana:
    - version 6.8.10 ([#373](https://github.com/mesosphere/kubernetes-base-addons/pull/373), [@sebbrandt87](https://github.com/sebbrandt87))
- Ops-portal:
  - Fix: Unable to change ops-portal password ([#379](https://github.com/mesosphere/kubernetes-base-addons/pull/379), [@GoelDeepak](https://github.com/GoelDeepak))
- Prometheus:
  - chore: bump chart to v9.3.1
    - refactor!: (breaking change) version 9 of the helm chart removes the existing `additionalScrapeConfigsExternal` in favor of `additionalScrapeConfigsSecret`. This change lets users specify the secret name and secret key to use for the additional scrape configuration of prometheus.
    - feat: add ingress configuration for Thanos sidecar, enabling external access from a centralized thanos querier running in another cluster
    - feat: add scrape timeout config to service monitor to avoid timeouts on slow kubelets
    - feat: add docker checksum option to improve security for deployed containers
    - feat: add option to disable availability rules
    - feat: enable scraping /metrics/resource for kubelet service
    - feat: [prometheus] enable namespace overrides
    - feat: [prometheus] allow additional volumes and volumeMounts
    - feat: [alertmanager] add volume and volume mounts to spec
    - feat: [alertmanager] add support for serviceAccount.annotations
    - feat: [grafana] enable adding annotations to all default dashboard configmaps
    - chore: bump prometheus to v2.18.2
    - chore: bump alertmanager to v0.21.0
    - chore: bump hyperkube to v1.16.12
    - chore: bump grafana to v5.3.0
    - fix: add missing grafana annotations to k8s-coredns dashboard
    - fix: reduced CPU utilization and time lag for code_verb:apiserver_request_total:increase30d scrape
    - fix: invalid image pull policy for the admission webhook patch
    - fix: alert "KubeNodeUnreachable" no longer fires on an autoscaling scale-down event ([#444](https://github.com/mesosphere/kubernetes-base-addons/pull/444), [@samvantran](https://github.com/samvantran))
  - disable ServiceMonitors for kube-controller-manager and kube-scheduler. kubernetes has determined the ports that were used for these tests was insecure and has limited it to localhost only. This causes these specific tests to fail. The state of the controller-manager and scheduler pods are still tracked in general as pods. ([#474](https://github.com/mesosphere/kubernetes-base-addons/pull/474), [@dkoshkin](https://github.com/dkoshkin))

## stable-1.15-2.2.0, stable-1.16-2.2.0, stable-1.17-2.2.0

* Prometheus
  * Fix an issue that may cause Grafana's home dashboard to be empty. ([#351](https://github.com/mesosphere/kubernetes-base-addons/pull/351), [@branden](https://github.com/branden))
  * disable ServiceMonitors for kube-controller-manager and kube-scheduler. kubernetes has determined the ports that were used for these tests was insecure and has limited it to localhost only. This causes these specific tests to fail. The state of the controller-manager and scheduler pods are still tracked in general as pods. ([#474](https://github.com/mesosphere/kubernetes-base-addons/pull/474), [@dkoshkin](https://github.com/dkoshkin))
  * Improve Grafana dashboard names and tags for dashboards tied to addons ([#352](https://github.com/mesosphere/kubernetes-base-addons/pull/352), [@gracedo](https://github.com/gracedo))
* Traefik
  * fix metrics access and reporting ([#349](https://github.com/mesosphere/kubernetes-base-addons/pull/349), [@gracedo](https://github.com/gracedo))

## stable-1.15-2.1.1, stable-1.16-2.1.1, stable-1.17-2.1.1

* dex-k8s-authenticator
  * Windows download support for the credentials plugin ([#377](https://github.com/mesosphere/kubernetes-base-addons/pull/377), [@jr0d](https://github.com/jr0d))

## stable-1.15-2.1.0, stable-1.16-2.1.0, stable-1.17-2.1.0

* traefik
  * fix the velero-minio entrypoint to inherit global ssl and proxy protocol configurations ([#259](https://github.com/mesosphere/kubernetes-base-addons/pull/259), [@jieyu](https://github.com/jieyu))
* elasticsearch
  * default data nodes has been increased to 4 ([#327](https://github.com/mesosphere/kubernetes-base-addons/pull/327), [@alejandroEsc](https://github.com/alejandroEsc))
* external-dns
  * disable by default ([#335](https://github.com/mesosphere/kubernetes-base-addons/pull/335), [@GoelDeepak](https://github.com/GoelDeepak))

## stable-1.15-2.0.1, stable-1.16-2.0.1, stable-1.17-2.0.1

- Traefik: fix metrics access and reporting ([#349](https://github.com/mesosphere/kubernetes-base-addons/pull/349), [@gracedo](https://github.com/gracedo))
- Prometheus: Improve Grafana dashboard names and tags for dashboards tied to addons ([#352](https://github.com/mesosphere/kubernetes-base-addons/pull/352), [@gracedo](https://github.com/gracedo))

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
