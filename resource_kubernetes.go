package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubernetes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesCreate,
		ReadContext:   resourceKubernetesRead,
		UpdateContext: resourceKubernetesUpdate,
		DeleteContext: resourceKubernetesDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Kubernetes cluster",
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Kubernetes version",
			},
			"network_driver": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network driver for the cluster (e.g., flannel, calico, etc.)",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ams-1",
				Description: "Availability zone for the cluster",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster ID",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cluster API endpoint",
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The kubeconfig for accessing the cluster",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the cluster",
			},
		},
	}
}

func resourceKubernetesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)

	payload := map[string]interface{}{
		"name":           d.Get("name").(string),
		"k8s_version":    d.Get("k8s_version").(string),
		"network_driver": d.Get("network_driver").(string),
	}

	if availabilityZone := d.Get("availability_zone").(string); availabilityZone != "" {
		payload["availability_zone"] = availabilityZone
	}

	// Using /api/v1/k8s/clusters endpoint for Kubernetes cluster operations
	body, err := makeRequest("POST", "https://hostman.com/api/v1/k8s/clusters", token, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return diag.FromErr(err)
	}

	cluster := resp["cluster"].(map[string]interface{})
	id := ""
	switch v := cluster["id"].(type) {
	case string:
		id = v
	case float64:
		id = fmt.Sprintf("%.0f", v)
	case int:
		id = strconv.Itoa(v)
	default:
		id = fmt.Sprintf("%v", v)
	}
	
	d.SetId(id)
	d.Set("cluster_id", id)

	// Wait for cluster to be ready and kubeconfig to be available
	maxWait := 30 * time.Minute
	interval := 10 * time.Second
	start := time.Now()
	
	for {
		if time.Since(start) > maxWait {
			return diag.Errorf("timeout waiting for cluster to become ready")
		}

		// Read current status
		if diags := resourceKubernetesRead(ctx, d, meta); diags.HasError() {
			return diags
		}

		status := d.Get("status").(string)
		if status == "ready" || status == "running" {
			break
		}

		time.Sleep(interval)
	}

	return resourceKubernetesRead(ctx, d, meta)
}

func resourceKubernetesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	body, err := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/k8s/clusters/%s", id), token, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return diag.FromErr(err)
	}

	cluster := resp["cluster"].(map[string]interface{})

	d.Set("name", cluster["name"])
	d.Set("cluster_id", cluster["id"])
	d.Set("status", cluster["status"])

	if k8sVersion, ok := cluster["k8s_version"].(string); ok {
		d.Set("k8s_version", k8sVersion)
	}

	if networkDriver, ok := cluster["network_driver"].(string); ok {
		d.Set("network_driver", networkDriver)
	}

	if availabilityZone, ok := cluster["availability_zone"].(string); ok {
		d.Set("availability_zone", availabilityZone)
	}

	if endpoint, ok := cluster["endpoint"].(string); ok {
		d.Set("endpoint", endpoint)
	}

	if kubeconfig, ok := cluster["kubeconfig"].(string); ok {
		d.Set("kubeconfig", kubeconfig)
	}

	return nil
}

func resourceKubernetesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	changes := make(map[string]interface{})
	if d.HasChange("name") {
		changes["name"] = d.Get("name").(string)
	}
	if d.HasChange("k8s_version") {
		changes["k8s_version"] = d.Get("k8s_version").(string)
	}
	if d.HasChange("network_driver") {
		changes["network_driver"] = d.Get("network_driver").(string)
	}

	if len(changes) > 0 {
		_, err := makeRequest("PUT", fmt.Sprintf("https://hostman.com/api/v1/k8s/clusters/%s", id), token, changes)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceKubernetesRead(ctx, d, meta)
}

func resourceKubernetesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	_, err := makeRequest("DELETE", fmt.Sprintf("https://hostman.com/api/v1/k8s/clusters/%s", id), token, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for deletion to complete
	maxWait := 15 * time.Minute
	interval := 10 * time.Second
	start := time.Now()

	for {
		if time.Since(start) > maxWait {
			return diag.Errorf("timeout waiting for cluster deletion")
		}

		// Check if cluster still exists
		_, err := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/k8s/clusters/%s", id), token, nil)
		if err != nil {
			// If we get an error (likely 404), the cluster is deleted
			break
		}

		time.Sleep(interval)
	}

	d.SetId("")
	return nil
}