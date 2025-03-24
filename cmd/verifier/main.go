package main

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/executor"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Version - package version ldflags
var Version = "local"

func main() {

	app := &cli.App{
		Name:     "Project Verifier",
		Version:  Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{{
			Name: "Idea authors: RCKHB-Intech / Integrations / CKPR Team",
		},
			{
				Name: "Kazbek Tokaev ksxack@gmail.com",
			},
		},
		Copyright:              "(c) 2025 Kazbek Tokaev",
		HelpName:               "verifier",
		Usage:                  "Checking Gitlab projects for compliance with requirements",
		UsageText:              "verifier [global options] directory",
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Path (`CI_PROJECT_NAMESPACE`) to the project in Gitlab",
				EnvVars:  []string{"CI_PROJECT_NAMESPACE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Identifier (`CI_PROJECT_NAME`) of the project in Gitlab",
				EnvVars:  []string{"CI_PROJECT_NAME"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"c"},
				Usage:    "Title (`CI_PROJECT_TITLE`) of the project in Gitlab",
				EnvVars:  []string{"CI_PROJECT_TITLE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "system",
				Value:    "",
				Usage:    "External system code (`SYSTEM_NAME`) of the project",
				EnvVars:  []string{"SYSTEM_NAME"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "namespace",
				Value:    "",
				Usage:    "Kubernetes namespace name (`NAMESPACE_NAME`) of the project",
				EnvVars:  []string{"NAMESPACE_NAME"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "environment",
				Value:    "",
				Usage:    "Environment name (`ENV_NAME`) of the project",
				EnvVars:  []string{"ENV_NAME"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "log-level",
				Value:    config.DefaultLogLevel,
				Aliases:  []string{"l"},
				Usage:    "Logging level (`LOG_LEVEL`): error / warn / info / debug / trace",
				EnvVars:  []string{"LOG_LEVEL"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "log-format",
				Value:    config.DefaultLogFormat,
				Aliases:  []string{"f"},
				Usage:    "Logging format (`LOG_FORMAT`): text / json / nested",
				EnvVars:  []string{"LOG_FORMAT"},
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "log-timestamp",
				Value:    false,
				Usage:    "Display timestamp in messages: true / false",
				EnvVars:  []string{"LOG_TIMESTAMP"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "tag",
				Aliases:  []string{"g"},
				Usage:    "Set tag in Gitlab",
				EnvVars:  []string{"CI_COMMIT_TAG"},
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "Type (`CHECKS_TYPES`) of checks",
				EnvVars:  []string{"CHECKS_TYPES"},
				Required: false,
			},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}
}

func run(context *cli.Context) error {
	conf := config.CreateConfig(context)

	resultError, resultWarn, err := executor.NewExecutor().Run(conf)
	if err != nil {
		return err
	}

	if resultError {
		log.Fatal("Project verification failed!")
	}

	if resultWarn {
		log.Error("Project verification failed!")
		os.Exit(2)
	}

	log.Info("All checks passed!")

	return nil
}
