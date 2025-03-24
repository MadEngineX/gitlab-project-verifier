package service

import (
	"MadEngineX/gitlab-project-verifier/config"
	"MadEngineX/gitlab-project-verifier/pkg/verifier"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

const (
	namespaceCheckWarningOnly = false // change on true, if this check is not necessary
)

type NamespaceCheck struct {
}

func (r NamespaceCheck) ID() string {
	return "SVC03"
}

func (r NamespaceCheck) Name() string {
	return "Check Kubernetes namespace existence"
}

func (r NamespaceCheck) Run(conf *config.Config) verifier.CheckResult {

	if !shouldRunCheck() {
		return verifier.CheckResult{
			Passed:  true,
			Message: "check skipped",
		}
	}

	var cfg *rest.Config
	var err error

	// Check if running inside a Kubernetes cluster
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		// In-cluster configuration
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return verifier.CheckResult{
				Passed:      false,
				Message:     fmt.Sprintf("unable to create in-cluster config: %s", err.Error()),
				WarningOnly: namespaceCheckWarningOnly,
			}
		}
	} else {
		// Out-of-cluster configuration
		cfg = &rest.Config{
			Host:        conf.KubeApiServer,
			BearerToken: conf.ServiceAccountToken,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("unable to create Kubernetes client: %s", err.Error()),
			WarningOnly: namespaceCheckWarningOnly,
		}
	}

	fmt.Println(conf.NamespaceName)
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), conf.NamespaceName, metav1.GetOptions{})
	if err != nil {
		return verifier.CheckResult{
			Passed:      false,
			Message:     fmt.Sprintf("namespace [%s] does not exist: %s", conf.NamespaceName, err.Error()),
			WarningOnly: namespaceCheckWarningOnly,
		}
	}

	return verifier.CheckResult{
		Passed:  true,
		Message: fmt.Sprintf("namespace [%s] exists", conf.NamespaceName),
	}
}
