# Week 3 – PaaS API Layer

A Go REST API for provisioning and managing CloudNativePG PostgreSQL instances in STACKIT SKE.

## Local commands

Every long command is wrapped in the `Makefile`. Run `make` to list them all.

| Command | Purpose |
| --- | --- |
| `make run` | Start the server with KUBECONFIG, API_TOKEN, PAAS_NAMESPACE and PORT already set |
| `make stop` | Kill whatever still listens on port 8080 |
| `make restart` | Free the port and start the server again |
| `make check` | Tidy, format, vet and test before a commit |
| `make unit-test` | Run only the REST API unit tests, with full output |
| `make health` | Call `/healthz` |
| `make create` | Create an instance, for example `make create NAME=my-db` |
| `make list` | List all instances |
| `make connection` | Show the credentials of one instance |
| `make delete` | Delete one instance |
| `make clusters` | Show the CloudNativePG clusters with kubectl |
| `make kube-refresh` | Write a fresh login kubeconfig with the STACKIT CLI |
| `make kube-context` | Show which cluster and user kubectl would use |
| `make docs` | Validate the contract and open Swagger Editor |
| `make docker-build` | Build the container image |

Values may be overridden per command, for example `make run PORT=9090`.
Permanent values, including the API token, live in `.env.secrets`.
That file stays on your machine and is never committed.

`eval "$(make load)"` puts `KUBECONFIG`, `API_TOKEN`, `PAAS_NAMESPACE` and
`PORT` into the current shell, so plain `kubectl` works afterwards. A make
recipe runs in a child shell and cannot export into its caller, which is why
the exports are evaluated rather than set.

## Cluster commands

These act on the deployment running in SKE, in namespace `paas-system`.

| Command | Purpose |
| --- | --- |
| `make preflight` | Check every prerequisite in one go: session, namespaces, CloudNativePG, both Secrets, Deployment, endpoint, RBAC and manifest drift |
| `make deploy` | Apply `deploy/01-rbac.yaml` and `deploy/02-api.yaml`, then wait for the new Pod |
| `make deploy-diff` | Show what `make deploy` would change, without changing it |
| `make rollout` | Restart the Pods and wait for them |
| `make status` | Deployment, Pods and Service in one view |
| `make image` | The image the running Pod was started from |
| `make logs` | Follow the Pod logs |
| `make describe` | Describe the Pod, including its events |
| `make events` | Recent events of the namespace |
| `make forward` | Tunnel `localhost:8080` to the Service |
| `make cluster-token` | Print the token from the `week3-paas-api-auth` Secret |

### Talking to the Pod

The Pod reads its bearer token from a Kubernetes Secret, so it differs from
the local one in `.env.secrets`. Add `REMOTE=1` and the curl helpers use the
cluster token instead:

```bash
make forward             # first terminal, keeps the tunnel open
make list REMOTE=1       # second terminal
make create REMOTE=1 NAME=demo-db
```

Without `REMOTE=1` the same commands reach a locally running `make run`.
The tunnel and the local server share port 8080, so only one may run at a
time; `make stop` frees the port for either.

To use plain `curl` instead of the helpers, put the token of the Pod into
the shell once:

```bash
eval "$(make load-remote)"
curl -sS localhost:8080/api/v1/instances -H "Authorization: Bearer $API_TOKEN"
```

`make load-remote` is `make load` with the token taken from the Secret
rather than from `.env.secrets`; `make cluster-token` prints it on its own.

## Cluster access

Local commands talk to the SKE cluster through the login kubeconfig
`~/.kube/lvl3-paas-login.yaml`, created with:

```bash
stackit ske kubeconfig create lvl3-paas --login --filepath ~/.kube/lvl3-paas-login.yaml
```

That file holds no credentials. Every call asks the STACKIT CLI for a fresh
token, so the kubeconfig does not expire, but you must stay logged in with
`stackit auth login`. `make kube-refresh` writes the file again.

The container is the one exception. The STACKIT CLI does not exist inside the
image, so `make docker-run` first writes a separate one-hour kubeconfig with
embedded credentials to `bin/docker-kubeconfig.yaml` and mounts that instead.
`make clean` removes it.

Inside SKE the API uses no kubeconfig at all: it authenticates with the
ServiceAccount token of its own Pod, using the permissions in
`deploy/01-rbac.yaml`.

## Planned capabilities

- Create PostgreSQL instances
- List PostgreSQL instances
- Retrieve instance status
- Delete PostgreSQL instances
- Retrieve connection data
- OpenAPI documentation
- Unit tests
- Docker image
- Deployment to STACKIT SKE


Responsibility of each files:

cmd/api/main.go
→ reads environment variables
→ creates the Kubernetes service
→ calls apihttp.NewRouter(...)
→ starts and stops the HTTP server


internal/api/routes.go
→ creates the Gin router
→ adds middleware
→ registers REST routes


internal/api/handlers.go
→ implements Create, List, Delete, Connection and Health handlers


internal/api/handlers_test.go
→ unit tests for the routes, run with make unit-test
→ replaces Kubernetes with a fake service, so no cluster is needed


internal/api/middleware.go
→ bearer-token authentication
→ conversion of platform errors into HTTP responses


internal/models/
→ request and response structs


internal/platform/kubernetes_client.go
→ creates Kubernetes clients
→ uses kubeconfig locally or ServiceAccount inside SKE


internal/platform/kubernetes_instances.go
→ creates, lists and deletes CloudNativePG Cluster resources
→ reads connection credentials from Kubernetes Secrets


api/openapi.yaml
→ official OpenAPI contract displayed through Swagger