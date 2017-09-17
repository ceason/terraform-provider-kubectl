package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"os/exec"
	"strings"
	"gopkg.in/yaml.v2"
)

const (
	KUBECTL_APPLY_PRUNELABEL = "terraformKubectlPrunelabel"
)

type kubectlCli struct {
	context string
}

func (cfg *kubectlCli) Apply(yamlObject string, pruneId string) error {
	cmdInput, err := mergePrunelabel(yamlObject, pruneId)
	if err != nil {
		return err
	}
	cmdStr := fmt.Sprintf("kubectl --context=%s apply --prune -l %s=%s -f -", cfg.context, KUBECTL_APPLY_PRUNELABEL, pruneId)
	_, _, err = executeCmd(cmdStr, cmdInput)
	return err
}

// delete all objects in the provided yaml
func (cfg *kubectlCli) Delete(yamlObject string) error {
	cmdStr := fmt.Sprintf("kubectl --context=%s delete -f -", cfg.context)
	_, _, err := executeCmd(cmdStr, yamlObject)
	return err
}

// Get exported objects in yaml format
func (cfg *kubectlCli) Get(yamlObject string) (string, error) {
	cmdStr := fmt.Sprintf("kubectl --context=%s get --export -oyaml -f -", cfg.context)
	stdout, _, err := executeCmd(cmdStr, yamlObject)
	return stdout, err
}

// handy wrapper to execute a CLI command and return the result
func executeCmd(cmdStr string, stdin string) (stdout string, stderr string, err error) {
	args := strings.Split(cmdStr, " ")
	cmd := exec.Command(args[0], args[1:]...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Stdin = strings.NewReader(stdin)
	err = cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	if err != nil {
		err = errors.New(err.Error() + stderr + stdout)
		glog.ErrorDepth(1, err)
	} else {
		glog.V(4).Infof("Exec success: %s\n%s%s", cmdStr, stdout, stderr)
	}
	return
}

// merge the prune label into the provided yaml
func mergePrunelabel(yamlStr string, labelValue string) (string, error) {
	parsed := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(yamlStr), &parsed); err != nil {
		return "", err
	}
	metadataKey, _ := parsed["metadata"]
	metadata := metadataKey.(map[interface{}]interface{})
	if labelsKey, hasLabels := metadata["labels"]; hasLabels {
		labels := labelsKey.(map[interface{}]interface{})
		labels[KUBECTL_APPLY_PRUNELABEL] = labelValue
	} else {
		metadata["labels"] = map[string]string{
			KUBECTL_APPLY_PRUNELABEL: labelValue,
		}
	}
	mergedYaml, err := yaml.Marshal(parsed)
	if err != nil {
		return "", err
	}
	return string(mergedYaml), nil
}
