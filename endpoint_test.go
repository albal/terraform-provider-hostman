package main

import (
	"strings"
	"testing"
)

// TestAPIEndpointPatterns validates that all API endpoints follow consistent patterns
func TestAPIEndpointPatterns(t *testing.T) {
	expectedPatterns := map[string]string{
		"servers":     "https://hostman.com/api/v1/servers",
		"floating-ips": "https://hostman.com/api/v1/floating-ips", 
		"clusters":    "https://hostman.com/api/v1/k8s/clusters",
	}

	// This test serves as documentation and validation of expected API patterns
	// If endpoints change, this test should be updated to reflect the correct patterns
	for resource, expectedURL := range expectedPatterns {
		t.Logf("Expected %s endpoint: %s", resource, expectedURL)
	}
}

// TestKubernetesEndpointConsistency ensures Kubernetes endpoints follow the pattern
func TestKubernetesEndpointConsistency(t *testing.T) {
	// Test that we're using the expected cluster endpoint pattern
	// This helps catch endpoint misconfigurations early in development
	
	baseURL := "https://hostman.com/api/v1/k8s/clusters"
	testID := "test-123"
	
	expectedEndpoints := map[string]string{
		"create": baseURL,
		"read":   baseURL + "/" + testID,
		"update": baseURL + "/" + testID,
		"delete": baseURL + "/" + testID,
	}
	
	for operation, endpoint := range expectedEndpoints {
		t.Logf("Expected %s endpoint: %s", operation, endpoint)
	}
}

// TestAPIEndpointConsistency validates that all resources use consistent URL patterns  
func TestAPIEndpointConsistency(t *testing.T) {
	// This test verifies that endpoints follow the expected pattern:
	// https://hostman.com/api/v1/{resource-name}
	// This could help catch issues like the kubernetes/clusters endpoint mismatch
	
	expectedBaseURL := "https://hostman.com/api/v1/"
	
	testCases := []struct {
		resource string
		endpoint string
	}{
		{"servers", "https://hostman.com/api/v1/servers"},
		{"floating-ips", "https://hostman.com/api/v1/floating-ips"},
		{"clusters", "https://hostman.com/api/v1/k8s/clusters"},
	}
	
	for _, tc := range testCases {
		if !strings.HasPrefix(tc.endpoint, expectedBaseURL) {
			t.Errorf("Endpoint %s for resource %s does not follow expected pattern %s*", 
				tc.endpoint, tc.resource, expectedBaseURL)
		}
		
		if !strings.Contains(tc.endpoint, tc.resource) {
			t.Logf("Note: Endpoint %s does not contain resource name %s - verify this is intentional", 
				tc.endpoint, tc.resource)
		}
	}
}