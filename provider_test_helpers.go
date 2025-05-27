package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func testAccProviders() map[string]*schema.Provider {
	return map[string]*schema.Provider{
		"hostman": Provider(),
	}
}
