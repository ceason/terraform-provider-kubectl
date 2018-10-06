package main

import (
	"flag"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"strings"
	"k8s.io/kubernetes/pkg/kubectl/cmd"
)

func init() {
	if extraFlags := os.Getenv("TF_PROVIDER_KUBECTL_FLAGS"); len(extraFlags) > 0 {
		os.Args = append(os.Args, strings.Split(extraFlags, " ")...)
	}
	flag.Parse()
}

func main() {
	// run the embedded kubectl command if the first arg is 'kubectl'
	if len(os.Args) > 1 && os.Args[1] == "kubectl" {
		command := cmd.NewDefaultKubectlCommand()
		command.SetArgs(os.Args[2:])
		command.Execute()
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: func() terraform.ResourceProvider {
				return Provider()
			},
		})
	}
}
