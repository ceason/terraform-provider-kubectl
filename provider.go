package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type providerConfig struct {
	kubectl *kubectlCli
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"context": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ConfigureFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"kubectl_generic_object": resourceGenericObject(),
		},
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	// todo: validate this context exists in kubectl config
	ctx := d.Get("context").(string)
	cfg := &providerConfig{
		kubectl: &kubectlCli{
			context: ctx,
		},
	}
	return cfg, nil
}
