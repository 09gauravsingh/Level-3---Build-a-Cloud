resource "stackit_ske_kubeconfig" "admin" {
  project_id   = var.project_id
  cluster_name = stackit_ske_cluster.paas.name

  refresh        = true
  expiration     = 604800
  refresh_before = 86400
}

resource "local_sensitive_file" "kubeconfig" {
  content         = stackit_ske_kubeconfig.admin.kube_config
  filename        = "${path.module}/kubeconfig.yaml"
  file_permission = "0600"
}
