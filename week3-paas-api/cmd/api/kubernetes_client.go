package main

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const managedByLabel = "platform.level3.io/managed-by"

// clusterGVR identifies the CloudNativePG Cluster resource.
var clusterGVR = schema.GroupVersionResource{
	Group:    "postgresql.cnpg.io",
	Version:  "v1",
	Resource: "clusters",
}

// KubeService contains the Kubernetes clients used by our API.
type KubeService struct {
	clusters  dynamic.ResourceInterface
	core      kubernetes.Interface
	namespace string
}

// NewKubeService creates clients for CloudNativePG resources and Secrets.
func NewKubeService(namespace string) (*KubeService, error) {
	config, err := kubeConfig()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

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

// kubeConfig uses the Pod ServiceAccount in SKE.
// On the Mac, it uses KUBECONFIG or the default kubeconfig.
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
