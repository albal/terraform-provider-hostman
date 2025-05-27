terraform {
  required_providers {
    hostman = {
      source  = "local/hostman/hostman"
      version = "0.1.0"
    }
  }
}

provider "hostman" {
  token = "EXAMPLE"        # dummy token for schema export
}

