package common

import (
	"fmt"
	"os"
	"strings"

	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
)

const (
	dockerignoreCheckWarningOnly = true // change on false, when this check become necessary
	dockerfileFilename           = "Dockerfile"
	dockerignoreFilename         = ".dockerignore"
	dockerignoreGitIgnore        = ".git"
)

type DockerignoreCheck struct {
}

func (d DockerignoreCheck) ID() string {
	return "CM02"
}

func (d DockerignoreCheck) Name() string {
	return "Check .dockerignore file existence and .git ignored"
}

func (d DockerignoreCheck) Run(conf *config.Config) verifier.CheckResult {
	dockerfilePath := conf.ProjectDir + "/" + dockerfileFilename

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return verifier.CheckResult{
			Passed:  true,
			Message: fmt.Sprintf("can't find file [%s]", dockerfileFilename),
		}
	}

	dockerignorePath := conf.ProjectDir + "/" + dockerignoreFilename
	if _, err := os.Stat(dockerignorePath); os.IsNotExist(err) {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("[%s] found, but can't find file [%s]", dockerfileFilename, dockerignoreFilename),
			WarningOnly: dockerignoreCheckWarningOnly,
		}
	}

	content, err := os.ReadFile(dockerignorePath)
	if err != nil {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("[%s] found, [%s] found, but can't read [%s] due to the error: %s", dockerfileFilename, dockerignoreFilename, dockerignoreFilename, err.Error()),
			WarningOnly: dockerignoreCheckWarningOnly,
		}
	}

	if !strings.Contains(string(content), dockerignoreGitIgnore) {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("[%s] found, [%s] found, but [%s] is not ingored", dockerfileFilename, dockerignoreFilename, dockerignoreGitIgnore),
			WarningOnly: dockerignoreCheckWarningOnly,
		}
	}

	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("[%s] found, [%s] found, [%s] if ignored", dockerfileFilename, dockerignoreFilename, dockerignoreGitIgnore),
	}
}
