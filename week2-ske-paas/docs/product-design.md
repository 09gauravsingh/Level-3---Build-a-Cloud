# Level3 Managed PostgreSQL

## Purpose

This project provides a simplified PostgreSQL Platform-as-a-Service product
on STACKIT Kubernetes Engine.

A platform user requests a PostgreSQL database by creating a CloudNativePG
Cluster Custom Resource. The CloudNativePG operator then provisions and
manages the required Kubernetes and PostgreSQL components.

## Initial service profile

- Product name: Level3 Managed PostgreSQL
- Kubernetes namespace: database-services
- PostgreSQL cluster: level3-postgres
- PostgreSQL instances: 2
- Persistent storage: 5 GiB per instance
- Application database: platformdb
- Application user: platformuser
- External demonstration: kubectl port-forward
- Internal connectivity: Kubernetes DNS and ClusterIP Service

## Responsibilities

The platform manages:

- PostgreSQL instances
- Primary and standby roles
- Persistent volumes
- Services
- Application credentials
- Health status
- Reconciliation and recovery

The product consumer supplies:

- Desired number of instances
- Storage size
- Database name
- Database owner
- Resource requirements

## Limitations

This is a learning PaaS prototype. Both SKE worker nodes are currently in the
same availability zone, eu01-1. It therefore provides node-level redundancy
but not availability-zone-level redundancy.

Backups, monitoring, tenant isolation, billing, disaster recovery and
production security policies are outside the first implementation stage.
