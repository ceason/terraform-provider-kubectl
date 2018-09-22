package main

type apiResource struct {
	name       string
	namespaced bool
	kind       string
	apiGroup   string
}

type kubectlObject struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Annotations struct {
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
		} `json:"annotations"`
	} `json:"metadata"`
	yaml        string
	apiResource *apiResource
}

