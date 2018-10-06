package main

import (
	"flag"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"strings"
)

func init() {
	if extraFlags := os.Getenv("TF_PROVIDER_KUBECTL_FLAGS"); len(extraFlags) > 0 {
		os.Args = append(os.Args, strings.Split(extraFlags, " ")...)
	}
	flag.Parse()
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}
