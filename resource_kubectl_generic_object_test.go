package main

import (
	"testing"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func Test_leafValues(t *testing.T) {

	inputRaw := mustReadFile(t, "testdata/leaf-values-input.yaml")
	input := make(map[interface{}]interface{})
	err := yaml.Unmarshal(inputRaw, input)
	if err != nil {
		panic(err)
	}
	want := map[string]string{
		"metadata.annotations[prometheus.io/scrape]":      "true",
		"spec.replicas":                                   "4",
		"spec.template.spec.automountServiceAccountToken": "false",
		"spec.template.spec.containers.0.name":            "sidecar",
		"spec.template.spec.containers.0.args.0":          "echo",
		"spec.template.spec.containers.0.args.1":          "Hello World",
		"spec.template.spec.containers.1.name":            "server",
		"spec.template.spec.containers.1.args.0":          "while true; do sleep 30; done",
	}

	got := leafValues("", input)
	// check for extra values
	for key, value := range got {
		_, ok := want[key]
		if !ok {
			t.Errorf("Extra/unwanted/undeclared property '%s=%s'", key, value)
		}
	}

	// check for missing/incorrect values
	for key, value := range want {
		gotVal, ok := got[key]
		if !ok {
			t.Errorf("Missing property '%s=%s'", key, value)
		} else if gotVal != value {
			t.Errorf("Wrong value for '%s=%s' (got %s)", key, value, gotVal)
		}
	}
}

func mustReadFile(t *testing.T, filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not open '%s': %s", filename, err.Error())
	}
	return data
}
