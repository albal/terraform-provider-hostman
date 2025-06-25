package main

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProviderConfigure(t *testing.T) {
	testCases := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "valid token",
			token:       "test-token-123",
			expectError: false,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: false, // Provider accepts empty token, validation happens at resource level
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider := Provider()
			
			// Create resource config
			d := schema.TestResourceDataRaw(t, provider.Schema, map[string]interface{}{
				"token": tc.token,
			})

			// Test configure function
			meta, diags := provider.ConfigureContextFunc(context.Background(), d)

			if tc.expectError && !diags.HasError() {
				t.Fatal("expected error, but got none")
			}
			if !tc.expectError && diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}

			if !tc.expectError {
				if meta.(string) != tc.token {
					t.Fatalf("expected token %q, got %q", tc.token, meta.(string))
				}
			}
		})
	}
}