package service

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

const (
	valuesCheckWarningOnly = false // change on true, if this check is not necessary
)

type ValuesCheck struct {
}

func (r ValuesCheck) ID() string {
	return "SVC02"
}

func (r ValuesCheck) Name() string {
	return "Check values.yaml for allowed root keys"
}

func (r ValuesCheck) Run(conf *config.Config) verifier.CheckResult {
	allowedKeys := map[string]bool{
		"deployment":    true,
		"ingress":       true,
		"hostName":      true,
		"domain":        true,
		"monitoring":    true,
		"celery":        true,
		"traefik":       true,
		"image":         true,
		"cronjob":       true,
		"consumer":      true,
		"taskScheduler": true,
		"taskManager":   true,
		"alertRules":    true,
	}

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
			WarningOnly: valuesCheckWarningOnly,
		}
	}

	for _, file := range valuesFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return verifier.CheckResult{
				Passed:      false,
				Message:     fmt.Sprintf("unable to read file [%s]: %v", file, err),
				WarningOnly: valuesCheckWarningOnly,
			}
		}
		var values map[string]interface{}
		if err := yaml.Unmarshal(data, &values); err != nil {
			return verifier.CheckResult{
				Passed:      false,
				Message:     fmt.Sprintf("unable to parse yaml in file [%s]: %v", file, err),
				WarningOnly: valuesCheckWarningOnly,
			}
		}
		for key := range values {
			if !allowedKeys[key] {
				return verifier.CheckResult{
					Passed:      false,
					Message:     fmt.Sprintf("file [%s] contains disallowed key [%s]", file, key),
					WarningOnly: valuesCheckWarningOnly,
				}
			}
		}
	}
	return verifier.CheckResult{
		Passed:  true,
		Message: "all values.yaml files contain only allowed root keys",
	}
}
