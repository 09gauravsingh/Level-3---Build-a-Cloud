// Previous client.go file created to coonect to the kubernetes and prepare clients. This new file uses those clients to perform real operations.

// Previous code
// → Connect to Kubernetes
// → Prepare clients

// This code
// → Create PostgreSQL
// → List PostgreSQL instances
// → Delete PostgreSQL
// → Read connection credentials
// → Convert Kubernetes data into clean API responses

package main

import (
	"context"
	"strconv"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Create creates one CloudNativePG Cluster resource.
func (s *KubeService) Create(
	ctx context.Context,
	request CreateInstanceRequest,
) (Instance, error) {
	cluster := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "postgresql.cnpg.io/v1",
			"kind":       "Cluster",

			"metadata": map[string]any{
				"name":      request.Name,
				"namespace": s.namespace,

				"labels": map[string]any{
					managedByLabel: "week3-paas-api",
				},
			},

			"spec": map[string]any{
				"instances": request.Instances,

				"storage": map[string]any{
					"size": request.StorageSize,
				},

				"bootstrap": map[string]any{
					"initdb": map[string]any{
						"database": request.Database,
						"owner":    request.Owner,
					},
				},
			},
		},
	}
	//Send create request to Kubernetes
	created, err := s.clusters.Create(
		ctx,
		cluster,
		metav1.CreateOptions{},
	)

	if err != nil {
		return Instance{}, err
	}

	return toInstance(created), nil
}

// List returns only instances created through this API.
func (s *KubeService) List(
	ctx context.Context,
) ([]Instance, error) {
	list, err := s.clusters.List(
		ctx,
		metav1.ListOptions{
			LabelSelector: managedByLabel +
				"=week3-paas-api",
		},
	)

	if err != nil {
		return nil, err
	}

	instances := make(
		[]Instance,
		0,
		len(list.Items),
	)

	for index := range list.Items {
		instances = append(
			instances,
			toInstance(&list.Items[index]),
		)
	}

	return instances, nil
}

// Delete removes one API-managed CloudNativePG Cluster.
func (s *KubeService) Delete(
	ctx context.Context,
	name string,
) error {
	if _, err := s.getManagedCluster(ctx, name); err != nil {
		return err
	}

	return s.clusters.Delete(
		ctx,
		name,
		metav1.DeleteOptions{},
	)
}

// Connection reads credentials from the CloudNativePG application Secret.
func (s *KubeService) Connection(
	ctx context.Context,
	name string,
) (ConnectionData, error) {
	if _, err := s.getManagedCluster(ctx, name); err != nil {
		return ConnectionData{}, err
	}

	secret, err := s.core.
		CoreV1().
		Secrets(s.namespace).
		Get(
			ctx,
			name+"-app",
			metav1.GetOptions{},
		)

	if err != nil {
		return ConnectionData{}, err
	}

	port, _ := strconv.Atoi(
		string(secret.Data["port"]),
	)

	if port == 0 {
		port = 5432
	}

	return ConnectionData{
		Host:     string(secret.Data["host"]),
		Port:     port,
		Database: string(secret.Data["dbname"]),
		Username: string(secret.Data["username"]),
		Password: string(secret.Data["password"]),
		URI:      string(secret.Data["uri"]),
	}, nil
}

// getManagedCluster prevents the API from changing manually created clusters.
func (s *KubeService) getManagedCluster(
	ctx context.Context,
	name string,
) (*unstructured.Unstructured, error) {
	cluster, err := s.clusters.Get(
		ctx,
		name,
		metav1.GetOptions{},
	)

	if err != nil {
		return nil, err
	}

	if cluster.GetLabels()[managedByLabel] != "week3-paas-api" {
		return nil, apierrors.NewNotFound(
			clusterGVR.GroupResource(),
			name,
		)
	}

	return cluster, nil
}

// toInstance converts a Kubernetes resource into simple API JSON.
func toInstance(
	cluster *unstructured.Unstructured,
) Instance {
	desired, _, _ := unstructured.NestedInt64(
		cluster.Object,
		"spec",
		"instances",
	)

	ready, _, _ := unstructured.NestedInt64(
		cluster.Object,
		"status",
		"readyInstances",
	)

	status, _, _ := unstructured.NestedString(
		cluster.Object,
		"status",
		"phase",
	)

	primary, _, _ := unstructured.NestedString(
		cluster.Object,
		"status",
		"currentPrimary",
	)

	storage, _, _ := unstructured.NestedString(
		cluster.Object,
		"spec",
		"storage",
		"size",
	)

	database, _, _ := unstructured.NestedString(
		cluster.Object,
		"spec",
		"bootstrap",
		"initdb",
		"database",
	)

	owner, _, _ := unstructured.NestedString(
		cluster.Object,
		"spec",
		"bootstrap",
		"initdb",
		"owner",
	)

	if status == "" {
		status = "Provisioning"
	}

	createdAt := ""

	// Store the Kubernetes creation timestamp in a variable first.
	created := cluster.GetCreationTimestamp()

	if !created.Time.IsZero() {
		createdAt = created.Time.UTC().Format(time.RFC3339)
	}

	return Instance{
		Name:             cluster.GetName(),
		Namespace:        cluster.GetNamespace(),
		Status:           status,
		DesiredInstances: desired,
		ReadyInstances:   ready,
		Primary:          primary,
		StorageSize:      storage,
		Database:         database,
		Owner:            owner,
		CreatedAt:        createdAt,
	}
}
