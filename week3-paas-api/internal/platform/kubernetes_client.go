//This file has only one responsibility: Connect the Go application to Kubernetes.
//file that connects your Go API to the Kubernetes cluster.

package platform

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// This label identifies PostgreSQL clusters created by our API.
	managedByLabel = "platform.level3.io/managed-by"
	managedByValue = "week3-paas-api"
)

// clusterGVR identifies the CloudNativePG Cluster resource.
var clusterGVR = schema.GroupVersionResource{
	Group:    "postgresql.cnpg.io",
	Version:  "v1",
	Resource: "clusters",
}

// KubeService contains the Kubernetes clients used by the API.
type KubeService struct {
	clusters  dynamic.ResourceInterface
	core      kubernetes.Interface
	namespace string
}

// NewKubeService connects the application to Kubernetes.
func NewKubeService(namespace string) (*KubeService, error) {
	config, err := kubeConfig()
	if err != nil {
		return nil, err
	}

	// Dynamic client manages CloudNativePG custom resources.
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Core client reads Kubernetes Secrets.
	coreClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeService{
		clusters: dynamicClient.
			Resource(clusterGVR).
			Namespace(namespace),
		core:      coreClient,
		namespace: namespace,
	}, nil
}

// kubeConfig uses the Pod ServiceAccount inside SKE.
// During local development, it uses the Mac kubeconfig.
func kubeConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return rest.InClusterConfig()
	}

	rules := clientcmd.NewDefaultClientConfigLoadingRules()

	if path := os.Getenv("KUBECONFIG"); path != "" {
		rules.ExplicitPath = path
	}

	return clientcmd.
		NewNonInteractiveDeferredLoadingClientConfig(
			rules,
			&clientcmd.ConfigOverrides{},
		).
		ClientConfig()
}
