-include .env.secrets

KUBECONFIG     ?= $(HOME)/.kube/lvl3-paas-login.yaml
CLUSTER        ?= lvl3-paas
PROJECT_ID     ?= 20a60c06-1e1a-406f-a840-37d1ff14f0e8

ADMIN_USERNAME ?= admin
ADMIN_PASSWORD ?= week4-password
JWT_SECRET     ?= my-week4-super-secret-key
PAAS_NAMESPACE ?= database-services
PORT           ?= 8080

APP            ?= week3-paas-api
API_NAMESPACE  ?= paas-system
API_IMAGE      ?= registry.onstackit.cloud/gaurav-paas-20a60c06/week3-paas-api:v1.4.10

export KUBECONFIG ADMIN_USERNAME ADMIN_PASSWORD JWT_SECRET PAAS_NAMESPACE PORT

.DEFAULT_GOAL := help

.PHONY: help run stop test vet kube-refresh deploy status logs forward \
	connect-prometheus connect-grafana publish

help:
	@echo "run                  start the API locally"
	@echo "stop                 free port $(PORT)"
	@echo "test                 run unit tests"
	@echo "vet                  run go vet ./..."
	@echo "kube-refresh         write a fresh login kubeconfig"
	@echo "deploy               apply manifests to the cluster"
	@echo "status               show Deployment, Pods and Service"
	@echo "logs                 follow Pod logs"
	@echo "forward              tunnel localhost:$(PORT) to the API"
	@echo "connect-prometheus   tunnel localhost:9090 to Prometheus"
	@echo "connect-grafana      tunnel localhost:3000 to Grafana"
	@echo "publish              build, push and roll out API_IMAGE"

run:
	go run ./cmd/api

stop:
	@lsof -ti tcp:$(PORT) | xargs kill 2>/dev/null || true

test:
	go test ./...

vet:
	go vet ./...

kube-refresh:
	stackit ske kubeconfig create $(CLUSTER) --login --assume-yes \
		--project-id $(PROJECT_ID) \
		--filepath "$(KUBECONFIG)"

deploy:
	kubectl apply -f deploy
	kubectl rollout status deployment/$(APP) -n $(API_NAMESPACE)

status:
	kubectl get deployment,pods,service -n $(API_NAMESPACE) -o wide

logs:
	kubectl logs deployment/$(APP) -n $(API_NAMESPACE) -f --tail=50

forward:
	kubectl port-forward -n $(API_NAMESPACE) service/$(APP) $(PORT):80

connect-prometheus:
	kubectl port-forward -n monitoring svc/monitoring-kube-prometheus-prometheus 9090:9090

connect-grafana:
	kubectl port-forward -n monitoring svc/monitoring-grafana 3000:80

publish:
	podman build --platform linux/amd64 -t "$(API_IMAGE)" .
	podman push $(API_IMAGE)
	kubectl set image deployment/$(APP) api=$(API_IMAGE) -n $(API_NAMESPACE)
	kubectl rollout status deployment/$(APP) -n $(API_NAMESPACE)

