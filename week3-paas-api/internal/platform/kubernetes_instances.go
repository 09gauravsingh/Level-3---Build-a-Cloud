// This file now contains only the required product operations:
// Create
// List
// Delete
// Connection data and read secrets from the Kubernetes cluster.



package platform

import (
	"context"
	"strconv"
	"time"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Create creates one CloudNativePG Cluster resource.
func (s *KubeService) Create(
	ctx context.Context,
	request models.CreateInstanceRequest,
) (models.Instance, error) {
	cluster := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "postgresql.cnpg.io/v1",
			"kind":       "Cluster",

			"metadata": map[string]any{
				"name":      request.Name,
				"namespace": s.namespace,
				"labels": map[string]any{
					managedByLabel: managedByValue,
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

	created, err := s.clusters.Create(
		ctx,
		cluster,
		metav1.CreateOptions{},
	)
	if err != nil {
		return models.Instance{}, err
	}

	return toInstance(created), nil
}

// List returns only PostgreSQL instances created by this API.
func (s *KubeService) List(
	ctx context.Context,
) ([]models.Instance, error) {
	list, err := s.clusters.List(
		ctx,
		metav1.ListOptions{
			LabelSelector: managedByLabel + "=" + managedByValue,
		},
	)
	if err != nil {
		return nil, err
	}

	instances := make([]models.Instance, 0, len(list.Items))

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
	// First verify that this instance belongs to our API.
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
) (models.ConnectionData, error) {
	// Prevent access to manually created PostgreSQL clusters.
	if _, err := s.getManagedCluster(ctx, name); err != nil {
		return models.ConnectionData{}, err
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
		return models.ConnectionData{}, err
	}

	port, _ := strconv.Atoi(string(secret.Data["port"]))
	if port == 0 {
		port = 5432
	}

	return models.ConnectionData{
		Host:     string(secret.Data["host"]),
		Port:     port,
		Database: string(secret.Data["dbname"]),
		Username: string(secret.Data["username"]),
		Password: string(secret.Data["password"]),
		URI:      string(secret.Data["uri"]),
	}, nil
}

// getManagedCluster verifies that a cluster was created by our API.
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

	if cluster.GetLabels()[managedByLabel] != managedByValue {
		return nil, apierrors.NewNotFound(
			clusterGVR.GroupResource(),
			name,
		)
	}

	return cluster, nil
}

// toInstance converts a Kubernetes Cluster into simple REST API data.
func toInstance(
	cluster *unstructured.Unstructured,
) models.Instance {
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
	created := cluster.GetCreationTimestamp()

	if !created.Time.IsZero() {
		createdAt = created.Time.UTC().Format(time.RFC3339)
	}

	return models.Instance{
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
