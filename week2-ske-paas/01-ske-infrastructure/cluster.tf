resource "stackit_ske_cluster" "paas" {
  project_id = var.project_id
  name       = var.cluster_name

  node_pools = [
    {
      name         = var.node_pool_name
      machine_type = var.worker_machine_type

      minimum = var.worker_minimum
      maximum = var.worker_maximum

      availability_zones = var.worker_availability_zones

      volume_type = "storage_premium_perf6"
      volume_size = var.worker_volume_size

      cri                     = "containerd"
      allow_system_components = true

      max_surge       = 1
      max_unavailable = 0

      labels = {
        platform     = "paas"
        environment  = "learning"
        "managed-by" = "terraform"
      }
    }
  ]

  network = {
    control_plane = {
      access_scope = "PUBLIC"
    }
  }

  maintenance = {
    start = "01:00:00Z"
    end   = "03:00:00Z"

    enable_kubernetes_version_updates    = true
    enable_machine_image_version_updates = true
  }
}
