package service

import "os"

func shouldRunCheck() bool {
	branch := os.Getenv("CI_COMMIT_REF_NAME")
	pipelineTag := os.Getenv("CI_COMMIT_TAG")

	if pipelineTag != "" {
		return true
	}

	switch branch {
	case "develop", "stage", "main":
		return true
	default:
		return false
	}
}
