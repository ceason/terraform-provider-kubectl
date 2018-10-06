package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"strconv"
	"gopkg.in/yaml.v2"
	"runtime"
	"io/ioutil"
	"os"
)

func NewKubectlCli(context string, defaultNamespace string) (*kubectlCli, error) {

	stdout, _, err := executeArgs(strings.Split(fmt.Sprintf("kubectl --context=%s api-resources", context), " ")...)
	if err != nil {
		return nil, err
	}
	stdoutLines := strings.Split(stdout, "\n")
	// find the starting index of relevant columns
	apiGroupIdx := strings.Index(stdoutLines[0], "APIGROUP")
	if apiGroupIdx == -1 {
		return nil, errors.New(fmt.Sprintf("Could not find APIGROUP column in 'kubectl api-resources' output: %s", stdout))
	}
	var resources []apiResource
	for _, line := range stdoutLines[1:] {
		if line == "" {
			continue
		}
		runes := []rune(line)
		substr := string(runes[apiGroupIdx:])
		fields := strings.Fields(substr)
		namespaced, err := strconv.ParseBool(fields[len(fields)-2])
		if err != nil {
			return nil, err
		}
		resource := apiResource{
			name:       strings.Fields(substr)[0],
			kind:       fields[len(fields)-1],
			namespaced: namespaced,
		}
		if len(fields) > 2 {
			resource.apiGroup = fields[0]
		}
		resources = append(resources, resource)
	}

	client := &kubectlCli{
		context:          context,
		apiResources:     resources,
		defaultNamespace: defaultNamespace,
	}
	if client.defaultNamespace == "" {
		stdout, _, err := executeArgs("kubectl", "--context="+context, "config", "view", "--minify", "--merge", "-ogo-template={{(index .contexts 0).context.namespace}}")
		if err != nil {
			return nil, err
		}
		if stdout == "<no value>" {
			client.defaultNamespace = "default"
		} else {
			client.defaultNamespace = stdout
		}
	}

	return client, nil
}

type kubectlCli struct {
	context          string
	defaultNamespace string
	apiResources     []apiResource
}

func (k *kubectlCli) Apply(cfg *kubectlObjectConfig, dryRun bool) (*kubectlObject, error) {
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("%s-%s*.yaml", cfg.Metadata.Name, cfg.Kind))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte(cfg.userProvidedYaml)); err != nil {
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}
	stdout, _, err := executeArgs("apply",
		"--namespace", cfg.Metadata.Namespace,
		"--context", k.context,
		"-ojson",
		fmt.Sprintf("--dry-run=%s", strconv.FormatBool(dryRun)),
		"-f", tmpfile.Name())
	if err != nil {
		return nil, err
	}

	obj := &kubectlObject{rawYaml: stdout}
	err = yaml.Unmarshal([]byte(stdout), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// delete all objects in the provided yaml
func (k *kubectlCli) Delete(namespace, kind, name string) error {
	cmdStr := fmt.Sprintf("kubectl --context=%s --namespace=%s delete %s %s", k.context, namespace, kind, name)
	_, _, err := executeArgs(strings.Split(cmdStr, " ")...)
	return err
}

func (k *kubectlCli) NewObjectConfig(yamlStr string) (*kubectlObjectConfig, error) {
	// todo: validation (single object, has name/kind/etc)

	obj := &kubectlObjectConfig{userProvidedYaml: yamlStr}
	yaml.Unmarshal([]byte(yamlStr), obj)
	apiGroup := ""
	parts := strings.Split(obj.ApiVersion, "/")
	if len(parts) > 1 {
		apiGroup = parts[0]
	}
	for _, r := range k.apiResources {
		if r.apiGroup == apiGroup && r.kind == obj.Kind {
			obj.apiResource = &r
			break
		}
	}
	if obj.apiResource == nil {
		return nil, errors.New(fmt.Sprintf("Could not find resource kind '%s' in apiGroup '%s' (yaml='%s')", obj.Kind, apiGroup, yamlStr))
	}
	if obj.Metadata.Namespace == "" && obj.apiResource.namespaced {
		obj.Metadata.Namespace = k.defaultNamespace
	}
	if ! obj.apiResource.namespaced {
		obj.Metadata.Namespace = ""
	}

	return obj, nil
}

func (k *kubectlCli) GetObject(resourceId string) (*kubectlObject, error) {
	parts := strings.Split(resourceId, "/")
	namespace := parts[0]
	kind := parts[1]
	name := parts[2]
	cmdStr := fmt.Sprintf("kubectl --namespace=%s --context=%s get %s %s -ojson", namespace, k.context, kind, name)
	stdout, _, err := executeArgs(strings.Split(cmdStr, " ")...)
	if err != nil {
		return nil, err
	}
	obj := &kubectlObject{rawYaml: stdout}
	err = yaml.Unmarshal([]byte(stdout), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (k *kubectlCli) ObjectExists(resourceId string) (bool, error) {
	//return false, errors.New(resourceId)
	namespace, kind, name, err := resourceIdParts(resourceId)
	if err != nil {
		panic(err)
	}

	stdout, _, err := executeArgs(strings.Split(fmt.Sprintf("kubectl get --context=%s --namespace=%s %s -oname", k.context, namespace, kind), " ")...)
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasSuffix(line, "/"+name) {
			return true, nil
		}
	}
	return false, nil
}

// Return the object's fully qualified 'kind' (ie suffixed with apiGroup if necessary)
func (o kubectlObjectConfig) FullKind() string {
	kind := strings.ToLower(o.Kind)
	if o.apiResource.apiGroup != "" {
		kind = fmt.Sprintf("%s.%s", kind, o.apiResource.apiGroup)
	}
	return kind
}

func (o kubectlObjectConfig) ResourceId() string {
	// ResourceId is '[namespace]/(full)kind/name'
	return fmt.Sprintf("%s/%s/%s", o.Metadata.Namespace, o.FullKind(), o.Metadata.Name)
}

// Flat list of all object properties and their (stringified) values
func (o kubectlObjectConfig) Properties() map[string]string {
	raw := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(o.userProvidedYaml), raw)
	if err != nil {
		// pretty certain this can't happen (as the yaml is already validated)
		panic(err)
	}
	out := leafValues("", raw)
	return out
}

// Flat list of all object properties and their (stringified) values
func (o kubectlObject) Properties() map[string]string {
	raw := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(o.rawYaml), raw)
	if err != nil {
		// pretty certain this can't happen (as the yaml is already validated)
		panic(err)
	}
	out := leafValues("", raw)
	return out
}

// handy wrapper to execute a CLI command and return the result
func executeArgs(args ... string) (stdout string, stderr string, err error) {
	if args[0] != "kubectl" {
		args = append([]string{"kubectl"}, args...)
	}
	cmd := exec.Command(os.Args[0], args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("%s:%d] ", file, line)
		err = errors.New(msg + stderr + stdout)
	}
	return
}
