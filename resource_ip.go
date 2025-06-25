package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPCreate,
		ReadContext:   resourceIPRead,
		UpdateContext: resourceIPUpdate,
		DeleteContext: resourceIPDelete,

		Schema: map[string]*schema.Schema{
			"is_ddos_guard": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ams-1",
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)

	payload := map[string]interface{}{
		"is_ddos_guard":     d.Get("is_ddos_guard").(bool),
		"availability_zone": d.Get("availability_zone").(string),
	}
	if comment := d.Get("comment").(string); comment != "" {
		payload["comment"] = comment
	}

	body, err := makeRequest("POST", "https://hostman.com/api/v1/floating-ips", token, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return diag.FromErr(err)
	}

	ip := resp["ip"].(map[string]interface{})
	id := ""
	switch v := ip["id"].(type) {
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
	d.Set("ip", ip["ip"])

	// Now bind if resource_type and resource_id are set
	resourceType := d.Get("resource_type").(string)
	resourceID := getResourceIDString(d)
	if resourceType != "" && resourceID != "" {
		bindPayload := map[string]interface{}{
			"resource_type": resourceType,
			"resource_id":   resourceID,
		}
		_, err := makeRequest("POST", fmt.Sprintf("https://hostman.com/api/v1/floating-ips/%s/bind", id), token, bindPayload)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIPRead(ctx, d, meta)
}

func resourceIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	body, err := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/floating-ips/%s", id), token, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return diag.FromErr(err)
	}

	ip := resp["ip"].(map[string]interface{})
	d.Set("ip", ip["ip"])
	d.Set("is_ddos_guard", ip["is_ddos_guard"])
	d.Set("availability_zone", ip["availability_zone"])
	d.Set("comment", ip["comment"])
	d.Set("resource_type", ip["resource_type"])
	if ip["resource_id"] != nil {
		d.Set("resource_id", fmt.Sprintf("%v", ip["resource_id"]))
	} else {
		d.Set("resource_id", "")
	}

	return nil
}

func resourceIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	// Only attempt to bind if resource_type or resource_id changed
	if d.HasChange("resource_type") || d.HasChange("resource_id") {
		resourceType := d.Get("resource_type").(string)
		resourceID := getResourceIDString(d)
		if resourceType != "" && resourceID != "" {
			// Read current binding first
			body, err := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/floating-ips/%s", id), token, nil)
			if err != nil {
				return diag.FromErr(err)
			}
			var resp map[string]interface{}
			if err := json.Unmarshal(body, &resp); err != nil {
				return diag.FromErr(err)
			}
			ip := resp["ip"].(map[string]interface{})
			currentType := ip["resource_type"]
			currentID := ""
			if ip["resource_id"] != nil {
				currentID = fmt.Sprintf("%v", ip["resource_id"])
			}
			// Only bind if not already bound to the same resource
			if currentType != resourceType || currentID != resourceID {
				bindPayload := map[string]interface{}{
					"resource_type": resourceType,
					"resource_id":   resourceID,
				}
				_, err := makeRequest("POST", fmt.Sprintf("https://hostman.com/api/v1/floating-ips/%s/bind", id), token, bindPayload)
				if err != nil {
					// If already bound, ignore error
					if !isAlreadyBoundError(err) {
						return diag.FromErr(err)
					}
				}
			}
		}
	}

	return resourceIPRead(ctx, d, meta)
}

// Helper to check for "floating_ip_already_bound" error
func isAlreadyBoundError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "floating_ip_already_bound"))
}

// Helper to safely convert resource_id to string
func getResourceIDString(d *schema.ResourceData) string {
	v := d.Get("resource_id")
	switch id := v.(type) {
	case string:
		return id
	case int:
		return strconv.Itoa(id)
	case float64:
		return fmt.Sprintf("%.0f", id)
	default:
		return fmt.Sprintf("%v", id)
	}
}

func resourceIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	_, err := makeRequest("DELETE", fmt.Sprintf("https://hostman.com/api/v1/floating-ips/%s", id), token, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
