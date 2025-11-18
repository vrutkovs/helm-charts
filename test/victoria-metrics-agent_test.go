package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// TestVictoriaMetricsAgentInstallDefault tests that the victoria-metrics-agent chart can be installed with default values.
func TestVictoriaMetricsAgentInstallDefault(t *testing.T) {
	t.Parallel()

	const retries = 60
	const sleepBetweenRetries = time.Second
	const helmChartPath = "../charts/victoria-metrics-agent"

	namespaceName := fmt.Sprintf("vmagent-%s", strings.ToLower(random.UniqueId()))
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    kubectlOptions,
		SetValues: map[string]string{
			"remoteWrite[0].url": "http://example.com:9428",
		},
	}

	// Install the chart and verify no errors occurred.
	releaseName := "victoria-metrics-agent"
	helm.Install(t, options, helmChartPath, releaseName)

	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, "victoria-metrics-agent", retries, sleepBetweenRetries)
}
