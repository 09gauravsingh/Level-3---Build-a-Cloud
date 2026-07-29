# LEVEL3 Week 2 — SKE PaaS Platform

This Terraform project provisions a multi-zone STACKIT Kubernetes Engine
cluster for the Week 2 Platform-as-a-Service project.

## Planned architecture

- Cluster: lvl3-paas
- Region: eu01
- Node pool: paas-workers
- Initial workers: 3
- Maximum workers: 6
- Availability zones:
  - eu01-1
  - eu01-2
  - eu01-3

The SKE infrastructure is managed separately from the Week 1 DevStack and
kubeadm environments.
