package verifier

import "MadEngineX/gitlab-project-verifier/config"

type Check interface {
	ID() string
	Name() string
	Run(conf *config.Config) CheckResult
}

type CheckResult struct {
	Passed      bool
	Message     string
	WarningOnly bool
}

func (result CheckResult) Pointer() *CheckResult {
	return &result
}
