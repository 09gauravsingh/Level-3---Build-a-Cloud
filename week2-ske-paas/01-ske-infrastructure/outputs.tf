output "cluster_name" {
  description = "Name of the SKE PaaS platform cluster."
  value       = stackit_ske_cluster.paas.name
}

output "kubernetes_version" {
  description = "Kubernetes version currently used by SKE."
  value       = stackit_ske_cluster.paas.kubernetes_version_used
}

output "worker_node_pool" {
  description = "Name of the main PaaS worker-node pool."
  value       = var.node_pool_name
}

output "kubeconfig_file" {
  description = "Path to the generated administrator kubeconfig."
  value       = local_sensitive_file.kubeconfig.filename
}

output "kubectl_test_command" {
  description = "Command for verifying the SKE worker nodes."
  value       = "KUBECONFIG=${local_sensitive_file.kubeconfig.filename} kubectl get nodes -o wide"
}
