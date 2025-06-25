package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceServer(t *testing.T) {
	resource := resourceServer()

	// Test that the resource has the correct schema
	expectedFields := []string{"name", "bandwidth", "preset_id", "os_id", "image_id", "is_ddos_guard", "root_pass"}
	for _, field := range expectedFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("expected field %q not found in schema", field)
		}
	}

	// Test required fields
	requiredFields := []string{"name", "bandwidth", "is_ddos_guard"}
	for _, field := range requiredFields {
		if !resource.Schema[field].Required {
			t.Errorf("expected field %q to be required", field)
		}
	}

	// Test computed fields
	computedFields := []string{"root_pass"}
	for _, field := range computedFields {
		if !resource.Schema[field].Computed {
			t.Errorf("expected field %q to be computed", field)
		}
	}

	// Test sensitive fields
	sensitiveFields := []string{"root_pass"}
	for _, field := range sensitiveFields {
		if !resource.Schema[field].Sensitive {
			t.Errorf("expected field %q to be sensitive", field)
		}
	}
}

func TestResourceServerSchema(t *testing.T) {
	resource := resourceServer()

	// Test that we can create a ResourceData with valid values
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"name":          "test-server",
		"bandwidth":     200,
		"is_ddos_guard": false,
		"preset_id":     123,
		"os_id":         99,
	})

	if data.Get("name").(string) != "test-server" {
		t.Errorf("expected name to be 'test-server', got %v", data.Get("name"))
	}

	if data.Get("bandwidth").(int) != 200 {
		t.Errorf("expected bandwidth to be 200, got %v", data.Get("bandwidth"))
	}

	if data.Get("is_ddos_guard").(bool) != false {
		t.Errorf("expected is_ddos_guard to be false, got %v", data.Get("is_ddos_guard"))
	}
}

func TestResourceIP(t *testing.T) {
	resource := resourceIP()

	// Test that the resource has the correct schema
	expectedFields := []string{"is_ddos_guard", "availability_zone", "comment", "ip", "resource_type", "resource_id"}
	for _, field := range expectedFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("expected field %q not found in schema", field)
		}
	}

	// Test computed fields
	computedFields := []string{"ip"}
	for _, field := range computedFields {
		if !resource.Schema[field].Computed {
			t.Errorf("expected field %q to be computed", field)
		}
	}

	// Test default values
	if resource.Schema["is_ddos_guard"].Default != false {
		t.Errorf("expected is_ddos_guard default to be false, got %v", resource.Schema["is_ddos_guard"].Default)
	}

	if resource.Schema["availability_zone"].Default != "ams-1" {
		t.Errorf("expected availability_zone default to be 'ams-1', got %v", resource.Schema["availability_zone"].Default)
	}
}

func TestResourceIPSchema(t *testing.T) {
	resource := resourceIP()

	// Test that we can create a ResourceData with valid values
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"is_ddos_guard":     false,
		"availability_zone": "ams-1",
		"comment":           "test comment",
		"resource_type":     "server",
		"resource_id":       "123",
	})

	if data.Get("is_ddos_guard").(bool) != false {
		t.Errorf("expected is_ddos_guard to be false, got %v", data.Get("is_ddos_guard"))
	}

	if data.Get("availability_zone").(string) != "ams-1" {
		t.Errorf("expected availability_zone to be 'ams-1', got %v", data.Get("availability_zone"))
	}

	if data.Get("comment").(string) != "test comment" {
		t.Errorf("expected comment to be 'test comment', got %v", data.Get("comment"))
	}
}