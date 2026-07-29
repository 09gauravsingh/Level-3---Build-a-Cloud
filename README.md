# build-a-cloud

Hands-on infrastructure work for the LEVEL3 programme. Each week is a
self-contained directory that builds one layer of a cloud platform.

| Directory | Topic |
| --- | --- |
| `week2-ske-paas` | A managed PostgreSQL Platform-as-a-Service on STACKIT Kubernetes Engine |

## Week 2 — Managed PostgreSQL on SKE

The goal of this week is to deliver a database *product* rather than a database
*server*. A platform user asks for a PostgreSQL cluster by submitting a single
Custom Resource, and the platform takes care of provisioning, failover,
storage and credentials.

The stack has three layers:

1. **Infrastructure** — a STACKIT Kubernetes Engine cluster provisioned with
   Terraform.
2. **Platform** — the CloudNativePG operator, installed with Helm, which teaches
   Kubernetes what a PostgreSQL cluster is.
3. **Product** — a `Cluster` Custom Resource that the operator reconciles into a
   running, self-healing database service.

### Infrastructure

Terraform in `week2-ske-paas/01-ske-infrastructure` creates the SKE cluster:

- Cluster `lvl3-paas` in region `eu01`
- Node pool `paas-workers`, machine type `g3i.2`, 2 workers, 32 GB premium disks
- Availability zone `eu01-1`
- Public control-plane access scope
- Automatic Kubernetes and machine-image updates during a 01:00–03:00 UTC
  maintenance window

The administrator kubeconfig is issued by the `stackit_ske_kubeconfig` resource
with a seven-day lifetime and written to disk with `0600` permissions. Provider
authentication uses a STACKIT service-account key held outside this repository
in `~/.config/stackit`.

### Platform

`week2-ske-paas/02-operator` records the CloudNativePG installation: Helm chart
version 0.29.0, released as `cnpg` into the `cnpg-system` namespace, running
operator image 1.30.0. Installing the chart registers the
`postgresql.cnpg.io` API group, which is what makes the product layer possible.

### Product

`week2-ske-paas/03-product` holds the service definition:

- `namespace.yaml` — the `database-services` namespace all databases live in
- `level3-postgres.yaml` — the `level3-postgres` cluster: two PostgreSQL 18.4
  instances, 5 GiB per instance, application database `platformdb` owned by
  `platformuser`, data checksums enabled, and required pod anti-affinity so the
  primary and standby never share a worker node
- `psql-client.yaml` — a throwaway client pod that reads its credentials from
  the operator-generated `level3-postgres-app` Secret and connects with
  `sslmode=require`

CloudNativePG publishes three ClusterIP Services, so applications choose their
own read/write semantics through DNS alone:

| Service | Purpose |
| --- | --- |
| `level3-postgres-rw` | Always routes to the current primary |
| `level3-postgres-ro` | Routes to standby replicas only |
| `level3-postgres-r` | Routes to any instance |

### Connectivity and verification

`week2-ske-paas/04-connectivity/test-postgres.sql` connects through the
read/write Service, reports which instance answered, creates a `paas_demo`
table and reads its rows back.

`week2-ske-paas/evidence` captures a deliberate failover. Before the test the
primary was `level3-postgres-1`; after deleting that pod the operator promoted
`level3-postgres-2` and reported `Cluster in healthy state` again, with the
`level3-postgres-rw` Service following the new primary automatically. That is
the property that makes this a platform rather than a deployment.

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

## What is deliberately not committed

Terraform state, plan files, generated kubeconfigs and service-account keys are
excluded by `.gitignore`. Terraform state stores resource attributes in
cleartext, including the full cluster-admin kubeconfig, so it must never reach
a remote repository. Anyone reproducing this work generates their own state and
credentials from the committed configuration.

## Current limitations

This is a learning prototype. Both workers sit in the same availability zone,
so it survives node loss but not zone loss. Backups, monitoring, tenant
isolation, billing and production security policy are out of scope for this
stage.
