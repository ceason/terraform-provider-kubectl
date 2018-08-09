package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"os"
)

type providerConfig struct {
	kubectl *kubectlCli
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"context": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					// first check kube_ctx environment var
					ctx, ok := os.LookupEnv("KUBE_CTX")
					if ok {
						return ctx, nil
					}
					// default to current context from kubectl
					stdout, _, err := executeCmd("kubectl config current-context", "")
					return stdout, err
				},
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
