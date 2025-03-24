package service

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"fmt"
	"os"
)

const (
	cmdDirCheckWarningOnly = false // change on true, if this check is not necessary
	deployFolder           = "deploy"
)

type DeployFolderCheck struct {
}

func (r DeployFolderCheck) ID() string {
	return "SVC01"
}

func (r DeployFolderCheck) Name() string {
	return "Check deploy folder existence"
}

func (r DeployFolderCheck) Run(conf *config.Config) verifier.CheckResult {
	if _, err := os.Stat(conf.ProjectDir + "/" + deployFolder); os.IsNotExist(err) {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("unable to find folder [%s]", deployFolder),
			WarningOnly: cmdDirCheckWarningOnly,
		}
	}
	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("successfully found folder [%s]", deployFolder),
	}
}
