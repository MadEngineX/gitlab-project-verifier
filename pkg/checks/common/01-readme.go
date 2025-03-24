package common

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"fmt"
	"os"
)

const (
	readmeCheckWarningOnly = true // change on false, when this check become necessary
	readmeFilename         = "README.md"
)

type ReadmeCheck struct {
}

func (r ReadmeCheck) ID() string {
	return "CM01"
}

func (r ReadmeCheck) Name() string {
	return "Check README.md file existence"
}

func (r ReadmeCheck) Run(conf *config.Config) verifier.CheckResult {
	if _, err := os.Stat(conf.ProjectDir + "/" + readmeFilename); os.IsNotExist(err) {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("can't find file [%s]", readmeFilename),
			WarningOnly: readmeCheckWarningOnly,
		}
	}
	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("successfully found file [%s]", readmeFilename),
	}
}
