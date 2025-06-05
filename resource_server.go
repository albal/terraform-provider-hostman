package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"preset_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"os_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_ddos_guard": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"root_pass": {
				Type:     schema.TypeString,
				Computed: true,
				Sensitive: true,
				Description: "The root password for the server. Only available after creation.",
			},
		},
	}
}

func makeRequest(method, url, token string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)

	payload := map[string]interface{}{
		"name":          d.Get("name").(string),
		"bandwidth":     d.Get("bandwidth").(int),
		"is_ddos_guard": d.Get("is_ddos_guard").(bool),
	}

	if imageID := d.Get("image_id").(string); imageID != "" {
		payload["image_id"] = imageID
	} else if osID := d.Get("os_id").(int); osID != 0 {
		payload["os_id"] = osID
	}

	if presetID := d.Get("preset_id").(int); presetID != 0 {
		payload["preset_id"] = presetID
	}

	body, err := makeRequest("POST", "https://hostman.com/api/v1/servers", token, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return diag.FromErr(err)
	}

	server := resp["server"].(map[string]interface{})
	id := int(server["id"].(float64))
	d.SetId(strconv.Itoa(id))

    // Poll for root_pass to become available
    var rootPass string
    maxWait := 30 * time.Minute
    interval := 5 * time.Second
    start := time.Now()
    for {
        // Fetch server details
        body, err := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/servers/%d", id), token, nil)
        if err != nil {
            return diag.FromErr(err)
        }
        var pollResp map[string]interface{}
        if err := json.Unmarshal(body, &pollResp); err != nil {
            return diag.FromErr(err)
        }
        srv := pollResp["server"].(map[string]interface{})
        if pass, ok := srv["root_pass"].(string); ok && pass != "" {
            rootPass = pass
            break
        }
        if time.Since(start) > maxWait {
            return diag.Errorf("timeout waiting for root_pass to become available")
        }
        time.Sleep(interval)
    }

	d.Set("root_pass", rootPass)

	return resourceServerRead(ctx, d, meta)
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	body, err := makeRequest("GET", fmt.Sprintf("https://hostman.com/api/v1/servers/%s", id), token, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return diag.FromErr(err)
	}

	server := resp["server"].(map[string]interface{})
	d.Set("name", server["name"])
	d.Set("root_pass", server["root_pass"])
	// Add more attributes as needed

	return nil
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	changes := make(map[string]interface{})
	if d.HasChange("name") {
		changes["name"] = d.Get("name").(string)
	}
	if d.HasChange("bandwidth") {
		changes["bandwidth"] = d.Get("bandwidth").(int)
	}
	if d.HasChange("preset_id") {
		changes["preset_id"] = d.Get("preset_id").(int)
	}
	if d.HasChange("os_id") {
		changes["os_id"] = d.Get("os_id").(int)
	}
	if d.HasChange("image_id") {
		changes["image_id"] = d.Get("image_id").(string)
	}
	if d.HasChange("is_ddos_guard") {
		changes["is_ddos_guard"] = d.Get("is_ddos_guard").(bool)
	}

	if len(changes) > 0 {
		_, err := makeRequest("PATCH", fmt.Sprintf("https://hostman.com/api/v1/servers/%s", id), token, changes)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceServerRead(ctx, d, meta)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token := meta.(string)
	id := d.Id()

	_, err := makeRequest("DELETE", fmt.Sprintf("https://hostman.com/api/v1/servers/%s", id), token, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
