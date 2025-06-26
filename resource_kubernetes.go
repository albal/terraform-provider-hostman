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
				Description: "Network driver for the cluster (e.g., kuberouter, flannel, calico, etc.)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Kubernetes cluster",
			},
			"master_nodes_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of master nodes in the cluster",
			},
			"preset_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Master node tariff ID (e.g., 403)",
			},
			"configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Master node configuration parameters. Cannot be provided together with preset_id.",
				ConfigMode:  schema.SchemaConfigModeBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configurator_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Configurator ID",
						},
						"disk": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Disk size in GB",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of CPU cores",
						},
						"ram": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "RAM size in MB",
						},
					},
				},
			},
			"worker_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Worker groups in the cluster",
				ConfigMode:  schema.SchemaConfigModeBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the worker group",
						},
						"preset_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Worker node tariff ID. Cannot be provided together with configuration.",
						},
						"configuration": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Worker node configuration parameters. Cannot be provided together with preset_id.",
							ConfigMode:  schema.SchemaConfigModeBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configurator_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Configurator ID",
									},
									"disk": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Disk size in GB",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of CPU cores",
									},
									"ram": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "RAM size in MB",
									},
								},
							},
						},
						"node_count": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Number of nodes in the group",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 100 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 100, got: %d", key, v))
								}
								return
							},
						},
						"labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Labels for the node group",
							ConfigMode:  schema.SchemaConfigModeBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Label key",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Label value",
									},
								},
							},
						},
						"is_autoscaling": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Autoscaling. Automatic increase and decrease in the number of nodes in the group depending on the current load",
						},
						"min_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Minimum number of nodes. To be used with is_autoscaling and max_size parameters",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 2 {
									errs = append(errs, fmt.Errorf("%q must be >= 2, got: %d", key, v))
								}
								return
							},
						},
						"max_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Maximum number of nodes. To be used with is_autoscaling and min_size parameters",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 2 {
									errs = append(errs, fmt.Errorf("%q must be >= 2, got: %d", key, v))
								}
								return
							},
						},
					},
				},
			},
			"is_ingress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable ingress controller",
			},
			"is_k8s_dashboard": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable Kubernetes dashboard",
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

	// Add optional description
	if description := d.Get("description").(string); description != "" {
		payload["description"] = description
	}

	// Add master nodes count
	if masterNodesCount := d.Get("master_nodes_count").(int); masterNodesCount > 0 {
		payload["master_nodes_count"] = masterNodesCount
	}

	// Add master node preset_id
	if presetId := d.Get("preset_id").(int); presetId > 0 {
		payload["preset_id"] = presetId
	}

	// Add master node configuration (alternative to preset_id)
	if configList := d.Get("configuration").([]interface{}); len(configList) > 0 {
		configMap := configList[0].(map[string]interface{})
		payload["configuration"] = map[string]interface{}{
			"configurator_id": configMap["configurator_id"].(int),
			"disk":           configMap["disk"].(int),
			"cpu":            configMap["cpu"].(int),
			"ram":            configMap["ram"].(int),
		}
	}

	// Add worker groups
	if workerGroupsList := d.Get("worker_groups").([]interface{}); len(workerGroupsList) > 0 {
		workerGroups := make([]map[string]interface{}, 0, len(workerGroupsList))
		
		for _, wg := range workerGroupsList {
			workerGroup := wg.(map[string]interface{})
			group := map[string]interface{}{
				"name":       workerGroup["name"].(string),
				"node_count": workerGroup["node_count"].(int),
			}

			// Add preset_id if provided
			if presetId := workerGroup["preset_id"].(int); presetId > 0 {
				group["preset_id"] = presetId
			}

			// Add configuration if provided (alternative to preset_id)
			if configList := workerGroup["configuration"].([]interface{}); len(configList) > 0 {
				configMap := configList[0].(map[string]interface{})
				group["configuration"] = map[string]interface{}{
					"configurator_id": configMap["configurator_id"].(int),
					"disk":           configMap["disk"].(int),
					"cpu":            configMap["cpu"].(int),
					"ram":            configMap["ram"].(int),
				}
			}

			// Add labels if provided
			if labelsList := workerGroup["labels"].([]interface{}); len(labelsList) > 0 {
				labels := make([]map[string]interface{}, 0, len(labelsList))
				for _, l := range labelsList {
					labelMap := l.(map[string]interface{})
					labels = append(labels, map[string]interface{}{
						"key":   labelMap["key"].(string),
						"value": labelMap["value"].(string),
					})
				}
				group["labels"] = labels
			}

			// Add autoscaling fields if provided
			if isAutoscaling := workerGroup["is_autoscaling"].(bool); isAutoscaling {
				group["is_autoscaling"] = isAutoscaling
				
				if minSize := workerGroup["min_size"].(int); minSize > 0 {
					group["min-size"] = minSize
				}
				
				if maxSize := workerGroup["max_size"].(int); maxSize > 0 {
					group["max-size"] = maxSize
				}
			}

			workerGroups = append(workerGroups, group)
		}
		
		payload["worker_groups"] = workerGroups
	}

	// Add optional flags
	if isIngress := d.Get("is_ingress").(bool); isIngress {
		payload["is_ingress"] = isIngress
	}

	if isK8sDashboard := d.Get("is_k8s_dashboard").(bool); isK8sDashboard {
		payload["is_k8s_dashboard"] = isK8sDashboard
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

		// Check cluster status first
		body, err := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/k8s/clusters/%s", id), token, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(body, &resp); err != nil {
			return diag.FromErr(err)
		}

		cluster := resp["cluster"].(map[string]interface{})
		
		// Check for status first
		status, statusOk := cluster["status"].(string)
		if statusOk && (status == "failed" || status == "error" || status == "deleted") {
			return diag.Errorf("cluster creation failed with status: %s", status)
		}
		
		// Check if cluster is in a ready state (ready or started)
		if statusOk && (status == "ready" || status == "started") {
			// Cluster status indicates it's ready, now we can proceed
			// The kubeconfig will be fetched in the read function
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
	
	// Convert cluster ID to string (it may come as float64 for large integers)
	clusterID := ""
	switch v := cluster["id"].(type) {
	case string:
		clusterID = v
	case float64:
		clusterID = fmt.Sprintf("%.0f", v)
	case int:
		clusterID = strconv.Itoa(v)
	default:
		clusterID = fmt.Sprintf("%v", v)
	}
	d.Set("cluster_id", clusterID)
	
	d.Set("status", cluster["status"])

	if k8sVersion, ok := cluster["k8s_version"].(string); ok {
		d.Set("k8s_version", k8sVersion)
	}

	if networkDriver, ok := cluster["network_driver"].(string); ok {
		d.Set("network_driver", networkDriver)
	}

	if description, ok := cluster["description"].(string); ok {
		d.Set("description", description)
	}

	if masterNodesCount, ok := cluster["master_nodes_count"].(float64); ok {
		d.Set("master_nodes_count", int(masterNodesCount))
	}

	if presetId, ok := cluster["preset_id"].(float64); ok {
		d.Set("preset_id", int(presetId))
	}

	if configuration, ok := cluster["configuration"].(map[string]interface{}); ok {
		configList := []interface{}{
			map[string]interface{}{
				"configurator_id": int(configuration["configurator_id"].(float64)),
				"disk":           int(configuration["disk"].(float64)),
				"cpu":            int(configuration["cpu"].(float64)),
				"ram":            int(configuration["ram"].(float64)),
			},
		}
		d.Set("configuration", configList)
	}

	if workerGroups, ok := cluster["worker_groups"].([]interface{}); ok {
		groups := make([]interface{}, 0, len(workerGroups))
		
		for _, wg := range workerGroups {
			workerGroup := wg.(map[string]interface{})
			group := map[string]interface{}{
				"name":       workerGroup["name"].(string),
				"node_count": int(workerGroup["node_count"].(float64)),
			}

			if presetId, ok := workerGroup["preset_id"].(float64); ok {
				group["preset_id"] = int(presetId)
			}

			if configuration, ok := workerGroup["configuration"].(map[string]interface{}); ok {
				configList := []interface{}{
					map[string]interface{}{
						"configurator_id": int(configuration["configurator_id"].(float64)),
						"disk":           int(configuration["disk"].(float64)),
						"cpu":            int(configuration["cpu"].(float64)),
						"ram":            int(configuration["ram"].(float64)),
					},
				}
				group["configuration"] = configList
			}

			if labels, ok := workerGroup["labels"].([]interface{}); ok {
				labelsList := make([]interface{}, 0, len(labels))
				for _, l := range labels {
					labelMap := l.(map[string]interface{})
					labelsList = append(labelsList, map[string]interface{}{
						"key":   labelMap["key"].(string),
						"value": labelMap["value"].(string),
					})
				}
				group["labels"] = labelsList
			}

			if isAutoscaling, ok := workerGroup["is_autoscaling"].(bool); ok {
				group["is_autoscaling"] = isAutoscaling
			}

			if minSize, ok := workerGroup["min-size"].(float64); ok {
				group["min_size"] = int(minSize)
			}

			if maxSize, ok := workerGroup["max-size"].(float64); ok {
				group["max_size"] = int(maxSize)
			}

			groups = append(groups, group)
		}
		
		d.Set("worker_groups", groups)
	}

	if isIngress, ok := cluster["is_ingress"].(bool); ok {
		d.Set("is_ingress", isIngress)
	}

	if isK8sDashboard, ok := cluster["is_k8s_dashboard"].(bool); ok {
		d.Set("is_k8s_dashboard", isK8sDashboard)
	}

	if availabilityZone, ok := cluster["availability_zone"].(string); ok {
		d.Set("availability_zone", availabilityZone)
	}

	if endpoint, ok := cluster["endpoint"].(string); ok {
		d.Set("endpoint", endpoint)
	}

	// Fetch kubeconfig from dedicated endpoint
	kubeconfigBody, kubeconfigErr := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/k8s/clusters/%s/kubeconfig", id), token, nil)
	if kubeconfigErr == nil {
		var kubeconfigResp map[string]interface{}
		if err := json.Unmarshal(kubeconfigBody, &kubeconfigResp); err == nil {
			if kubeconfig, ok := kubeconfigResp["kubeconfig"].(string); ok {
				d.Set("kubeconfig", kubeconfig)
			}
		}
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
	if d.HasChange("description") {
		changes["description"] = d.Get("description").(string)
	}
	if d.HasChange("master_nodes_count") {
		changes["master_nodes_count"] = d.Get("master_nodes_count").(int)
	}
	if d.HasChange("preset_id") {
		changes["preset_id"] = d.Get("preset_id").(int)
	}
	if d.HasChange("configuration") {
		if configList := d.Get("configuration").([]interface{}); len(configList) > 0 {
			configMap := configList[0].(map[string]interface{})
			changes["configuration"] = map[string]interface{}{
				"configurator_id": configMap["configurator_id"].(int),
				"disk":           configMap["disk"].(int),
				"cpu":            configMap["cpu"].(int),
				"ram":            configMap["ram"].(int),
			}
		}
	}
	if d.HasChange("worker_groups") {
		if workerGroupsList := d.Get("worker_groups").([]interface{}); len(workerGroupsList) > 0 {
			workerGroups := make([]map[string]interface{}, 0, len(workerGroupsList))
			
			for _, wg := range workerGroupsList {
				workerGroup := wg.(map[string]interface{})
				group := map[string]interface{}{
					"name":       workerGroup["name"].(string),
					"node_count": workerGroup["node_count"].(int),
				}

				if presetId := workerGroup["preset_id"].(int); presetId > 0 {
					group["preset_id"] = presetId
				}

				if configList := workerGroup["configuration"].([]interface{}); len(configList) > 0 {
					configMap := configList[0].(map[string]interface{})
					group["configuration"] = map[string]interface{}{
						"configurator_id": configMap["configurator_id"].(int),
						"disk":           configMap["disk"].(int),
						"cpu":            configMap["cpu"].(int),
						"ram":            configMap["ram"].(int),
					}
				}

				if labelsList := workerGroup["labels"].([]interface{}); len(labelsList) > 0 {
					labels := make([]map[string]interface{}, 0, len(labelsList))
					for _, l := range labelsList {
						labelMap := l.(map[string]interface{})
						labels = append(labels, map[string]interface{}{
							"key":   labelMap["key"].(string),
							"value": labelMap["value"].(string),
						})
					}
					group["labels"] = labels
				}

				if isAutoscaling := workerGroup["is_autoscaling"].(bool); isAutoscaling {
					group["is_autoscaling"] = isAutoscaling
					
					if minSize := workerGroup["min_size"].(int); minSize > 0 {
						group["min-size"] = minSize
					}
					
					if maxSize := workerGroup["max_size"].(int); maxSize > 0 {
						group["max-size"] = maxSize
					}
				}

				workerGroups = append(workerGroups, group)
			}
			
			changes["worker_groups"] = workerGroups
		}
	}
	if d.HasChange("is_ingress") {
		changes["is_ingress"] = d.Get("is_ingress").(bool)
	}
	if d.HasChange("is_k8s_dashboard") {
		changes["is_k8s_dashboard"] = d.Get("is_k8s_dashboard").(bool)
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