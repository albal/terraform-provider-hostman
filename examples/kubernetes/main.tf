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
  k8s_version       = var.kubernetes_version
  network_driver    = var.network_driver
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

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.28"
}

variable "network_driver" {
  description = "Network driver for the cluster"
  type        = string
  default     = "flannel"
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