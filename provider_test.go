package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServerResource_basic(t *testing.T) {
	token := os.Getenv("HOSTMAN_TOKEN")
	if token == "" {
		t.Skip("HOSTMAN_TOKEN must be set for acceptance tests")
	}

	resourceName := "hostman_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if token == "" {
				t.Fatal("HOSTMAN_TOKEN must be set for acceptance tests")
			}
		},
		Providers: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccServerConfig(token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-test-server"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "200"),
				),
			},
		},
	})
}

func TestAccIPResource_basic(t *testing.T) {
	token := os.Getenv("HOSTMAN_TOKEN")
	if token == "" {
		t.Skip("HOSTMAN_TOKEN must be set for acceptance tests")
	}

	resourceName := "hostman_ip.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if token == "" {
				t.Fatal("HOSTMAN_TOKEN must be set for acceptance tests")
			}
		},
		Providers: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccIPConfig(token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "is_ddos_guard", "false"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "ams-1"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tf-test-ip"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
				),
			},
		},
	})
}

// testAccServerConfig returns the Terraform configuration for acceptance testing.
// It includes the provider block (with token) and the hostman_server resource.
func testAccServerConfig(token string) string {
	return fmt.Sprintf(`
provider "hostman" {
  token = %q
}

resource "hostman_server" "test" {
  name      = "tf-test-server"
  bandwidth = 200
}
`, token)
}

func TestAccKubernetesResource_basic(t *testing.T) {
	token := os.Getenv("HOSTMAN_TOKEN")
	if token == "" {
		t.Skip("HOSTMAN_TOKEN must be set for acceptance tests")
	}

	resourceName := "hostman_kubernetes.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if token == "" {
				t.Fatal("HOSTMAN_TOKEN must be set for acceptance tests")
			}
		},
		Providers: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesConfig(token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-test-cluster"),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", "1.28"),
					resource.TestCheckResourceAttr(resourceName, "network_driver", "flannel"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "ams-1"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

// testAccIPConfig returns the Terraform configuration for IP resource acceptance testing.
func testAccIPConfig(token string) string {
	return fmt.Sprintf(`
provider "hostman" {
  token = %q
}

resource "hostman_ip" "test" {
  is_ddos_guard     = false
  availability_zone = "ams-1"
  comment           = "tf-test-ip"
}
`, token)
}

// testAccKubernetesConfig returns the Terraform configuration for Kubernetes resource acceptance testing.
func testAccKubernetesConfig(token string) string {
	return fmt.Sprintf(`
provider "hostman" {
  token = %q
}

resource "hostman_kubernetes" "test" {
  name              = "tf-test-cluster"
  k8s_version       = "1.28"
  network_driver    = "flannel"
  availability_zone = "ams-1"
}
`, token)
}
