# week5-observability

In-cluster metrics and logs for the PostgreSQL PaaS. The Go API already
exposes `/metrics` and writes JSON to stdout; this directory is the
pipeline that scrapes, stores and displays that data.

How Prometheus, Loki and Alloy sit next to the API is in
[../architecture.md](../architecture.md).

## Monitoring

`monitoring/values.yaml` is the kube-prometheus-stack Helm values:

- Grafana with a 2Gi persistent volume
- Prometheus, 7-day retention, 10Gi volume
- Alertmanager and SKE control-plane scrapes left off

`monitoring/api-servicemonitor.yaml` tells Prometheus to scrape
`week3-paas-api` in `paas-system` on `/metrics` every 30 seconds. The
series are `total_http_requests` (method, path, status) and
`http_request_duration_seconds` (method, path).

## Logging

| File | Role |
| --- | --- |
| `logging/loki-values.yaml` | Loki install |
| `logging/alloy-values.yaml` | Alloy DaemonSet |
| `logging/config.alloy` | Discover pods on this node, relabel namespace/pod/container/app, push to Loki |
| `logging/loki-datasource.yaml` | Grafana datasource for Loki |

Alloy reads logs through the Kubernetes API (`loki.source.kubernetes`)
rather than a privileged host filesystem, and writes to
`http://loki-gateway.logging.svc.cluster.local/loki/api/v1/push`.
