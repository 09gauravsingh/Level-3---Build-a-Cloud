# build-a-cloud

Hands-on infrastructure work for the LEVEL3 programme. Each week is a
directory that adds one layer of a PostgreSQL Platform-as-a-Service on
STACKIT Kubernetes Engine. Together they are one product: a public
dashboard where people register, sign in, and provision their own
managed PostgreSQL clusters.

The flows, namespaces, auth model and observability path are drawn in
[architecture.md](architecture.md).

| Directory | Layer |
| --- | --- |
| `week2-ske-paas` | SKE cluster, CloudNativePG operator, first product `Cluster` |
| `week3-paas-api` | Go REST API: accounts, JWT, instance lifecycle |
| `week4-web-ui` | Vue dashboard over HTTPS |
| `week5-observability` | Prometheus, Grafana, Loki, Alloy |

Public hostname: `https://gaurav-paas.runs.onstackit.cloud`.

---

## What you can do with it

1. Open the dashboard and create an account, or sign in as the environment
   administrator.
2. Ask for a PostgreSQL cluster (name, replica count, storage, database).
3. Watch CloudNativePG provision it. The list refreshes every ten seconds.
4. Copy connection credentials from the generated Secret.
5. Delete the cluster when you are done.

A registered user only sees clusters they created. The administrator sees
every cluster the API manages. Isolation is a Kubernetes label
(`platform.level3.io/owner`) on the `Cluster` object, not a namespace per
user. Cluster names are unique in `database-services`; a taken name returns
`409 Conflict`.

---

## Week 2 — Managed PostgreSQL on SKE

The goal of this week is to deliver a database *product* rather than a
database *server*. A platform user asks for a PostgreSQL cluster by
submitting a single Custom Resource, and the platform takes care of
provisioning, failover, storage and credentials.

The stack has three layers:

1. **Infrastructure** — a STACKIT Kubernetes Engine cluster provisioned
   with Terraform.
2. **Platform** — the CloudNativePG operator, installed with Helm, which
   teaches Kubernetes what a PostgreSQL cluster is.
3. **Product** — a `Cluster` Custom Resource that the operator reconciles
   into a running, self-healing database service.

### Infrastructure

Terraform in `week2-ske-paas/01-ske-infrastructure` creates the SKE cluster:

- Cluster `lvl3-paas` in region `eu01`
- Node pool `paas-workers`, machine type `g3i.2`, 2 workers, 32 GB premium disks
- Availability zone `eu01-1`
- Public control-plane access scope
- Automatic Kubernetes and machine-image updates during a 01:00–03:00 UTC
  maintenance window

The administrator kubeconfig is issued by the `stackit_ske_kubeconfig`
resource with a seven-day lifetime and written to disk with `0600`
permissions. Provider authentication uses a STACKIT service-account key
held outside this repository in `~/.config/stackit`.

### Platform

`week2-ske-paas/02-operator` records the CloudNativePG installation: Helm
chart version 0.29.0, released as `cnpg` into the `cnpg-system` namespace,
running operator image 1.30.0. Installing the chart registers the
`postgresql.cnpg.io` API group, which is what makes the product layer
possible.

### Product

`week2-ske-paas/03-product` holds the service definition:

- `namespace.yaml` — the `database-services` namespace all databases live in
- `level3-postgres.yaml` — the `level3-postgres` cluster: two PostgreSQL 18.4
  instances, 5 GiB per instance, application database `platformdb` owned by
  `platformuser`, data checksums enabled, and required pod anti-affinity so
  the primary and standby never share a worker node
- `psql-client.yaml` — a throwaway client pod that reads its credentials
  from the operator-generated `level3-postgres-app` Secret and connects with
  `sslmode=require`

CloudNativePG publishes three ClusterIP Services, so applications choose
their own read/write semantics through DNS alone:

| Service | Purpose |
| --- | --- |
| `level3-postgres-rw` | Always routes to the current primary |
| `level3-postgres-ro` | Routes to standby replicas only |
| `level3-postgres-r` | Routes to any instance |

### Connectivity and verification

`week2-ske-paas/04-connectivity/test-postgres.sql` connects through the
read/write Service, reports which instance answered, creates a `paas_demo`
table and reads its rows back.

`week2-ske-paas/evidence` captures a deliberate failover. Before the test
the primary was `level3-postgres-1`; after deleting that pod the operator
promoted `level3-postgres-2` and reported `Cluster in healthy state` again,
with the `level3-postgres-rw` Service following the new primary
automatically. That is the property that makes this a platform rather than
a deployment.

## Reproducing the week 2 build

```bash
cd week2-ske-paas/01-ske-infrastructure
terraform init
terraform apply
export KUBECONFIG="$PWD/kubeconfig.yaml"

helm repo add cnpg https://cloudnative-pg.github.io/charts
helm upgrade --install cnpg cnpg/cloudnative-pg \
  --namespace cnpg-system --create-namespace --version 0.29.0

cd ../03-product
kubectl apply -f namespace.yaml
kubectl apply -f level3-postgres.yaml
kubectl apply -f psql-client.yaml

kubectl -n database-services get cluster level3-postgres
```

---

## Week 3 — REST control plane

`week3-paas-api` is a Go service that turns HTTP into CloudNativePG
`Cluster` objects. It does not run PostgreSQL. The operator does.

What it adds on top of week 2:

- Self-service **register** and **login**. Registered users live in SQLite
  (`USER_DB_PATH`). The environment admin (`ADMIN_USERNAME` /
  `ADMIN_PASSWORD`) remains a superuser.
- JWT sessions (HS256, one hour) with `sub` and `role`.
- Per-user dashboards via the `platform.level3.io/owner` label.
- Prometheus metrics on `/metrics` and a public `/healthz`.
- OpenAPI contract in `week3-paas-api/api/openapi.yaml`.

Local loop (see `week3-paas-api/README.md` for the full Makefile):

```bash
cd week3-paas-api
make run          # API on :8080, talks to SKE through kubeconfig
make unit-test    # routes and auth, no cluster required
```

In the cluster the API runs in `paas-system`, writes SQLite to a PVC at
`/data/users.db`, and uses a ServiceAccount that may only manage clusters
and read Secrets in `database-services`.

---

## Week 4 — Web console

`week4-web-ui` is a Vue 3 + Vite dashboard. The bundle only ever calls
same-origin `/api/...` paths:

- **Development:** Vite proxies `/api` to `http://localhost:8080`.
- **Production:** Traefik Ingress sends `/api` to the API Service and `/`
  to the UI Service, with a Let's Encrypt certificate on
  `gaurav-paas.runs.onstackit.cloud`.

The login screen toggles between sign-in and create-account. After
register the UI signs the user in and shows only their instances. The
administrator sees every instance and an admin badge.

```bash
cd week4-web-ui
npm install
npm run dev       # http://localhost:5173
```

---

## Week 5 — Observability

`week5-observability` installs the viewing side of the platform:

- **Prometheus + Grafana** (`monitoring/values.yaml`) scrape the API
  ServiceMonitor (`monitoring/api-servicemonitor.yaml`) for
  `total_http_requests` and `http_request_duration_seconds`.
- **Loki + Alloy** (`logging/`) collect pod logs through the Kubernetes
  API and store them for Grafana.

The API already emits JSON logs to stdout and Prometheus text on
`/metrics`; this week is the in-cluster pipeline that scrapes and stores
them.

---

## What is deliberately not committed

Terraform state, plan files, generated kubeconfigs, `.env.secrets` and
service-account keys are excluded by `.gitignore`. Terraform state stores
resource attributes in cleartext, including the full cluster-admin
kubeconfig, so it must never reach a remote repository. The local SQLite
file `users.db` is ignored for the same reason. Anyone reproducing this
work generates their own state and credentials from the committed
configuration.

## Current limitations

This is a learning prototype. Both workers sit in the same availability
zone, so it survives node loss but not zone loss. Tenant isolation is a
label on a shared namespace, not a namespace per user. SQLite is a single
replica. There is no backup of user databases, no password reset, and no
refresh token. Billing and production security policy remain out of scope.
