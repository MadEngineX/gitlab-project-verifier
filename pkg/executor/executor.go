package executor

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/generated"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sort"
)

type Executor struct {
	verifications map[string]generated.CheckWrapper
	keys          []string
}

func (e *Executor) Run(conf *config.Config) (bool, bool, error) {

	// По-умолчанию считаем проверки пройденными
	resultError := false
	resultWarning := false

	// Обход списка чеков
	for _, checkType := range conf.ChecksTypes {
		for _, key := range e.keys {
			if checkType != e.verifications[key].Group {
				log.Tracef(
					"- Skip [%v], because target check [%+v] was not specified",
					e.verifications[key].Check.Name(),
					checkType,
				)
				continue
			}
			verification := e.verifications[key].Check

			contextLogger := log.WithField("id", verification.ID())

			checkResult := e.executeProjectCheck(verification, conf, contextLogger)
			// If 1 check doesn't pass, all result become False
			if !checkResult.Passed {
				if checkResult.WarningOnly {
					resultWarning = true
				} else {
					resultError = true
				}
			}

		}
	}
	return resultError, resultWarning, nil
}

func makeVerifications() map[string]generated.CheckWrapper {
	checksList := generated.RegisterChecks()
	checkMap := createCheckMap(checksList)
	return checkMap
}

func createCheckMap(checks []generated.CheckWrapper) map[string]generated.CheckWrapper {
	checkMap := make(map[string]generated.CheckWrapper)
	for _, check := range checks {
		if check, exists := checkMap[check.Check.ID()]; exists {
			log.Fatalf("Identifier [%v] of check [%v] already used in check [%v]", check.Check.ID(), check.Check.Name(), checkMap[check.Check.ID()].Check.Name())
		}
		checkMap[check.Check.ID()] = check
	}
	return checkMap
}

func sortMapKeys(checkMap map[string]generated.CheckWrapper) []string {
	var sortedKeys []string

	for key := range checkMap {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)
	return sortedKeys
}

func (e *Executor) executeProjectCheck(verification verifier.Check, conf *config.Config, log *log.Entry) *verifier.CheckResult {
	return e.executeCheck(conf, log, verification)
}

func (e *Executor) executeCheck(
	conf *config.Config,
	log *log.Entry,
	check verifier.Check,
) *verifier.CheckResult {

	// Init and run check
	group := e.verifications[check.ID()].Group
	verificationName := fmt.Sprintf("'%v'", check.Name())

	result := check.Run(conf)

	// Output the result of the verifier checks
	if result.Passed {
		log.Debugf("- Successfully passed: [%v] %v.", group, verificationName)
		return &result
	}

	if result.WarningOnly {
		log.Warnf("- Failed: [%v] %s, Message: %v", group, verificationName, result.Message)
		return &result
	}
	log.Errorf("- Failed: [%v] %s, Message: %v", group, verificationName, result.Message)

	return &result
}

func NewExecutor() *Executor {
	checkMap := makeVerifications()
	return &Executor{
		verifications: checkMap,
		keys:          sortMapKeys(checkMap),
	}
}
