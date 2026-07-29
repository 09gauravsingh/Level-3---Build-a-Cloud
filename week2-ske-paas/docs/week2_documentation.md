# Week 2 Technical Documentation

## Building a Managed PostgreSQL PaaS on STACKIT Kubernetes Engine

**Terraform-based SKE infrastructure, CloudNativePG Operator, secure in-cluster connectivity, replication and automatic failover**

<img src="docs/assets/media/image1.png" style="width:6.29921in;height:1.18039in" />

Week 2 solution at a glance

| **Document field**         | **Value**                                                                                                                      |
|----------------------------|--------------------------------------------------------------------------------------------------------------------------------|
| Project                    | LEVEL3 - Build a Cloud by STACKIT                                                                                              |
| Author                     | Kumar Gaurav                                                                                                                   |
| STACKIT project            | Gaurav (region eu01)                                                                                                           |
| Implementation environment | Mac workstation; SKE cluster managed through Terraform and kubectl                                                             |
| Document version           | 1.0 - verified implementation through Task 5.7                                                                                 |
| Status date                | 29 July 2026                                                                                                                   |
| Current status             | SKE, operator, PostgreSQL product, connectivity, replication and failover completed; persistence/recovery finalisation remains |

| **Important scope statement:** This document records the verified Week 2 implementation up to and including the successful automatic failover test. It does not claim that backup/restore or disaster recovery has been completed. |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

# Contents

1. Executive summary and Week 2 outcome

2. Scope, learning objectives and completion status

3. Final architecture and responsibility model

4. Working environment, repository and security rules

5. Phase A - Provision STACKIT Kubernetes Engine with Terraform

6. Phase B - Extend Kubernetes with the CloudNativePG Operator

7. Phase C - Define the managed PostgreSQL product

8. Phase D - Demonstrate connectivity and product usage

9. Phase E - Demonstrate replication and automatic failover

10. Validation against the Week 2 requirements

11. Troubleshooting, corrections and lessons learned

12. Security review and production-readiness gaps

13. Remaining work and recommended next task

14. Operational command reference

15. Evidence index and assessment checklist

16. Knowledge-check questions

References and appendices

| **Reading guide:** Sections 1-4 explain the design. Sections 5-9 form the implementation runbook. Sections 10-15 provide assessment evidence, troubleshooting and operational reference. |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

# 1. Executive summary and Week 2 outcome

The Week 2 goal was to move beyond virtual machines and build a platform service on managed Kubernetes. The implemented product is a self-built managed PostgreSQL service running on STACKIT Kubernetes Engine (SKE). Infrastructure is declared with Terraform, PostgreSQL lifecycle management is delegated to the CloudNativePG Operator, and client access is demonstrated from a separate Kubernetes Pod.

This is a PaaS-style platform layer because a consumer does not manually create PostgreSQL processes, configure replication, attach storage, or track the active primary. The consumer declares a PostgreSQL Cluster Custom Resource, while the operator continuously creates and repairs the required Kubernetes and PostgreSQL resources.

| **Verified result:** The SKE cluster is operational with two worker nodes. CloudNativePG created a healthy two-instance PostgreSQL cluster. Writes through the read-write Service replicated to the standby. A write through the read-only Service was rejected as expected. Deleting the current primary triggered automatic promotion of the standby and the cluster recovered to 2/2 Ready. |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

## 1.1 Final verified state

| **Area**            | **Verified state**                                                                                         |
|---------------------|------------------------------------------------------------------------------------------------------------|
| Infrastructure      | STACKIT SKE cluster `lvl3-paas` in `eu01`, managed through Terraform from the Mac                      |
| Workers             | Node pool `paas-workers`, two fixed `g3i.2` workers in availability zone `eu01-1`                    |
| Operator            | CloudNativePG installed with Helm; CRDs and controller verified                                            |
| Product namespace   | `database-services`                                                                                      |
| PostgreSQL product  | `level3-postgres`, PostgreSQL 18.4, two instances, 5 GiB persistent storage per instance                 |
| Database identity   | Database `platformdb`, application owner `platformuser`                                                |
| Connectivity        | `level3-psql-client` uses Secret-injected credentials and Kubernetes Service DNS                         |
| Read/write endpoint | `level3-postgres-rw` routes to the current primary                                                       |
| Read-only endpoint  | `level3-postgres-ro` routes to the standby                                                               |
| Failover evidence   | Primary changed from `level3-postgres-1` to `level3-postgres-2`; cluster returned to healthy 2/2 state |
| Evidence files      | `evidence/5.7-before-failover.txt` and `evidence/5.7-after-failover.txt`                               |

# 2. Scope, learning objectives and completion status

## 2.1 What Week 2 was designed to teach

- **Managed Kubernetes:** understand the boundary between the cloud provider, Kubernetes control plane, worker nodes and workloads.

- **Infrastructure as Code:** create repeatable SKE infrastructure using Terraform rather than manual portal operations.

- **Kubernetes extensibility:** understand CustomResourceDefinitions, Custom Resources, controllers, operators and reconciliation.

- **PaaS product design:** expose a database service through a simple declarative product definition.

- **Stateful workloads:** use persistent volumes, stable Services, Secrets and replication.

- **Operational validation:** connect to the product, run SQL, verify read/write behaviour and test failover.

- **Documentation and evidence:** retain commands and outputs that prove the platform works.

## 2.2 Completion matrix

| **Work package**                   | **Status**      | **Evidence / note**                                                             |
|------------------------------------|-----------------|---------------------------------------------------------------------------------|
| Mac-local Week 2 workspace         | Complete        | Authoritative files and state are under `~/build-a-cloud/week2-ske-paas`      |
| SKE Terraform configuration        | Complete        | Cluster and node pool created; Terraform state is Mac-local                     |
| kubectl access                     | Complete        | `KUBECONFIG` points to `01-ske-infrastructure/kubeconfig.yaml`; nodes Ready |
| CloudNativePG Operator             | Complete        | Helm installation and CRDs verified                                             |
| Managed PostgreSQL Custom Resource | Complete        | `level3-postgres` healthy with 2/2 instances                                  |
| Connectivity demonstration         | Complete        | Client Pod connected as `platformuser` to `platformdb`                      |
| Product usage demonstration        | Complete        | Table creation, insert and select succeeded                                     |
| Replication and read-only test     | Complete        | Data visible through `-ro`; write rejected on standby                         |
| Automatic failover                 | Complete        | Primary changed and cluster recovered                                           |
| Persistence/recovery final test    | Pending         | Next task; do not confuse Pod survival with backup/restore                      |
| Backup and restore                 | Not implemented | Production gap; requires a backup design and restore test                       |

# 3. Final architecture and responsibility model

<img src="docs/assets/media/image1.png" style="width:6.37795in;height:1.19514in" />

Figure 1 - Complete Week 2 architecture

## 3.1 End-to-end request flow

1.  Terraform running on the Mac authenticates to STACKIT and declares the SKE cluster and node pool.

2.  STACKIT manages the Kubernetes control plane; the user receives a kubeconfig for API access.

3.  Helm installs CloudNativePG and its CRDs in the SKE cluster.

4.  The `level3-postgres` Custom Resource declares the desired database product.

5.  CloudNativePG creates PostgreSQL Pods, persistent volume claims, Services, Secrets and replication configuration.

6.  The client Pod resolves `level3-postgres-rw` or `level3-postgres-ro` through Kubernetes DNS and authenticates with credentials injected from a Secret.

7.  During failover, CloudNativePG promotes a standby and changes Service endpoints; the client continues to use the same stable Service name.

## 3.2 Who manages what?

| **Layer**                 | **Main responsibility**                                                                                       |
|---------------------------|---------------------------------------------------------------------------------------------------------------|
| STACKIT cloud             | Physical infrastructure, networking foundation, storage services and SKE service APIs                         |
| STACKIT Kubernetes Engine | Managed Kubernetes control plane and cluster lifecycle integration                                            |
| Terraform                 | Declarative source of truth for the SKE infrastructure and generated kubeconfig workflow                      |
| Kubernetes                | Scheduling, Pod lifecycle, Services, Secrets, PersistentVolumeClaims and API state                            |
| CloudNativePG Operator    | PostgreSQL topology, initialization, credentials/certificates, replication, role changes, health and failover |
| PostgreSQL                | SQL execution, transactions, database objects and physical streaming replication                              |
| Platform engineer         | Product policy, manifests, resource sizing, access controls, testing, monitoring, backup design and evidence  |

| **Correct interpretation:** This is not the same as buying a provider-native managed database. It is a PaaS-style database product built by the platform engineer on top of managed Kubernetes. |
|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

# 4. Working environment, repository and security rules

## 4.1 Authoritative execution environment

All Week 2 work is performed directly from the Mac. The earlier `gaurav-openstack` server is not the source of truth for this phase. Terraform state, lock file, kubeconfig and manifests must remain aligned in the Mac-local workspace.

    cd ~/build-a-cloud/week2-ske-paas
    uname -s
    # Expected: Darwin

## 4.2 Repository structure

    week2-ske-paas/
    ├── 01-ske-infrastructure/
    │   ├── versions.tf
    │   ├── provider.tf
    │   ├── variables.tf
    │   ├── cluster.tf
    │   ├── kubeconfig.tf
    │   ├── outputs.tf
    │   ├── .terraform.lock.hcl
    │   ├── terraform.tfstate               # sensitive, local only
    │   ├── kubeconfig.yaml                 # sensitive, local only
    │   └── evidence/
    ├── 02-operator/
    │   ├── chart-version.txt
    │   ├── cnpg-airgap-resources.txt
    │   ├── helm-install.log
    │   └── cnpg-install.log
    ├── 03-product/
    │   ├── namespace.yaml
    │   ├── level3-postgres.yaml
    │   └── psql-client.yaml
    ├── 04-connectivity/
    │   └── test-postgres.sql
    ├── 05-demo/
    ├── docs/
    └── evidence/
        ├── 5.7-before-failover.txt
        └── 5.7-after-failover.txt

## 4.3 Source-of-truth rule

- Terraform is the source of truth for SKE infrastructure.

- YAML manifests are the source of truth for Kubernetes/operator product resources.

- The STACKIT Portal is used to observe and verify the result, not to recreate the same resources manually.

- Terminal output and screenshots are evidence, but they are not configuration source files.

## 4.4 Sensitive-file handling

| **Never commit:** `terraform.tfstate`, state backups, kubeconfig files, service-account private keys, Terraform plan files and any file containing decoded Secret values. |
|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

    # Recommended .gitignore concepts
    .terraform/
    *.tfstate
    *.tfstate.*
    *.tfplan
    kubeconfig*.yaml
    kubeconfig-expired-*
    *service-account*.json
    *.key
    .DS_Store

The provider lock file `.terraform.lock.hcl` is different: it should normally be committed because it records selected provider versions and checksums, improving reproducibility.

# 5. Phase A - Provision STACKIT Kubernetes Engine with Terraform

## 5.1 Why use SKE?

SKE provides Kubernetes as a managed cloud service. STACKIT operates the control-plane layer, while the platform engineer defines worker capacity and deploys workloads. This reduces the operational burden compared with creating and maintaining the Kubernetes control plane manually.

## 5.2 Actual infrastructure values

| **Setting**          | **Value**                            | **Purpose**                                              |
|----------------------|--------------------------------------|----------------------------------------------------------|
| STACKIT project name | Gaurav                               | Cloud project containing the SKE resources               |
| Project ID           | 20a60c06-1e1a-406f-a840-37d1ff14f0e8 | API identifier used by Terraform                         |
| Region               | eu01                                 | STACKIT region                                           |
| Cluster name         | lvl3-paas                            | Short, stable SKE cluster identity                       |
| Node pool            | paas-workers                         | Worker group for platform workloads                      |
| Machine type         | g3i.2                                | Worker capacity selected for the learning platform       |
| Worker count         | minimum 2 / maximum 2                | Fixed two-worker design; no autoscaling in this exercise |
| Availability zone    | eu01-1                               | Worker placement zone                                    |
| Maintenance time     | approximately 01:00                  | Planned maintenance window configuration                 |

## 5.3 Terraform file responsibilities

| **File**                | **Responsibility**                                                                                         |
|-------------------------|------------------------------------------------------------------------------------------------------------|
| `versions.tf`         | Declares Terraform and provider requirements.                                                              |
| `provider.tf`         | Configures the STACKIT provider and default region/authentication behaviour.                               |
| `variables.tf`        | Defines reusable inputs such as project ID, cluster name, worker type and counts.                          |
| `cluster.tf`          | Declares the `stackit_ske_cluster` and node-pool desired state.                                          |
| `kubeconfig.tf`       | Requests current SKE access credentials and writes `kubeconfig.yaml` with restrictive local permissions. |
| `outputs.tf`          | Prints useful non-secret identifiers and cluster information after apply.                                  |
| `.terraform.lock.hcl` | Pins provider selections/checksums for reproducible runs.                                                  |
| `terraform.tfstate`   | Maps Terraform configuration to real cloud objects; sensitive and authoritative.                           |

## 5.4 Terraform lifecycle used in the project

    cd ~/build-a-cloud/week2-ske-paas/01-ske-infrastructure

    terraform init
    terraform fmt -recursive
    terraform validate
    terraform plan
    terraform apply

| **Command**            | **Meaning**                                                           |
|------------------------|-----------------------------------------------------------------------|
| `terraform init`     | Downloads providers/modules and initializes the working directory.    |
| `terraform fmt`      | Normalizes Terraform formatting; does not change cloud resources.     |
| `terraform validate` | Checks syntax and internal configuration consistency.                 |
| `terraform plan`     | Calculates the proposed changes without applying them.                |
| `terraform apply`    | Executes the approved plan through the STACKIT API and updates state. |

## 5.5 Kubeconfig generation and use

The SKE kubeconfig contains the API server address, certificate authority data, client credentials, context definitions and current context. It is authentication material and must be protected. Refreshing the kubeconfig resource changes access credentials; it does not recreate the SKE cluster.

    cd ~/build-a-cloud/week2-ske-paas
    export KUBECONFIG="$HOME/build-a-cloud/week2-ske-paas/01-ske-infrastructure/kubeconfig.yaml"

    kubectl config current-context
    kubectl cluster-info
    kubectl get nodes -o wide

| **Expected result:** The current context identifies `lvl3-paas`, and both SKE workers report `Ready`. |
|-----------------------------------------------------------------------------------------------------------|

## 5.6 Portal verification

After Terraform apply, the same cluster should be visible in the STACKIT Portal under Runtime -/> Kubernetes Engine. Portal verification confirms that the declared cluster name, region, node pool and worker count match Terraform. No duplicate manual cluster should be created.

# 6. Phase B - Extend Kubernetes with the CloudNativePG Operator

## 6.1 Built-in resource, CRD, CR and Operator

| **Term**                       | **Simple definition**                                                           | **Project example**                                        |
|--------------------------------|---------------------------------------------------------------------------------|------------------------------------------------------------|
| Built-in resource              | A type Kubernetes already understands.                                          | Pod, Service, Secret, Namespace, PersistentVolumeClaim     |
| CustomResourceDefinition (CRD) | Adds a new API resource type to Kubernetes.                                     | `clusters.postgresql.cnpg.io`                            |
| Custom Resource (CR)           | One object created using the new type.                                          | `kind: Cluster`, name `level3-postgres`                |
| Controller                     | A loop that watches resources and moves actual state toward desired state.      | Kubernetes controllers                                     |
| Operator                       | A specialised controller that encodes operational knowledge for an application. | CloudNativePG                                              |
| Reconciliation                 | Repeated comparison and repair of desired versus actual state.                  | Create Pods/PVCs, promote standby, rebuild failed instance |

<img src="docs/assets/media/image2.png" style="width:6.14173in;height:1.33516in" />

Figure 2 - CloudNativePG reconciliation loop

## 6.2 Operator installation with Helm

    helm version
    helm repo add cnpg https://cloudnative-pg.github.io/charts
    helm repo update

    helm install cnpg   --namespace cnpg-system   --create-namespace   cnpg/cloudnative-pg

The project retains chart/version and installation logs under `02-operator/` so the installation can be audited without relying only on terminal history.

## 6.3 Verification commands

    kubectl get pods --namespace cnpg-system
    kubectl get deployments --namespace cnpg-system
    kubectl get crd | grep postgresql.cnpg.io
    kubectl describe crd clusters.postgresql.cnpg.io
    kubectl explain cluster
    kubectl explain cluster.spec

| **What successful installation means:** Kubernetes now recognises `apiVersion: postgresql.cnpg.io/v1` and `kind: Cluster`, and the CloudNativePG controller is running and watching these resources. |
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

# 7. Phase C - Define the managed PostgreSQL product

## 7.1 Product namespace

The dedicated namespace separates database platform resources from unrelated workloads and gives a clean boundary for listing, access control, policies and evidence.

    apiVersion: v1
    kind: Namespace
    metadata:
      name: database-services
      labels:
        platform.level3.io/purpose: managed-databases
    kubectl apply --filename 03-product/namespace.yaml
    kubectl get namespace database-services

## 7.2 Exact PostgreSQL product manifest

    apiVersion: postgresql.cnpg.io/v1
    kind: Cluster

    metadata:
      name: level3-postgres
      namespace: database-services

      labels:
        platform.level3.io/product: managed-postgresql
        platform.level3.io/environment: learning

    spec:
      description: Level3 managed PostgreSQL PaaS learning service

      # One primary and one standby PostgreSQL instance.
      instances: 2

      # Fixed image for reproducible PostgreSQL versioning.
      imageName: ghcr.io/cloudnative-pg/postgresql:18.4-system-trixie

      bootstrap:
        initdb:
          database: platformdb
          owner: platformuser
          dataChecksums: true

      storage:
        size: 5Gi

      resources:
        requests:
          cpu: 250m
          memory: 512Mi
        limits:
          cpu: "1"
          memory: 1Gi

      affinity:
        enablePodAntiAffinity: true
        topologyKey: kubernetes.io/hostname
        podAntiAffinityType: required

## 7.3 Manifest explanation

| **Field**                     | **Meaning in this product**                                                            |
|-------------------------------|----------------------------------------------------------------------------------------|
| `apiVersion`                | Selects the API group/version added by the CloudNativePG CRD.                          |
| `kind: Cluster`             | Requests a CloudNativePG-managed PostgreSQL cluster, not a generic Kubernetes cluster. |
| `metadata.name`             | Stable product identity used as the prefix for generated resources.                    |
| `namespace`                 | Places all namespaced product objects in `database-services`.                        |
| `labels`                    | Platform metadata for selection, ownership and learning environment classification.    |
| `instances: 2`              | Desired topology is one primary plus one standby.                                      |
| `imageName`                 | Pins PostgreSQL 18.4 to improve reproducibility.                                       |
| `bootstrap.initdb.database` | Creates the initial application database `platformdb`.                               |
| `bootstrap.initdb.owner`    | Creates `platformuser` as the application database owner.                            |
| `dataChecksums: true`       | Enables PostgreSQL page checksums for detecting certain storage corruption.            |
| `storage.size: 5Gi`         | Requests one persistent 5 GiB volume for each PostgreSQL instance.                     |
| `resources.requests`        | Scheduler reservation used to decide whether a worker has enough CPU/memory.           |
| `resources.limits`          | Maximum container CPU/memory consumption permitted by Kubernetes.                      |
| Required Pod anti-affinity    | Prevents both PostgreSQL instances from being scheduled on the same worker hostname.   |

## 7.4 Create the product

    kubectl apply --filename 03-product/level3-postgres.yaml

    kubectl get cluster level3-postgres   --namespace database-services   --watch

The apply command creates only the Custom Resource. CloudNativePG then performs the operational work. The watch command does not change anything; it keeps displaying status changes until Control+C is pressed.

## 7.5 Resources generated by the Operator

| **Generated object**               | **Purpose**                                                                       |
|------------------------------------|-----------------------------------------------------------------------------------|
| PostgreSQL Pods                    | Run the primary and standby database instances.                                   |
| PVC per instance                   | Store PostgreSQL data independently of the Pod/container lifetime.                |
| `level3-postgres-rw` Service     | Stable endpoint for the current primary; accepts reads and writes.                |
| `level3-postgres-ro` Service     | Stable endpoint for ready standby replicas; read-only workload path.              |
| `level3-postgres-r` Service      | Endpoint covering readable PostgreSQL instances.                                  |
| `level3-postgres-app` Secret     | Application username/password and connection information generated by convention. |
| CA, server and replication Secrets | TLS identity and secure replication support.                                      |
| Status/events                      | Expose cluster phase, current primary and reconciliation history.                 |

    kubectl get cluster,pods,services,pvc,secrets   --namespace database-services

    kubectl describe cluster level3-postgres   --namespace database-services

# 8. Phase D - Demonstrate connectivity and product usage

## 8.1 Connectivity design decision

The original local-client approach would have installed `psql` on the Mac and used a port-forward. Because the Mac did not have administrator access for installation and the goal was to test realistic in-cluster connectivity, the workflow was changed: `psql` runs inside a dedicated client Pod, while the Mac runs only `kubectl`.

| **Why this is better for the demonstration:** It proves that another Kubernetes workload can discover the database by Service name, receive credentials from a Secret and use the same internal network path an application would use. |
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

## 8.2 Exact client Pod manifest

    apiVersion: v1
    kind: Pod

    metadata:
      name: level3-psql-client
      namespace: database-services

      labels:
        platform.level3.io/component: postgresql-client
        platform.level3.io/environment: learning

    spec:
      restartPolicy: Never

      containers:
        - name: psql
          image: ghcr.io/cloudnative-pg/postgresql:18.4-system-trixie
          imagePullPolicy: IfNotPresent

          command:
            - sleep
            - infinity

          env:
            - name: PGHOST
              value: level3-postgres-rw

            - name: PGPORT
              value: "5432"

            - name: PGDATABASE
              value: platformdb

            - name: PGUSER
              valueFrom:
                secretKeyRef:
                  name: level3-postgres-app
                  key: username

            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: level3-postgres-app
                  key: password

            - name: PGSSLMODE
              value: require

            - name: PGAPPNAME
              value: level3-psql-client

          resources:
            requests:
              cpu: 10m
              memory: 32Mi
            limits:
              cpu: 100m
              memory: 128Mi

## 8.3 Create and enter the client Pod

    kubectl apply --filename 03-product/psql-client.yaml

    kubectl wait   --namespace database-services   --for=condition=Ready   pod/level3-psql-client   --timeout=120s

    kubectl exec   --namespace database-services   --stdin   --tty   level3-psql-client   --container psql   -- psql

In `kubectl exec ... -- psql`, everything before `--` configures Kubernetes execution. Everything after `--` is the command executed inside the container. The environment variables remove the need to type the host, port, database, username or password manually.

## 8.4 Connection path

<img src="docs/assets/media/image3.png" style="width:6.29921in;height:0.72135in" />

Figure 3 - Secure in-cluster connectivity path

## 8.5 SQL usage demonstration

    CREATE TABLE IF NOT EXISTS platform_service_test (
        service_id TEXT PRIMARY KEY,
        service_name TEXT NOT NULL,
        service_tier TEXT NOT NULL,
        database_owner TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

    INSERT INTO platform_service_test (
        service_id,
        service_name,
        service_tier,
        database_owner
    )
    VALUES (
        'managed-postgresql',
        'Level3 Managed PostgreSQL',
        'PaaS',
        current_user
    )
    ON CONFLICT (service_id)
    DO UPDATE SET
        service_name = EXCLUDED.service_name,
        service_tier = EXCLUDED.service_tier,
        database_owner = EXCLUDED.database_owner;

    SELECT service_id, service_name, service_tier, database_owner, created_at
    FROM platform_service_test;

| **What this proves:** The client authenticated as the application user, connected to the intended database, created a schema object, wrote data and read it back. |
|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|

## 8.6 Repository SQL file

    /set ON_ERROR_STOP on

    CREATE TABLE IF NOT EXISTS platform_connectivity_test (
        id INTEGER PRIMARY KEY,
        message TEXT NOT NULL,
        checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

    INSERT INTO platform_connectivity_test (id, message, checked_at)
    VALUES (
        1,
        'Connection to the Level3 managed PostgreSQL service works',
        NOW()
    )
    ON CONFLICT (id)
    DO UPDATE SET
        message = EXCLUDED.message,
        checked_at = NOW();

    SELECT id, message, checked_at
    FROM platform_connectivity_test
    ORDER BY id;

`cat 04-connectivity/test-postgres.sql` only prints the file. It does not execute SQL. To execute the local file while keeping `psql` inside Kubernetes, pipe the file into the Pod:

    kubectl exec   --namespace database-services   --stdin   level3-psql-client   --container psql   -- psql < 04-connectivity/test-postgres.sql

## 8.7 Primary and standby verification

    # Read-write Service: expected is_standby = f
    kubectl exec -n database-services level3-psql-client -c psql --   psql --host=level3-postgres-rw   --command="SELECT current_user, pg_is_in_recovery() AS is_standby;"

    # Read-only Service: expected is_standby = t
    kubectl exec -n database-services level3-psql-client -c psql --   psql --host=level3-postgres-ro   --command="SELECT current_user, pg_is_in_recovery() AS is_standby;"

| **Value**                   | **Meaning**                                                      |
|-----------------------------|------------------------------------------------------------------|
| `pg_is_in_recovery() = f` | The connected PostgreSQL instance is the primary.                |
| `pg_is_in_recovery() = t` | The connected PostgreSQL instance is a standby in recovery mode. |

## 8.8 Replication and standby protection

A row written through the `-rw` Service was successfully selected through the `-ro` Service. This proves that the standby had received the committed data through PostgreSQL replication. An attempted INSERT through `-ro` returned `ERROR: cannot execute INSERT in a read-only transaction`. That error is a successful test result, not a platform failure.

<img src="docs/assets/media/image4.png" style="width:6.37795in;height:2.95093in" />

Figure 4 - Connectivity, replicated read, expected standby write rejection and final resource state

## 8.9 Credential handling

The client Pod obtained the `platformuser` username and password from the CloudNativePG-generated `level3-postgres-app` Secret through `secretKeyRef`. The password was not copied into the manifest, printed to the terminal or stored in the documentation.

| **Security nuance:** A Kubernetes Secret improves separation from normal manifests, but base64 is not encryption. Access must still be protected with RBAC, secret-at-rest encryption and least privilege. |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

# 9. Phase E - Demonstrate replication and automatic failover

## 9.1 Test objective

The failover exercise tested whether the database service can recover when the current primary Pod disappears. The test deleted only the primary Pod. It did not delete the Cluster Custom Resource, PVCs, namespace, Services, Secrets or Terraform-managed SKE infrastructure.

## 9.2 Discover the primary dynamically

    PRIMARY_BEFORE=$(kubectl get cluster level3-postgres   --namespace database-services   --output jsonpath='{.status.currentPrimary}')

    echo "Primary before failover: $PRIMARY_BEFORE"

Using `.status.currentPrimary` avoids hard-coding a Pod name that can become incorrect after a role change.

## 9.3 Establish a pre-failover data marker

    kubectl exec -n database-services level3-psql-client -c psql --   psql --host=level3-postgres-rw --set=ON_ERROR_STOP=1   --command="
    INSERT INTO platform_service_test (
      service_id, service_name, service_tier, database_owner
    )
    VALUES (
      'before-failover',
      'Record created before primary failure',
      'HA-Test',
      current_user
    )
    ON CONFLICT (service_id) DO UPDATE SET
      service_name = EXCLUDED.service_name,
      service_tier = EXCLUDED.service_tier,
      database_owner = EXCLUDED.database_owner;"

The row was then read through `level3-postgres-ro` before deleting the primary. This proved that the marker had already reached the standby and made the subsequent survival test meaningful.

## 9.4 Monitor and simulate failure

    # Terminal 2: watch cluster status
    kubectl get cluster level3-postgres   --namespace database-services   --watch

    # Terminal 3: repeatedly connect through the stable rw Service
    while true; do
      printf "/n[%s] " "$(date '+%H:%M:%S')"
      kubectl exec -n database-services level3-psql-client -c psql --     psql --host=level3-postgres-rw --tuples-only --no-align     --command="SELECT current_database(), current_user,
                          inet_server_addr(), pg_is_in_recovery();" 2>&1 || true
      sleep 2
    done

    # Terminal 1: delete only the current primary Pod
    kubectl delete pod "$PRIMARY_BEFORE"   --namespace database-services

## 9.5 What CloudNativePG did

8.  Detected that the current primary was unhealthy or absent.

9.  Selected the standby as the failover target.

10. Promoted `level3-postgres-2` to primary.

11. Updated the `level3-postgres-rw` Service to point to the promoted primary.

12. Recreated and resynchronised `level3-postgres-1` as a standby.

13. Returned the cluster to two instances and two Ready instances.

<img src="docs/assets/media/image5.png" style="width:6.29921in;height:0.48666in" />

Figure 5 - Automatic failover sequence

## 9.6 Verified failover result

| **Check**           | **Observed result**                                                             |
|---------------------|---------------------------------------------------------------------------------|
| Primary before      | `level3-postgres-1`                                                           |
| Primary after       | `level3-postgres-2`                                                           |
| Final cluster phase | `Cluster in healthy state`                                                    |
| Final readiness     | `2` instances, `2` Ready                                                    |
| Pods                | `level3-postgres-1`, `level3-postgres-2` and `level3-psql-client` Running |
| PVCs                | Both 5 GiB claims remained Bound                                                |
| Events              | `FailingOver`, `FailoverTarget`, Pod create/start events                    |
| Evidence            | Before and after output files stored under `evidence/`                        |

<img src="docs/assets/media/image6.png" style="width:6.37795in;height:3.6312in" />

Figure 6 - Saved post-failover evidence and terminal output

## 9.7 RTO and RPO interpretation

| **Concept**                    | **Meaning in this test**                                                                                         |
|--------------------------------|------------------------------------------------------------------------------------------------------------------|
| Recovery Time Objective (RTO)  | Approximate time between the first failed client connection and the first successful connection after promotion. |
| Recovery Point Objective (RPO) | Potential amount of recently committed data that could be lost if it had not reached the standby before failure. |

The pre-failover row was deliberately verified on the standby, so that specific row was expected to survive. This does not prove a universal zero-data-loss guarantee for every failure timing and replication configuration.

## 9.8 High availability is not backup

| **Critical distinction:** Replication and failover protect service availability. They do not protect against accidental DELETE/DROP, logical corruption, compromised credentials, or cluster-wide loss. Those risks require independent backups and tested restore procedures. |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

# 10. Validation against the Week 2 requirements

| **Requirement**            | **How it was satisfied**                                                            | **Status** |
|----------------------------|-------------------------------------------------------------------------------------|------------|
| Managed Kubernetes         | SKE cluster created with two workers and accessed from Mac through kubeconfig.      | Met        |
| Infrastructure as Code     | Terraform files and state manage the SKE cluster; portal used for verification.     | Met        |
| PaaS product               | CloudNativePG-backed managed PostgreSQL product defined by a small Custom Resource. | Met        |
| Operator deployment        | CloudNativePG installed with Helm; CRDs/controller verified.                        | Met        |
| Custom Resource            | `level3-postgres` applied in `database-services`.                               | Met        |
| Connectivity documentation | Client Pod, Service DNS, Secret injection and commands documented.                  | Met        |
| Connectivity demonstration | Successful login to `platformdb` as `platformuser` from another Pod.            | Met        |
| Using the product          | SQL table, insert, update-on-conflict and select demonstrated.                      | Met        |
| Persistent storage         | Two 5 GiB RWO PVCs Bound using `premium-perf1-stackit`.                           | Met        |
| Replication                | Write through `-rw` visible through `-ro`.                                      | Met        |
| Read-only behaviour        | INSERT through standby rejected as expected.                                        | Met        |
| High availability          | Primary Pod deletion caused automatic promotion and service rerouting.              | Met        |
| Recovery evidence          | Cluster recovered to 2/2 and before/after evidence saved.                           | Met        |
| Backup/restore             | Not implemented or tested in the verified scope.                                    | Gap        |

## 10.1 Direct answer to the connectivity rubric

| **Connectivity: documentation and demonstration of connecting to and using the PaaS product:** Completed. The documentation explains the endpoint, credentials, client environment and commands; the demonstration includes authentication, table creation, writes, reads, replication, standby protection and connectivity continuity after failover. |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

# 11. Troubleshooting, corrections and lessons learned

| **Problem / question**                  | **Cause or clarification**                                                                              | **Resolution / lesson**                                                                             |
|-----------------------------------------|---------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `kubectl` tried `localhost:8080`    | `kubectl` did not have the SKE kubeconfig selected.                                                   | Export the Mac-local `KUBECONFIG` and verify current context before any cluster command.          |
| `level3-postgres` not found           | The CRD existed, but the product manifest had only been created locally and not yet applied.            | Apply namespace and product YAML, then query the correct namespace.                                 |
| No permission to install local `psql` | Mac administrator/Homebrew installation was unavailable.                                                | Use a dedicated in-cluster client Pod; no local `psql` or port-forward required.                  |
| Why create a client Pod?                | Executing inside the database Pod proves SQL access but not normal application-to-Service connectivity. | A separate Pod provides a realistic consumer path and avoids coupling to an instance Pod.           |
| `cat test-postgres.sql`               | `cat` reads and prints text only.                                                                     | Pipe the file into `kubectl exec ... -- psql` or use psql `--file` inside a container.          |
| Standby INSERT returned an error        | The `-ro` Service correctly reached a read-only standby.                                              | Treat `cannot execute INSERT in a read-only transaction` as successful evidence.                  |
| Red command exit after expected error   | `ON_ERROR_STOP=1` deliberately makes psql return non-zero on SQL error.                               | Expected for the negative test; explain it in evidence rather than hiding it.                       |
| Hard-coded primary name                 | Roles can change during failover.                                                                       | Read `.status.currentPrimary` into a shell variable.                                              |
| `kubectl ... --watch` appears to hang | Watch mode intentionally keeps running and waits for updates.                                           | Use Control+C after the desired transition is observed.                                             |
| Pod deletion versus data deletion       | Pods are replaceable compute; PostgreSQL data is on PVCs.                                               | Delete only the Pod for failure simulation; never delete CR/PVC without understanding consequences. |

## 11.1 Most important architectural lessons

- Connect applications to stable Services, never to a particular primary Pod name.

- A Custom Resource declares intent; the Operator owns the lifecycle work.

- Persistent volumes separate data lifetime from container and Pod lifetime.

- A Secret should be referenced, not copied into manifests or screenshots.

- A two-instance topology can fail over, but read-only capacity temporarily disappears while the only standby is being promoted and rebuilt.

- A healthy final state is not enough evidence by itself; record before/after state, events and functional SQL tests.

# 12. Security review and production-readiness gaps

## 12.1 Security controls already demonstrated

- Database is exposed through internal ClusterIP Services, not a public LoadBalancer.

- Client uses TLS-required mode (`PGSSLMODE=require`).

- Application credentials are injected from a Kubernetes Secret.

- The product runs under a dedicated namespace.

- CPU and memory requests/limits reduce uncontrolled resource use.

- Required Pod anti-affinity reduces the risk of both database instances sharing one worker.

- Sensitive values are intentionally absent from the documentation and evidence.

## 12.2 Gaps before calling this production-ready

| **Gap**                   | **Why it matters**                                                                 | **Production direction**                                                                             |
|---------------------------|------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| Backup and restore        | HA replicas copy mistakes and corruption as well as good data.                     | Configure independent physical backups/snapshots and perform a restore drill.                        |
| NetworkPolicy             | Namespace placement alone does not restrict which Pods can connect.                | Allow PostgreSQL only from approved application namespaces/labels.                                   |
| RBAC and secret access    | Users able to create/read workloads may indirectly access Secrets.                 | Use least-privilege roles and separate platform/admin/application duties.                            |
| Secret-at-rest protection | Kubernetes Secrets are not automatically equivalent to an external secret manager. | Enable encryption at rest and consider external secret integration/rotation.                         |
| Monitoring and alerting   | Manual `kubectl get` is not continuous operations.                               | Add metrics, dashboards and alerts for replication lag, storage, failover and connection saturation. |
| Capacity testing          | Learning resource values are not workload sizing evidence.                         | Benchmark storage and database workload; tune CPU, memory, IOPS and connections.                     |
| Multi-zone resilience     | Both workers are currently listed in one availability zone.                        | For stronger fault isolation, use supported multi-zone node placement and validate storage topology. |
| Upgrade strategy          | Pinned versions improve reproducibility but still need lifecycle planning.         | Test operator and PostgreSQL upgrades in a non-production environment.                               |

# 13. Remaining work and recommended next task

## 13.1 Task 5.8 - Persistence and recovery verification

The next task should explicitly verify that database data remains available after Pod recreation and that the expected PVC is reattached. Although the failover test already showed persistent storage remaining Bound, a dedicated persistence test should capture the data marker, Pod/PVC identity before and after, and the final query result.

14. Create a uniquely timestamped persistence marker through `level3-postgres-rw`.

15. Record the current primary, Pod UID, PVC name and bound PV.

16. Delete only the selected PostgreSQL Pod using a normal graceful deletion.

17. Wait for CloudNativePG to return the cluster to a healthy state.

18. Verify the PVC remained Bound and data marker still exists.

19. Save before/after output under `evidence/`.

| **Do not overclaim:** This proves persistence across Pod recreation. It still does not prove point-in-time recovery, backup durability or recovery from accidental data deletion. |
|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

## 13.2 Final Week 2 close-out after Task 5.8

- Add the persistence evidence to this document.

- Create a short demo script for mentor presentation.

- Run a final `terraform plan` to confirm no unintended infrastructure drift.

- Run final Kubernetes health and SQL checks.

- Keep the client Pod for demos or delete it when no longer needed.

- Confirm no secrets, kubeconfig or Terraform state are staged for Git.

# 14. Operational command reference

## 14.1 Start every Mac session

    cd ~/build-a-cloud/week2-ske-paas
    export KUBECONFIG="$HOME/build-a-cloud/week2-ske-paas/01-ske-infrastructure/kubeconfig.yaml"
    kubectl config current-context
    kubectl get nodes

## 14.2 Platform health

    kubectl get pods --namespace cnpg-system
    kubectl get cluster,pods,services,pvc,secrets   --namespace database-services
    kubectl describe cluster level3-postgres   --namespace database-services

## 14.3 Interactive SQL

    kubectl exec -it   --namespace database-services   level3-psql-client   --container psql   -- psql

## 14.4 One-shot SQL through primary and standby

    kubectl exec -n database-services level3-psql-client -c psql --   psql --host=level3-postgres-rw   --command="SELECT current_user, pg_is_in_recovery();"

    kubectl exec -n database-services level3-psql-client -c psql --   psql --host=level3-postgres-ro   --command="SELECT current_user, pg_is_in_recovery();"

## 14.5 List tables in psql

    /dt
    /d platform_service_test
    /conninfo
    /q

## 14.6 Identify current primary

    kubectl get cluster level3-postgres   --namespace database-services   --output jsonpath='{.status.currentPrimary}{"/n"}' 

## 14.7 Evidence snapshot

    {
      echo "=== WEEK 2 PLATFORM STATE ==="
      date
      kubectl get cluster level3-postgres -n database-services
      kubectl get pods -n database-services -o wide
      kubectl get services -n database-services
      kubectl get pvc -n database-services
      kubectl get events -n database-services     --sort-by='.lastTimestamp' | tail -n 30
    } | tee evidence/week2-final-state.txt

## 14.8 Remove only the temporary client when finished

    kubectl delete --filename 03-product/psql-client.yaml

Deleting the client Pod does not delete PostgreSQL data. Do not delete the PostgreSQL Cluster Custom Resource or its PVCs as part of routine client cleanup.

# 15. Evidence index and assessment checklist

## 15.1 Evidence already available

| **Evidence**                | **What it proves**                                                               |
|-----------------------------|----------------------------------------------------------------------------------|
| Connectivity screenshot     | Standby identity, replicated read, expected write rejection and final resources. |
| `5.7-before-failover.txt` | Original primary and pre-failure cluster/resource state.                         |
| `5.7-after-failover.txt`  | New primary, healthy 2/2 recovery, Pods, Services, PVCs and failover events.     |
| Terraform files and state   | SKE infrastructure was created and is managed declaratively.                     |
| Operator logs/version files | CloudNativePG installation traceability.                                         |
| Product/client YAML         | Declarative definition of the PaaS product and consumer path.                    |
| SQL file                    | Repeatable functional connectivity test.                                         |

## 15.2 Final submission checklist

- /[x/] Terraform SKE configuration present and validated

- /[x/] Two worker nodes Ready

- /[x/] CloudNativePG Operator and CRDs installed

- /[x/] PostgreSQL Custom Resource applied

- /[x/] Two PostgreSQL instances healthy

- /[x/] Services, PVCs and Secrets inspected

- /[x/] Secure application-style connection demonstrated

- /[x/] Table create/insert/select demonstrated

- /[x/] Primary and standby roles verified

- /[x/] Replication demonstrated

- /[x/] Standby write rejection demonstrated

- /[x/] Automatic failover demonstrated

- /[x/] Before/after evidence saved

- /[ /] Dedicated persistence test and evidence

- /[ /] Backup/restore implementation and restore evidence (production enhancement)

- /[ /] Final secret/state Git hygiene check

# 16. Knowledge-check questions

20. Why is the SKE control plane considered managed, while the PostgreSQL product is still your responsibility?

21. What is the difference between a CRD and the `level3-postgres` Custom Resource?

22. Why does `kubectl apply` of the Cluster CR not directly contain Pod, Service and PVC definitions?

23. Why should applications connect to `level3-postgres-rw` instead of `level3-postgres-1`?

24. Where did the client Pod obtain the `platformuser` password, and why was it not written into the manifest?

25. What does `pg_is_in_recovery()` tell you?

26. Why is an INSERT failure through `level3-postgres-ro` a successful test result?

27. What survives a Pod deletion, and what would be lost if a PVC were deleted under a Delete reclaim policy?

28. How does required Pod anti-affinity improve availability?

29. What is the difference between high availability, persistence and backup?

30. What are RTO and RPO in the context of the failover test?

31. Why must `terraform.tfstate` and `kubeconfig.yaml` stay out of Git?

# References

/[1/] STACKIT Kubernetes Engine documentation: [<u>https://docs.stackit.cloud/products/runtime/kubernetes-engine/</u>](https://docs.stackit.cloud/products/runtime/kubernetes-engine/)

/[2/] STACKIT Terraform Provider user guide: [<u>https://docs.stackit.cloud/developer-tools/stackit-iac/stackit-terraform-provider/</u>](https://docs.stackit.cloud/developer-tools/stackit-iac/stackit-terraform-provider/)

/[3/] Terraform Registry - stackit_ske_cluster resource: [<u>https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs/resources/ske_cluster</u>](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs/resources/ske_cluster)

/[4/] Kubernetes - Operator pattern: [<u>https://kubernetes.io/docs/concepts/extend-kubernetes/operator/</u>](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

/[5/] Kubernetes - Custom Resources: [<u>https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/</u>](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

/[6/] CloudNativePG - Installation and upgrades: [<u>https://cloudnative-pg.io/docs/1.30/installation_upgrade/</u>](https://cloudnative-pg.io/docs/1.30/installation_upgrade/)

/[7/] CloudNativePG - Architecture and application Services: [<u>https://cloudnative-pg.github.io/docs/1.28/architecture/</u>](https://cloudnative-pg.github.io/docs/1.28/architecture/)

/[8/] CloudNativePG - Connecting from an application: [<u>https://cloudnative-pg.io/docs/1.29/applications</u>](https://cloudnative-pg.io/docs/1.29/applications)

/[9/] CloudNativePG - Storage: [<u>https://cloudnative-pg.github.io/docs/1.28/storage/</u>](https://cloudnative-pg.github.io/docs/1.28/storage/)

/[10/] Kubernetes - Persistent Volumes: [<u>https://kubernetes.io/docs/concepts/storage/persistent-volumes/</u>](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

/[11/] Kubernetes - Secrets and good practices: [<u>https://kubernetes.io/docs/concepts/security/secrets-good-practices/</u>](https://kubernetes.io/docs/concepts/security/secrets-good-practices/)

/[12/] Kubernetes - kubectl exec: [<u>https://kubernetes.io/docs/reference/kubectl/generated/kubectl_exec/</u>](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_exec/)

*Official sources were checked on 29 July 2026. Project-specific values and results are taken from the user's verified terminal output, manifests and evidence files.*

# Appendix A - Important observed resources

    # Final post-failover summary (observed)
    Cluster: level3-postgres
    Namespace: database-services
    Instances: 2
    Ready: 2
    Status: Cluster in healthy state
    Primary before: level3-postgres-1
    Primary after:  level3-postgres-2

    Pods:
    - level3-postgres-1   Running
    - level3-postgres-2   Running
    - level3-psql-client  Running

    Services:
    - level3-postgres-r
    - level3-postgres-ro
    - level3-postgres-rw

    Storage:
    - one 5 GiB RWO PVC per PostgreSQL instance
    - StorageClass: premium-perf1-stackit
    - both claims Bound

# Appendix B - Evidence screenshots

<img src="docs/assets/media/image7.png" style="width:6.37795in;height:1.684in" />

Figure B1 - Connectivity and resource verification detail

<img src="docs/assets/media/image8.png" style="width:5.19685in;height:7.04596in" />

Figure B2 - Saved post-failover evidence file

<img src="docs/assets/media/image9.png" style="width:6.37795in;height:5.85566in" />

Figure B3 - CloudNativePG failover events and final terminal state
