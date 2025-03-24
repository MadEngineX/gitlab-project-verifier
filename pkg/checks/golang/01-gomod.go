package golang

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"fmt"
	"os"
)

const (
	gomodCheckWarningOnly = true // change on false, when this check become necessary
	gomodFilename         = "go.mod"
)

type GomodCheck struct {
}

func (r GomodCheck) ID() string {
	return "GO01"
}

func (r GomodCheck) Name() string {
	return "Check go.mod file existence"
}

func (r GomodCheck) Run(conf *config.Config) verifier.CheckResult {
	if _, err := os.Stat(conf.ProjectDir + "/" + gomodFilename); os.IsNotExist(err) {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("unable to find file [%s]", gomodFilename),
			WarningOnly: gomodCheckWarningOnly,
		}
	}
	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("successfully found file [%s]", gomodFilename),
	}
}
