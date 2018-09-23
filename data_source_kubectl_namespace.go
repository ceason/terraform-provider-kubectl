package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKubectlNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubectlNamespaceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubectlNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	k := meta.(*providerConfig)
	d.Set("name", k.defaultNamespace)
	d.SetId(k.defaultNamespace)
	return nil
}
