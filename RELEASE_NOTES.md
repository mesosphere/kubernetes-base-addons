# Release Notes
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
