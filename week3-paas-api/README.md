# Week 3 – PaaS API Layer

A Go REST API for provisioning and managing CloudNativePG PostgreSQL instances in STACKIT SKE.

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