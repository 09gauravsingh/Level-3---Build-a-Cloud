# Short commands for the Week 3 PaaS API.
#
# Run "make" or "make help" to list everything.
# Every value below can be overridden on the command line,
# for example: make run PORT=9090
# Persistent overrides belong in the local .env.secrets file.

SHELL := /bin/bash

# Local settings and keys. This file is optional and is never committed.
-include .env.secrets

# --- Configuration -----------------------------------------------------------

# The SKE login kubeconfig written by the STACKIT CLI.
# It carries no credentials itself: kubectl and the API server ask the
# STACKIT CLI for a fresh token on every call, so it never expires.
KUBECONFIG ?= $(HOME)/.kube/lvl3-paas-login.yaml  
# need to remove brackets from Home and (?) to run the command separately.(remove spaces between spaces in the path)

# The SKE cluster that the login kubeconfig points to, and the STACKIT
# project that owns it. Both match the Week 2 Terraform variables.
CLUSTER    ?= lvl3-paas
PROJECT_ID ?= 20a60c06-1e1a-406f-a840-37d1ff14f0e8

# Legacy env var still required by main.go at startup (unused by JWT middleware).
API_TOKEN ?= local-week3-development-token

# Login credentials for POST /api/v1/login (Web UI and curl).
ADMIN_USERNAME ?= admin
ADMIN_PASSWORD ?= week4-password

# Secret used to sign and verify JWTs. Keep stable across local restarts.
JWT_SECRET ?= my-week4-super-secret-key

# Namespace that holds the CloudNativePG clusters.
PAAS_NAMESPACE ?= database-services

PORT     ?= 8080
GIN_MODE ?= debug

# --- The deployment running inside SKE ---------------------------------------

# Name shared by the Deployment, Service, ServiceAccount and app label.
APP ?= week3-paas-api

# Namespace of the control plane itself, not of the databases it manages.
API_NAMESPACE ?= paas-system

DEPLOY_DIR      ?= deploy
ROLLOUT_TIMEOUT ?= 300s

# The API server reads these from the environment.
export KUBECONFIG API_TOKEN PAAS_NAMESPACE PORT GIN_MODE
export ADMIN_USERNAME ADMIN_PASSWORD JWT_SECRET

# --- Values used by the curl helpers ----------------------------------------

BASE_URL  ?= http://localhost:$(PORT)
NAME      ?= week3-demo-db
INSTANCES ?= 1
STORAGE   ?= 1Gi
DATABASE  ?= demoapp
OWNER     ?= demoowner

# --- Build and tooling paths -------------------------------------------------

BINARY ?= bin/api
IMAGE  ?= week3-paas-api:local

# The login kubeconfig calls the STACKIT CLI, which does not exist inside
# the image, so the container needs its own kubeconfig with embedded
# credentials. It lands in bin/ and is removed by "make clean".
DOCKER_KUBECONFIG     ?= bin/docker-kubeconfig.yaml
DOCKER_KUBECONFIG_TTL ?= 1h

# gofmt lives in the Go toolchain directory and not always on PATH.
GOFMT := $(shell go env GOROOT)/bin/gofmt

# Pretty-print JSON responses when jq is installed, otherwise print raw.
JQ := $(shell command -v jq 2>/dev/null || echo cat)

# The Pod authenticates with the token in the week3-paas-api-auth Secret,
# which differs from the local one in .env.secrets. Reading it needs a
# cluster call, so this stays lazy and runs only when REMOTE=1 is used.
CLUSTER_TOKEN = $(shell kubectl get secret $(APP)-auth \
	--namespace $(API_NAMESPACE) \
	--output jsonpath='{.data.API_TOKEN}' 2>/dev/null | base64 --decode)

# Point the curl helpers at the Pod instead of a local server:
#   make forward        (in one terminal)
#   make list REMOTE=1  (in another)
ifeq ($(REMOTE),1)
API_TOKEN = $(CLUSTER_TOKEN)
endif

# Recursive on purpose: API_TOKEN may resolve to a cluster lookup.
AUTH = Authorization: Bearer $(API_TOKEN)

.DEFAULT_GOAL := help

.PHONY: help env load load-remote token run restart start build test test-verbose unit-test cover \
	fmt vet tidy check ci clean stop kubeconfig-check \
	health create list connection delete smoke \
	nodes clusters pods kube-refresh kube-context \
	openapi docs docker-build docker-kubeconfig docker-run docker-shell \
	deploy deploy-diff deploy-delete rollout rollout-status \
	status logs describe events image forward cluster-token preflight

help: ## Show every available command
	@awk 'BEGIN { FS = ":.*##"; printf "\nUsage: make <target> [VAR=value]\n" } \
		/^##@/ { printf "\n%s\n", substr($$0, 5); next } \
		/^[a-zA-Z0-9_-]+:.*##/ { printf "  %-18s%s\n", $$1, $$2 } \
		END { printf "\n" }' $(firstword $(MAKEFILE_LIST))

##@ Server

run: kubeconfig-check ## Start the API server with all keys set
	go run ./cmd/api

restart: stop run ## Free the port and start the server again

build: ## Compile the server into bin/api
	go build -o $(BINARY) ./cmd/api

start: build kubeconfig-check ## Run the compiled binary
	./$(BINARY)

# A hung server ignores SIGTERM, so the port is checked again
# afterwards and SIGKILL is used when the process is still there.
stop: ## Kill whatever is listening on the API port
	@pid="$$(lsof -ti tcp:$(PORT) -sTCP:LISTEN 2>/dev/null)"; \
	if [ -z "$$pid" ]; then \
		echo "Port $(PORT) is already free."; \
		exit 0; \
	fi; \
	echo "Stopping process $$pid on port $(PORT)."; \
	kill $$pid 2>/dev/null; \
	for _ in 1 2 3 4 5; do \
		sleep 1; \
		lsof -ti tcp:$(PORT) -sTCP:LISTEN >/dev/null 2>&1 || break; \
	done; \
	pid="$$(lsof -ti tcp:$(PORT) -sTCP:LISTEN 2>/dev/null)"; \
	if [ -n "$$pid" ]; then \
		echo "Process $$pid ignored SIGTERM, sending SIGKILL."; \
		kill -9 $$pid 2>/dev/null; \
		sleep 1; \
	fi; \
	echo "Port $(PORT) is free."

##@ Quality

test: ## Run the unit tests
	go test ./...

test-verbose: ## Run the unit tests with full output
	go test -v ./...

unit-test: ## Run the unit tests of the REST API on their own
	go test ./internal/api -v

cover: ## Run the tests and show total coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

fmt: ## Format all Go files
	go fmt ./...

vet: ## Run the Go static checks
	go vet ./...

tidy: ## Synchronise go.mod and go.sum
	go mod tidy

check: tidy fmt vet test ## Format, vet and test before committing

ci: ## Verify formatting, vet, tests and build without changing files
	@unformatted="$$("$(GOFMT)" -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files need gofmt:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	go vet ./...
	go test ./...
	go build -o $(BINARY) ./cmd/api

clean: ## Remove build output, coverage files and the container kubeconfig
	rm -rf bin coverage.out

##@ API requests

health: ## GET /healthz
	@curl -sS $(BASE_URL)/healthz | $(JQ)

create: ## Create an instance (override with NAME=my-db)
	@curl -sS -X POST $(BASE_URL)/api/v1/instances \
		-H "$(AUTH)" \
		-H "Content-Type: application/json" \
		-d '{"name":"$(NAME)","instances":$(INSTANCES),"storageSize":"$(STORAGE)","database":"$(DATABASE)","owner":"$(OWNER)"}' \
		| $(JQ)

list: ## List every instance created by this API
	@curl -sS $(BASE_URL)/api/v1/instances -H "$(AUTH)" | $(JQ)

connection: ## Show the connection data of one instance
	@curl -sS $(BASE_URL)/api/v1/instances/$(NAME)/connection -H "$(AUTH)" | $(JQ)

delete: ## Delete one instance
	@curl -sS -X DELETE $(BASE_URL)/api/v1/instances/$(NAME) -H "$(AUTH)" | $(JQ)

smoke: health create list ## Health check, create and list in one command

##@ Kubernetes

kubeconfig-check: ## Confirm the kubeconfig file exists
	@test -f "$(KUBECONFIG)" || { \
		echo "Kubeconfig was not found: $(KUBECONFIG)"; \
		echo "Create it with: make kube-refresh"; \
		exit 1; \
	}

nodes: kubeconfig-check ## Show the SKE cluster nodes
	kubectl get nodes

clusters: kubeconfig-check ## Show the CloudNativePG clusters
	kubectl get clusters -n $(PAAS_NAMESPACE)

pods: kubeconfig-check ## Show the pods of the product namespace
	kubectl get pods -n $(PAAS_NAMESPACE)

# The CLI always asks before touching a kubeconfig, so --assume-yes is
# passed. It updates the entry in the file instead of replacing the file.
kube-refresh: ## Write a fresh login kubeconfig with the STACKIT CLI
	stackit ske kubeconfig create $(CLUSTER) --login --assume-yes \
		--project-id $(PROJECT_ID) \
		--filepath "$(KUBECONFIG)"

kube-context: kubeconfig-check ## Show which cluster and user kubectl would use
	kubectl config get-contexts

##@ Deployment on SKE

deploy: kubeconfig-check ## Apply the manifests and wait for the new Pod
	kubectl apply -f $(DEPLOY_DIR)/01-rbac.yaml
	kubectl apply -f $(DEPLOY_DIR)/02-api.yaml
	@$(MAKE) --no-print-directory rollout-status

# kubectl diff reports differences with exit code 1, which is not a
# failure here, so the result is turned into a readable sentence.
deploy-diff: kubeconfig-check ## Show what "make deploy" would change
	@if kubectl diff -f $(DEPLOY_DIR)/02-api.yaml; then \
		echo "The cluster already matches $(DEPLOY_DIR)/02-api.yaml."; \
	fi

deploy-delete: kubeconfig-check ## Remove the API Deployment and Service
	kubectl delete -f $(DEPLOY_DIR)/02-api.yaml --ignore-not-found

rollout: kubeconfig-check ## Restart the Pods and wait for them
	kubectl rollout restart deployment/$(APP) --namespace $(API_NAMESPACE)
	@$(MAKE) --no-print-directory rollout-status

rollout-status: kubeconfig-check ## Wait until the Deployment is available
	kubectl rollout status deployment/$(APP) \
		--namespace $(API_NAMESPACE) \
		--timeout=$(ROLLOUT_TIMEOUT)

status: kubeconfig-check ## Show the Deployment, Pods and Service in SKE
	@kubectl get deployment,pods,service \
		--namespace $(API_NAMESPACE) \
		--output wide

image: kubeconfig-check ## Show the image the running Pod was started from
	@kubectl get deployment/$(APP) --namespace $(API_NAMESPACE) \
		--output jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'

logs: kubeconfig-check ## Follow the Pod logs, stop with Ctrl-C
	kubectl logs deployment/$(APP) \
		--namespace $(API_NAMESPACE) \
		--follow --tail=50

describe: kubeconfig-check ## Describe the API Pod, including its events
	@kubectl describe pod \
		--namespace $(API_NAMESPACE) \
		--selector app=$(APP)

events: kubeconfig-check ## Show recent events of the API namespace
	@kubectl get events \
		--namespace $(API_NAMESPACE) \
		--sort-by=.lastTimestamp

# The tunnel occupies the same port as "make run", so only one of the two
# can be active. "make stop" frees the port again in both cases.
forward: kubeconfig-check ## Tunnel the local API port to the Service in SKE
	kubectl port-forward \
		--namespace $(API_NAMESPACE) \
		service/$(APP) \
		$(PORT):80

cluster-token: kubeconfig-check ## Print the API token stored in the cluster Secret
	@token="$(CLUSTER_TOKEN)"; \
	if [ -z "$$token" ]; then \
		echo "Secret $(APP)-auth was not found in namespace $(API_NAMESPACE)."; \
		exit 1; \
	fi; \
	echo "$$token"

# Every prerequisite of the live demo in one command. Each condition is
# checked separately so that a failure names the part that needs work.
preflight: ## Check everything the live demo depends on
	@failed=0; \
	check() { \
		if eval "$$2" >/dev/null 2>&1; then \
			printf '  ok    %s\n' "$$1"; \
		else \
			printf '  FAIL  %s\n' "$$1"; \
			failed=1; \
		fi; \
	}; \
	echo; \
	echo "Checking the demo prerequisites."; \
	echo; \
	check "login kubeconfig exists" \
		'test -f "$(KUBECONFIG)"'; \
	check "STACKIT session is valid and the API server answers" \
		'kubectl get --raw /version'; \
	check "namespace $(API_NAMESPACE) exists" \
		'kubectl get namespace $(API_NAMESPACE)'; \
	check "namespace $(PAAS_NAMESPACE) exists" \
		'kubectl get namespace $(PAAS_NAMESPACE)'; \
	check "CloudNativePG is installed" \
		'kubectl get crd clusters.postgresql.cnpg.io'; \
	check "Secret $(APP)-auth holds the API token" \
		'kubectl get secret $(APP)-auth --namespace $(API_NAMESPACE)'; \
	check "Secret stackit-registry can pull the image" \
		'kubectl get secret stackit-registry --namespace $(API_NAMESPACE)'; \
	check "Deployment $(APP) is available" \
		'kubectl wait --for=condition=Available deployment/$(APP) --namespace $(API_NAMESPACE) --timeout=10s'; \
	check "Service $(APP) has a ready endpoint" \
		'test -n "$$(kubectl get endpoints $(APP) --namespace $(API_NAMESPACE) --output jsonpath=\{.subsets[0].addresses[0].ip\})"'; \
	check "the ServiceAccount may manage databases" \
		'kubectl auth can-i create clusters.postgresql.cnpg.io --namespace $(PAAS_NAMESPACE) --as=system:serviceaccount:$(API_NAMESPACE):$(APP)'; \
	check "the ServiceAccount may not touch nodes" \
		'! kubectl auth can-i delete nodes --as=system:serviceaccount:$(API_NAMESPACE):$(APP)'; \
	check "the manifests match the cluster" \
		'kubectl diff -f $(DEPLOY_DIR)/02-api.yaml'; \
	echo; \
	if [ "$$failed" -ne 0 ]; then \
		echo "Some checks failed. Fix them before the demo."; \
		exit 1; \
	fi; \
	if lsof -ti tcp:$(PORT) -sTCP:LISTEN >/dev/null 2>&1; then \
		echo "Everything is ready. Port $(PORT) is in use, so a tunnel or"; \
		echo "a local server is already running."; \
	else \
		echo "Everything is ready. Start the tunnel with: make forward"; \
	fi; \
	echo

##@ Documentation

openapi: ## Check that api/openapi.yaml is valid
	@ruby -ryaml -e 'spec = YAML.load_file("api/openapi.yaml"); printf("OpenAPI %s, %d documented paths, YAML is valid\n", spec["openapi"], spec["paths"].size)'

docs: openapi ## Copy the contract to the clipboard and open Swagger Editor
	@pbcopy < api/openapi.yaml
	@open https://editor.swagger.io
	@echo "The contract is on the clipboard. Paste it into the editor."

##@ Container

# --load is required because buildx defaults to the docker-container driver,
# which otherwise leaves the result in the build cache only.
docker-build: ## Build the container image
	docker build --load -t $(IMAGE) .

docker-kubeconfig: ## Write a short-lived kubeconfig the container can use on its own
	@mkdir -p "$$(dirname "$(DOCKER_KUBECONFIG)")"
	stackit ske kubeconfig create $(CLUSTER) --assume-yes \
		--project-id $(PROJECT_ID) \
		--expiration $(DOCKER_KUBECONFIG_TTL) \
		--filepath "$(DOCKER_KUBECONFIG)"

# The container runs with the host user id so that it may read
# the mounted kubeconfig, which is readable only by its owner.
docker-run: docker-kubeconfig ## Run the container image against a mounted kubeconfig
	docker run --rm -p $(PORT):8080 \
		--user "$$(id -u):$$(id -g)" \
		-e API_TOKEN="$(API_TOKEN)" \
		-e PAAS_NAMESPACE="$(PAAS_NAMESPACE)" \
		-e KUBECONFIG=/kubeconfig \
		-v "$$(realpath "$(DOCKER_KUBECONFIG)")":/kubeconfig:ro \
		$(IMAGE)

docker-shell: ## Open a shell inside the image for inspection
	docker run --rm -it --entrypoint /bin/sh $(IMAGE)

##@ Environment

# The path is resolved by realpath and not by make's abspath, because
# make splits paths on whitespace and the home directory contains a tab.
env: ## Print the exports for a plain shell: eval "$(make env)"
	@printf 'export KUBECONFIG="%s"\n' \
		"$$(realpath "$(KUBECONFIG)" 2>/dev/null || echo "$(KUBECONFIG)")"
	@echo 'export API_TOKEN="$(API_TOKEN)"'
	@echo 'export ADMIN_USERNAME="$(ADMIN_USERNAME)"'
	@echo 'export ADMIN_PASSWORD="$(ADMIN_PASSWORD)"'
	@echo 'export JWT_SECRET="$(JWT_SECRET)"'
	@echo 'export PAAS_NAMESPACE="$(PAAS_NAMESPACE)"'
	@echo 'export PORT="$(PORT)"'

# A recipe runs in a child shell and cannot change the environment of the
# shell that called it, so the exports have to be evaluated by the caller.
# A terminal on stdout means nobody captured the output and the exports
# would be thrown away, so that case prints the correct command instead.
load: ## Load the exports into the current shell: eval "$(make load)"
	@if [ -t 1 ]; then \
		echo "make cannot export into the shell that called it."; \
		echo "Run this instead:"; \
		echo '  eval "$$(make $@)"'; \
		exit 1; \
	fi
	@$(MAKE) --no-print-directory env

# Same exports, except that API_TOKEN is read from the Secret the Pod
# uses, so plain curl reaches the API running in SKE through the tunnel.
load-remote: ## Load the exports with the token of the Pod in SKE
	@if [ -t 1 ]; then \
		echo "make cannot export into the shell that called it."; \
		echo "Run this instead:"; \
		echo '  eval "$$(make $@)"'; \
		exit 1; \
	fi
	@$(MAKE) --no-print-directory env REMOTE=1

token: ## Print the bearer token used by the curl helpers
	@echo "$(API_TOKEN)"
