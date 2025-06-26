terraform {
  required_providers {
    hostman = {
      source = "albal/hostman"
    }
  }
}

provider "hostman" {
  token = var.hostman_token
}

resource "hostman_kubernetes" "main" {
  name               = var.cluster_name
  k8s_version        = var.kubernetes_version
  network_driver     = var.network_driver
  description        = var.cluster_description
  master_nodes_count = var.master_nodes_count
  preset_id          = var.preset_id
  availability_zone  = var.availability_zone
  is_ingress         = var.is_ingress
  is_k8s_dashboard   = var.is_k8s_dashboard

  worker_groups = [
    {
      name       = "default-workers"
      preset_id  = var.worker_preset_id
      node_count = var.worker_node_count
    }
  ]
}

variable "hostman_token" {
  description = "Hostman API token"
  type        = string
  sensitive   = true
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "my-k8s-cluster"
}

variable "cluster_description" {
  description = "Description of the Kubernetes cluster"
  type        = string
  default     = "Terraform managed Kubernetes cluster"
}

variable "kubernetes_version" {
  description = "Kubernetes version (e.g., v1.28.0+k0s.0)"
  type        = string
  default     = "v1.28.0+k0s.0"
}

variable "network_driver" {
  description = "Network driver for the cluster"
  type        = string
  default     = "kuberouter"
}

variable "master_nodes_count" {
  description = "Number of master nodes in the cluster"
  type        = number
  default     = 1
}

variable "preset_id" {
  description = "Master node tariff ID (e.g., 403)"
  type        = number
  default     = 403
}

variable "worker_preset_id" {
  description = "Worker node tariff ID"
  type        = number
  default     = 1745
}

variable "worker_node_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 3
}

variable "availability_zone" {
  description = "Availability zone for the cluster"
  type        = string
  default     = "ams-1"
}

variable "is_ingress" {
  description = "Enable ingress controller"
  type        = bool
  default     = true
}

variable "is_k8s_dashboard" {
  description = "Enable Kubernetes dashboard"
  type        = bool
  default     = true
}

output "cluster_id" {
  description = "The Kubernetes cluster ID"
  value       = hostman_kubernetes.main.cluster_id
}

output "cluster_endpoint" {
  description = "The Kubernetes cluster API endpoint"
  value       = hostman_kubernetes.main.endpoint
}

output "cluster_status" {
  description = "The current status of the cluster"
  value       = hostman_kubernetes.main.status
}

output "kubeconfig" {
  description = "The kubeconfig for accessing the cluster"
  value       = hostman_kubernetes.main.kubeconfig
  sensitive   = true
}