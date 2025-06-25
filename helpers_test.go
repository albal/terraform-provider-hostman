package main

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestIsAlreadyBoundError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "floating_ip_already_bound error",
			err:      errors.New("API error (400): floating_ip_already_bound"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("API error (500): internal server error"),
			expected: false,
		},
		{
			name:     "partial match",
			err:      errors.New("floating_ip error"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isAlreadyBoundError(tc.err)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetResourceIDString(t *testing.T) {
	testCases := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "string value",
			value:    "test-id",
			expected: "test-id",
		},
		{
			name:     "int value",
			value:    123,
			expected: "123",
		},
		{
			name:     "float64 value",
			value:    456.0,
			expected: "456",
		},
		{
			name:     "other type (bool)",
			value:    true,
			expected: "1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock resource data
			resource := &schema.Resource{
				Schema: map[string]*schema.Schema{
					"resource_id": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			}

			data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
				"resource_id": tc.value,
			})

			result := getResourceIDString(data)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
