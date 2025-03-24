package golang

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"fmt"
	"os"
)

const (
	cmdDirCheckWarningOnly = true // change on false, when this check become necessary
	cmdDirname             = "cmd"
)

type CmdDirCheck struct {
}

func (r CmdDirCheck) ID() string {
	return "GO02"
}

func (r CmdDirCheck) Name() string {
	return "Check cmd directory existence"
}

func (r CmdDirCheck) Run(conf *config.Config) verifier.CheckResult {
	if _, err := os.Stat(conf.ProjectDir + "/" + cmdDirname); os.IsNotExist(err) {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("unable to find directory [%s]", cmdDirname),
			WarningOnly: cmdDirCheckWarningOnly,
		}
	}
	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("successfully found directory [%s]", cmdDirname),
	}
}
