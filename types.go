package main

type apiResource struct {
	name       string
	namespaced bool
	kind       string
	apiGroup   string
}

type kubectlObjectBase struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
}

type kubectlObject struct {
	kubectlObjectBase `yaml:",inline"`
	rawYaml string
}

type kubectlObjectConfig struct {
	kubectlObjectBase `yaml:",inline"`
	userProvidedYaml string
	apiResource      *apiResource
}
