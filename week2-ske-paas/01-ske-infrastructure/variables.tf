variable "project_id" {
  description = "STACKIT project in which the SKE cluster is created."
  type        = string
  default     = "20a60c06-1e1a-406f-a840-37d1ff14f0e8"
}

variable "region" {
  description = "STACKIT region used for the SKE cluster."
  type        = string
  default     = "eu01"
}

variable "cluster_name" {
  description = "Name of the STACKIT Kubernetes Engine PaaS platform cluster."
  type        = string
  default     = "lvl3-paas"
}

variable "node_pool_name" {
  description = "Name of the multi-zone SKE worker-node pool."
  type        = string
  default     = "paas-workers"
}

variable "worker_machine_type" {
  description = "STACKIT machine type used for each SKE worker."
  type        = string
  default     = "g3i.2"
}

variable "worker_minimum" {
  description = "Minimum number of worker nodes maintained by SKE."
  type        = number
  default     = 2
}

variable "worker_maximum" {
  description = "Maximum number of worker nodes allowed by the cluster autoscaler."
  type        = number
  default     = 2
}

variable "worker_availability_zones" {
  description = "Availability zones across which SKE workers are distributed."
  type        = list(string)

  default = [
    "eu01-1"
  ]
}

variable "worker_volume_size" {
  description = "Root disk size in GB for every SKE worker node."
  type        = number
  default     = 32
}