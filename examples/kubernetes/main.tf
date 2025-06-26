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
  name              = var.cluster_name
  node_count        = var.node_count
  version           = var.kubernetes_version
  node_type         = var.node_type
  availability_zone = var.availability_zone
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

variable "node_count" {
  description = "Number of nodes in the cluster"
  type        = number
  default     = 3
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.28"
}

variable "node_type" {
  description = "Node type for cluster nodes"
  type        = string
  default     = "standard"
}

variable "availability_zone" {
  description = "Availability zone for the cluster"
  type        = string
  default     = "ams-1"
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