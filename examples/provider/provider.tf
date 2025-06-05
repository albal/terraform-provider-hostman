terraform {
  required_providers {
    hostman = {
      source = "albal/hostman"
    }
  }
}

provider "hostman" {
  token = var.token
}

resource "hostman_server" "test-server" {
  name          = "test-server-1"
  bandwidth     = 200
  preset_id     = 3933
  is_ddos_guard = false
  os_id         = 99
}

resource "hostman_ip" "test-ip" {
  is_ddos_guard     = false
  availability_zone = "ams-1"
  comment           = "test-ip"
  resource_type     = "server"
  resource_id       = hostman_server.test-server.id
  depends_on        = [hostman_server.test-server]
}

variable "token" {}

output "server_ip" {
  value = hostman_ip.test-ip.ip
}

output "root_password" {
  value     = hostman_server.test-server.root_pass
  sensitive = true
}
