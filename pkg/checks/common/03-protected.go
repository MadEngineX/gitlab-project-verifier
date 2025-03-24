package common

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"fmt"
	"os"
	"regexp"
)

const (
	protectedCheckWarningOnly = false // change on true, if this check is not necessary
)

// ProtectedCheck checks if CI_COMMIT_REF_NAME is one of [develop, stage, main, *.*.* tag]
// and ensures CI_COMMIT_REF_PROTECTED is true.
type ProtectedCheck struct {
}

func (r ProtectedCheck) ID() string {
	return "CM03"
}

func (r ProtectedCheck) Name() string {
	return "Check if CI_COMMIT_REF_PROTECTED is true for certain branches and tags"
}

func (r ProtectedCheck) Run(conf *config.Config) verifier.CheckResult {
	refName := os.Getenv("CI_COMMIT_REF_NAME")
	refProtected := os.Getenv("CI_COMMIT_REF_PROTECTED")

	matched, err := regexp.MatchString(`^(develop|stage|main|\d+\.\d+\.\d+)$`, refName)
	if err != nil {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("Error matching ref name: %v", err),
			WarningOnly: protectedCheckWarningOnly,
		}
	}

	if matched && refProtected != "true" {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("Branch/tag [%s] should have CI_COMMIT_REF_PROTECTED set to true", refName),
			WarningOnly: protectedCheckWarningOnly,
		}
	}

	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("Branch/tag [%s] is correctly protected", refName),
	}
}
