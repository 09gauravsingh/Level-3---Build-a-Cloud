# Week 3 – PaaS API Layer

A Go REST API for provisioning and managing CloudNativePG PostgreSQL
instances in STACKIT SKE. It is the control plane for the product: people
register or sign in, receive a JWT, and then create, list, connect to and
delete PostgreSQL clusters. CloudNativePG still does the actual
provisioning.

How those calls move through the rest of the platform is in
[../architecture.md](../architecture.md).

## What the API does

- `POST /api/v1/register` — create a SQLite account (bcrypt password).
  Usernames must be Kubernetes-label safe (3–32 lowercase letters, numbers
  or hyphens). The environment admin name is reserved.
- `POST /api/v1/login` — environment admin or a registered user. Returns a
  one-hour HS256 JWT with `sub` and `role` (`admin` or `user`).
- `POST/GET/DELETE /api/v1/instances` and
  `GET /api/v1/instances/:name/connection` — require `Authorization: Bearer`.
  Ordinary users only see clusters labelled `platform.level3.io/owner` with
  their username. The admin sees every API-managed cluster. A name that
  belongs to someone else is returned as `404`, not `403`.
- `GET /healthz` and `GET /metrics` — public probes.

Instance names are unique in `database-services`. A second create of the
same name returns `409 Conflict`.

## Environment

| Variable | Required | Default | Purpose |
| --- | --- | --- | --- |
| `ADMIN_USERNAME` | yes | — | Superuser login |
| `ADMIN_PASSWORD` | yes | — | Superuser password |
| `JWT_SECRET` | yes | — | HS256 signing key |
| `PAAS_NAMESPACE` | no | `database-services` | Where Clusters are created |
| `PORT` | no | `8080` | Listen address |
| `USER_DB_PATH` | no | `./users.db` | SQLite file (`/data/users.db` in the cluster) |
| `KUBECONFIG` | local only | default kubeconfig | Cluster access outside SKE |

Permanent local values live in `.env.secrets` and are never committed.
In the cluster the SQLite file sits on a PVC because the container root
filesystem is read-only. The Deployment uses `replicas: 1` and
`strategy: Recreate` so two Pods never share that `ReadWriteOnce` volume.

## Commands

Run `make` to list them. Secrets come from `.env.secrets` (never committed)
or from the command line, for example `make run PORT=9090`.

| Command | Purpose |
| --- | --- |
| `make run` | Start the server locally |
| `make stop` | Free port 8080 |
| `make test` | Run unit tests |
| `make vet` | Run `go vet ./...` |
| `make kube-refresh` | Write a fresh login kubeconfig with the STACKIT CLI |
| `make deploy` | Apply the manifests and wait for the Pod |
| `make status` | Deployment, Pods and Service |
| `make logs` | Follow the Pod logs |
| `make forward` | Tunnel `localhost:8080` to the Service |
| `make connect-prometheus` | Tunnel `localhost:9090` to Prometheus |
| `make connect-grafana` | Tunnel `localhost:3000` to Grafana |
| `make publish` | Build, push and roll out `API_IMAGE` |

Publish a new image (always override the tag):

```bash
make publish API_IMAGE=registry.onstackit.cloud/gaurav-paas-20a60c06/week3-paas-api:v1.4.3
```

`make forward` and `make run` share port 8080; `make stop` frees it.

```bash
make forward
TOKEN=$(curl -sS -X POST localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"week4-password"}' | jq -r .token)
curl -sS localhost:8080/api/v1/instances -H "Authorization: Bearer $TOKEN"
```

## Cluster access

Local commands talk to the SKE cluster through the login kubeconfig
`~/.kube/lvl3-paas-login.yaml`, created with:

```bash
stackit ske kubeconfig create lvl3-paas --login --filepath ~/.kube/lvl3-paas-login.yaml
```

That file holds no credentials. Every call asks the STACKIT CLI for a fresh
token, so the kubeconfig does not expire, but you must stay logged in with
`stackit auth login`. `make kube-refresh` writes the file again.

Inside SKE the API uses no kubeconfig at all: it authenticates with the
ServiceAccount token of its own Pod, using the permissions in
`deploy/01-rbac.yaml`.

## Files

```
cmd/api/main.go
  reads environment variables
  opens the Kubernetes client and the SQLite user store
  calls apihttp.NewRouter(service, logger, users)
  starts and stops the HTTP server

internal/api/routes.go
  Gin router, CORS, metrics middleware
  public routes: /healthz, /metrics, /api/v1/login, /api/v1/register
  protected /api/v1/instances* group

internal/api/auth.go
  register handler, username/password rules, ownerScope helper

internal/api/users.go
  SQLite UserStore: CreateUser (bcrypt) and Authenticate

internal/api/handlers.go
  login (admin env, then SQLite) plus Create, List, Delete, Connection, Health

internal/api/middleware.go
  JWT verification, username/isAdmin on the gin context
  Kubernetes errors → HTTP status

internal/api/metrics.go
  total_http_requests and http_request_duration_seconds

internal/api/handlers_test.go
  route and auth tests against a fake PlatformService and a temp SQLite file

internal/models/
  request and response structs, including ownedBy on Instance

internal/platform/kubernetes_client.go
  dynamic + core clients; kubeconfig locally, ServiceAccount in SKE
  label constants managed-by and owner

internal/platform/kubernetes_instances.go
  create / list / delete Cluster resources
  list and getManagedCluster filter by owner when the scope is non-empty
  connection credentials from the <name>-app Secret

api/openapi.yaml
  official OpenAPI contract

deploy/01-rbac.yaml
  paas-system namespace, ServiceAccount, Role in database-services

deploy/02-api.yaml
  PVC for SQLite, Deployment, Service
```