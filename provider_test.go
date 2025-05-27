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
