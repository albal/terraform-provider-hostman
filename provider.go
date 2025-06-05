package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOSTMAN_TOKEN", nil),
				Description: "API token for Hostman",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hostman_server": resourceServer(),
			"hostman_ip":     resourceIP(),
		},
		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			token := d.Get("token").(string)
			return token, nil
		},
	}
}


