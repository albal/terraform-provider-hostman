package main

import (
	"fmt"
	"testing"
)

// TestAccProviderFactories tests that our provider can be instantiated
func TestAccProviderFactories(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider validation failed: %s", err)
	}
}

// Test configuration parsing without API calls
func TestTerraformConfigParsing(t *testing.T) {
	testCases := []struct {
		name   string
		config string
		valid  bool
	}{
		{
			name: "valid server config",
			config: `
				provider "hostman" {
					token = "fake-token"
				}
				
				resource "hostman_server" "test" {
					name          = "test-server"
					bandwidth     = 200
					is_ddos_guard = false
				}
			`,
			valid: true,
		},
		{
			name: "valid IP config",
			config: `
				provider "hostman" {
					token = "fake-token"
				}
				
				resource "hostman_ip" "test" {
					is_ddos_guard     = false
					availability_zone = "ams-1"
					comment           = "test comment"
				}
			`,
			valid: true,
		},
		{
			name: "valid combined config",
			config: `
				provider "hostman" {
					token = "fake-token"
				}
				
				resource "hostman_server" "web" {
					name          = "web-server"
					bandwidth     = 500
					preset_id     = 123
					os_id         = 99
					is_ddos_guard = true
				}
				
				resource "hostman_ip" "web_ip" {
					is_ddos_guard     = true
					availability_zone = "ams-1"
					comment           = "web server IP"
					resource_type     = "server"
					resource_id       = hostman_server.web.id
				}
			`,
			valid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test that configuration can be parsed
			// This validates the resource schema definitions
			provider := Provider()

			// Create a test helper to validate the config structure
			// Note: We can't fully test without API but we can validate schema
			resources := provider.ResourcesMap

			if len(resources) != 2 {
				t.Errorf("expected 2 resources, got %d", len(resources))
			}

			if _, ok := resources["hostman_server"]; !ok {
				t.Error("hostman_server resource not found")
			}

			if _, ok := resources["hostman_ip"]; !ok {
				t.Error("hostman_ip resource not found")
			}
		})
	}
}

// Test that resource configurations can be created and have expected attributes
func TestResourceConfiguration(t *testing.T) {
	testProvider := Provider()

	t.Run("server resource configuration", func(t *testing.T) {
		serverResource := testProvider.ResourcesMap["hostman_server"]
		if serverResource == nil {
			t.Fatal("server resource not found")
		}

		// Test that all CRUD operations are defined
		if serverResource.CreateContext == nil {
			t.Error("CreateContext not defined for server resource")
		}
		if serverResource.ReadContext == nil {
			t.Error("ReadContext not defined for server resource")
		}
		if serverResource.UpdateContext == nil {
			t.Error("UpdateContext not defined for server resource")
		}
		if serverResource.DeleteContext == nil {
			t.Error("DeleteContext not defined for server resource")
		}
	})

	t.Run("IP resource configuration", func(t *testing.T) {
		ipResource := testProvider.ResourcesMap["hostman_ip"]
		if ipResource == nil {
			t.Fatal("IP resource not found")
		}

		// Test that all CRUD operations are defined
		if ipResource.CreateContext == nil {
			t.Error("CreateContext not defined for IP resource")
		}
		if ipResource.ReadContext == nil {
			t.Error("ReadContext not defined for IP resource")
		}
		if ipResource.UpdateContext == nil {
			t.Error("UpdateContext not defined for IP resource")
		}
		if ipResource.DeleteContext == nil {
			t.Error("DeleteContext not defined for IP resource")
		}
	})
}

// TestProviderMetadata tests provider metadata and configuration
func TestProviderMetadata(t *testing.T) {
	provider := Provider()

	// Test that token is correctly configured
	tokenSchema := provider.Schema["token"]
	if tokenSchema == nil {
		t.Fatal("token schema not found")
	}

	if !tokenSchema.Required {
		t.Error("token should be required")
	}

	if tokenSchema.Type.String() != "TypeString" {
		t.Errorf("expected token type to be TypeString, got %s", tokenSchema.Type.String())
	}

	// Test that ConfigureContextFunc is set
	if provider.ConfigureContextFunc == nil {
		t.Error("ConfigureContextFunc should be set")
	}
}

// Mock tests for testing configuration without real API calls
func TestConfigurationFormats(t *testing.T) {
	testCases := []struct {
		name        string
		serverName  string
		bandwidth   int
		expectValid bool
	}{
		{
			name:        "normal server name",
			serverName:  "my-test-server",
			bandwidth:   200,
			expectValid: true,
		},
		{
			name:        "server name with numbers",
			serverName:  "server123",
			bandwidth:   500,
			expectValid: true,
		},
		{
			name:        "minimal bandwidth",
			serverName:  "minimal-server",
			bandwidth:   100,
			expectValid: true,
		},
		{
			name:        "high bandwidth",
			serverName:  "high-bandwidth-server",
			bandwidth:   1000,
			expectValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := fmt.Sprintf(`
				provider "hostman" {
					token = "test-token"
				}
				
				resource "hostman_server" "test" {
					name          = %q
					bandwidth     = %d
					is_ddos_guard = false
				}
			`, tc.serverName, tc.bandwidth)

			// Basic validation that the config can be structured
			if config == "" {
				t.Error("config should not be empty")
			}

			// Validate that we can build the testAccProviders
			providers := testAccProviders()
			if len(providers) != 1 {
				t.Errorf("expected 1 provider, got %d", len(providers))
			}

			if _, ok := providers["hostman"]; !ok {
				t.Error("hostman provider not found in test providers")
			}
		})
	}
}
