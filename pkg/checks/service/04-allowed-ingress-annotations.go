package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"gopkg.in/yaml.v2"
)

var allowedAnnotations = []string{
	"nginx.ingress.kubernetes.io/proxy-connect-timeout",
	"nginx.ingress.kubernetes.io/proxy-read-timeout",
	"nginx.ingress.kubernetes.io/proxy-send-timeout",
	"nginx.ingress.kubernetes.io/configuration-snippet",
}

type IngressAnnotationsCheck struct {
}

func (i IngressAnnotationsCheck) ID() string {
	return "SVC04"
}

func (i IngressAnnotationsCheck) Name() string {
	return "Check Ingress Annotations in values.yaml files"
}

func (i IngressAnnotationsCheck) Run(conf *config.Config) verifier.CheckResult {
	valuesFiles := []string{}
	err := filepath.Walk(conf.ProjectDir+"/deploy/envs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "values.yaml") {
			valuesFiles = append(valuesFiles, path)
		}
		return nil
	})
	if err != nil {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("error walking the path: %v", err),
			WarningOnly: false,
		}
	}

	for _, file := range valuesFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return verifier.CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("error reading file [%s]: %v", file, err),
			}
		}

		// Parse YAML file
		var content map[string]interface{}
		err = yaml.Unmarshal(data, &content)
		if err != nil {
			return verifier.CheckResult{
				Passed:  false,
				Message: fmt.Sprintf("error parsing YAML file [%s]: %v", file, err),
			}
		}

		// Traverse into ingress.annotations
		ingress, ok := content["ingress"].(map[interface{}]interface{})
		if !ok {
			continue
		}

		annotations, ok := ingress["annotations"].(map[interface{}]interface{})
		if !ok {
			continue
		}

		for key := range annotations {
			keyStr, _ := key.(string)
			if !isAllowed(keyStr) {
				return verifier.CheckResult{
					Passed:  false,
					Message: fmt.Sprintf("disallowed annotation [%s] found in file [%s]", keyStr, file),
				}
			}
		}
	}

	return verifier.CheckResult{
		Passed:  true,
		Message: "all ingress annotations are allowed",
	}
}

func isAllowed(annotation string) bool {
	for _, allowed := range allowedAnnotations {
		if annotation == allowed {
			return true
		}
	}
	return false
}
