package config

import (
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	DefaultLogLevel  = "info"
	DefaultLogFormat = "nested"
)

type Config struct {
	ProjectPath   string
	ProjectName   string
	ProjectTitle  string
	ProjectDir    string
	ChecksTypes   []string
	ProjectTag    string
	ProjectSystem string
	NamespaceName string
	EnvName       string
	Mode          string

	KubeApiServer          string `env:"K8S_API_SERVER"`
	ServiceAccountToken    string `env:"K8S_SA_TOKEN"`
	DynamicVerifierAddress string `env:"DYNAMIC_VERIFIER_ADDRESS"`

	ProjectSystemType string
}

func CreateConfig(context *cli.Context) *Config {
	if context.NArg() == 0 {
		log.Fatal("Mandatory argument missing - project directory path")
	}

	logLevelValue := context.String("log-level")
	logLevel, err := log.ParseLevel(logLevelValue)
	if err != nil {
		log.Warnf("Invalid log level specified [%v], default [%v] will be used", logLevelValue, DefaultLogLevel)
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	logFormatValue := strings.ToLower(context.String("log-format"))
	if logFormatValue != "text" && logFormatValue != "json" && logFormatValue != DefaultLogFormat {
		log.Warnf("Invalid log format specified [%v], default [%v] will be used", logFormatValue, DefaultLogFormat)
		logFormatValue = DefaultLogFormat
	}

	// Show timestamp for each log message
	var showTimestamp = context.Bool("log-timestamp")

	if logFormatValue == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else if logFormatValue == "text" {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp:       !showTimestamp,
			FullTimestamp:          showTimestamp,
			DisableLevelTruncation: true,
		})
	} else if logFormatValue == DefaultLogFormat {
		var timestampFormat = ">"
		if showTimestamp {
			timestampFormat = time.StampMilli
		}
		log.SetFormatter(&nested.Formatter{
			HideKeys:        true,
			ShowFullLevel:   true,
			TimestampFormat: timestampFormat,
		})
	}

	var result = &Config{
		ProjectPath:   context.String("path"),
		ProjectName:   context.String("name"),
		ProjectTitle:  context.String("title"),
		ChecksTypes:   context.StringSlice("type"),
		ProjectDir:    context.Args().Get(0),
		ProjectTag:    context.String("tag"),
		NamespaceName: context.String("namespace"),
		EnvName:       context.String("environment"),
	}

	if err := cleanenv.ReadEnv(result); err != nil {
		log.Fatalf("Unable to parse envs: %s", err.Error())
	}

	log.Info(" -------- ")
	log.Infof("  Verifier version: %v", context.App.Version)
	log.Infof("  Project path: %v", result.ProjectPath)
	log.Infof("  Project code: %v", result.ProjectName)
	log.Infof("  Project title: %v", result.ProjectTitle)
	log.Infof("  Check types: %v", result.ChecksTypes)
	log.Infof("  Project tag: %v", result.ProjectTag)
	log.Infof("  Project environment: %v", result.EnvName)
	log.Infof("  Project namespace name: %v", result.NamespaceName)
	log.Infof("  Project system: %v [%v]", result.ProjectSystem, result.ProjectSystemType)
	log.Info(" -------- ")

	return result
}
