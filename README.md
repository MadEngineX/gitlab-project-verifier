# Verifier 

Checking **Gitlab** projects for compliance with requirements. You can write any check with Golang for your Gitlab project and Verifier will help you to organize your checks.

Inspired by [RSKHB-Intech](https://rshb-intech.ru/rshb-intech)

✨ Key Features
- AI-Assisted Development - Generate checks using natural language prompts
- Multi-Standard Validation - Kubernetes, Vault, and language-specific rules, you can write any type of checks
- CI/CD Native - Designed for seamless pipeline integration
- Extensible Architecture - Add custom checks in minutes

## Contents
- [Usage](#Usage)
- [Local Run](#Local Run)
- [Develop your Checks](#develop-your-checks)
  - [Write your Checks](#write-your-checks)
  - [🛠️ AI-Powered Check Development](#AI-Powered Check Development)
  - [Disabling Check](#disabling-checks)
  
## Usage

Set environment variables (optional):

| Name           | Type     | Description                          |
|----------------|----------|--------------------------------------|
| K8S_API_SERVER | string   | URL of k8s api-server                |
| K8S_SA_TOKEN   | string   | SA token to make requests to K8s API |
| VAULT_ADDR     | string   | URL of Hashicorp Vault               |
| VAULT_TOKEN_RO | string   | Read-Only Token for Hashicorp Vault  |


```shell
verifier [global options] directory

GLOBAL OPTIONS:
   --path PATH, -p PATH                           Path (PATH) to the project in Gitlab [$CI_PROJECT_NAMESPACE]
   --name NAME, -n NAME                           Identifier (NAME) of the project in Gitlab [$CI_PROJECT_NAME]
   --title TITLE, -c TITLE                        Title (TITLE) of the project in Gitlab [$CI_PROJECT_TITLE]
   --system SYSTEM_NAME                           External system code (SYSTEM_NAME) of the project [$SYSTEM_NAME]
   --namespace NAMESPACE_NAME                     Kubernetes namespace name (NAMESPACE_NAME) of the project [$NAMESPACE_NAME]
   --environment ENV_NAME                         Environment name (ENV_NAME) of the project [$ENV_NAME]
   --log-level LEVEL, -l LEVEL                    Logging level (LEVEL): error / warn / info / debug / trace (default: "info") [$LOG_LEVEL]
   --log-format FORMAT, -f FORMAT                 Logging format (FORMAT): text / json / nested (default: "nested") [$LOG_FORMAT]
   --log-timestamp                                Display timestamp in messages: true / false (default: false) [$LOG_TIMESTAMP]
   --tag value, -g value                          Set tag in Gitlab [$CI_COMMIT_TAG]
   --type TYPE, -t TYPE [ --type TYPE, -t TYPE ]  Type (TYPE) of checks [$CHECKS_TYPE]
   --help, -h                                     show help
   --version, -v                                  print the version
```
### Example

```shell
git clone https://gitlab.com/ksxack/weather-bot.git
 cd weather-bot/

verifier --path ksxack/weather-bot --title weather-bot --name "Weather Bot" --environment dev  --log-level debug --type common,golang,service .
## Where --path is $CI_PROJECT_PATH,  --title is $CI_PROJECT_NAME, etc
```
![1](/assets/verifier-01.png)

Types of checks `--type` depends on the type of project, for example you can write group of Python project checks and just add this group in `--type` flag
```shell
verifier --path ksxack/python-tg-bot --title python-tg-bot --name "Telegram Bot" --environment dev  --log-level debug --type common,python,service .
```

## Local Run

If you need K8s integration:
```shell
kns infra-gitlab-runners
kubectl create token gitlab-runner > sa
export K8S_SA_TOKEN=$(cat sa)
export K8S_API_SERVER="https://104.197.203.229:443"

kubectl -n ci-cd-dev port-forward svc/dynamic-verifier-service 8080:8080
export DYNAMIC_VERIFIER_ADDRESS="http://localhost:8080/"
```

Set Gitlab CI variables to enable dynamic Checks (there are disabled fore branches like feature/*)

```shell
## For branches
export CI_COMMIT_REF_NAME="develop"
```

```shell
## For tags
export CI_COMMIT_TAG="1.0.0"
```

```shell
go run cmd/verifier/main.go --path ksxack/weather-bot --title weather-bot --name "Weather Bot" --environment dev  --log-level debug --type common,golang,service  ../weather-bot/
```

Local build 
```shell
go build -o verifier cmd/verifier/main.go
```

## Develop your Checks

### Write your Checks

All Checks are grouped and located in the `pkg/checks` folder. Creating a new check is very simple:

1) If you want to create a new Check Group, create a folder inside `pkg/checks`, for example, `pkg/checks/java`.

2) Inside this folder, create a `.go` file containing your Check. Each Check should be in its own `.go` file. For example, a check for pom.xml can be in a file named pom.go.

3) Your file (pom.go) should be in the package java (package name = folder name).

4) In the pom.go file, you need to define the Check. The Check must satisfy the following interface:

```go
type Check interface {
    ID() string
    Name() string
    Run(conf *config.Config) CheckResult
}
```
To do this, create an empty struct, for example, `PomCheck`, and define the following three methods for it:

- _ID()_ - returns the Check ID.
  ```go
  func (r PomCheck) ID() string {
      return "JV01"
  }
  ```
- _Name()_ - returns the Check Name.
  ```go
  func (r PomCheck) Name() string {
        return "Check pom.xml file existence"
  }
  ```
- _Run()_ - contains your logic. To fail the check, _Run()_ should return `false` in `CheckResult.Passed`. It is also strongly recommended to provide a `Message` to help users understand what went wrong.

5) After creating the Group and defining the Checks, run the code generation:
```bash
go run generate_checks.go
```
The generator will add the Group and Checks to pkg/generated.
Now you can run the Verifier. To call Checks from your Group, the Group name must be listed in the type flag (`--type common,service,java`).

### AI-Powered Check Development

**Checks written by ChatGPT**

Let's say I need to create a new Check Group called system. This group will contain checks for "System configurations" projects.

1) Create a folder named system inside the pkg/checks directory.

2) I want to add a check that ensures "The project path in GitLab (CI_PROJECT_NAME) must be 'system-configurations'."

3) In the pkg/checks/system folder, create a file named path.go.

4) Write a request for ChatGPT:
```
Please create a structure `PathCheck` that satisfies the following interface:

type Check interface {
    ID() string
    Name() string
    Run(conf *config.Config) CheckResult
}

The `Run` method should check that the project name in GitLab (conf.ProjectTitle) is 'system-configurations'.

Example:

package common

import (
	"fmt"
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"os"
)

const (
	readmeFilename = "README.md"
)

type ReadmeCheck struct {
}

func (r ReadmeCheck) ID() string {
	return "CM01"
}

func (r ReadmeCheck) Name() string {
	return "Check README.md file existence"
}

func (r ReadmeCheck) Run(conf *config.Config) verifier.CheckResult {
	if _, err := os.Stat(conf.ProjectDir + "/" + readmeFilename); os.IsNotExist(err) {
		return verifier.CheckResult{
			Passed:  false,
			Message: fmt.Sprintf("can't find file [%s]", readmeFilename),
		}
	}
	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("successfully found file [%s]", readmeFilename),
	}
}
```
5) Ensure everything is generated correctly, and update the messages to make them clear for users.

6) Run the code generator:
```bash
go run generate_checks.go
```
The generator will add the Group and Checks to `pkg/generated`.
Now you can run the Verifier. To call Checks from your Group, the Group name must be listed in the type flag (`--type common,system`).

### Disabling Checks 

In case to disable check, just add `_` in the beginning of Checks filename, like `_04-vault-secrets.go`.

Then re-run generator: 
```bash
go run generate_checks.go
```