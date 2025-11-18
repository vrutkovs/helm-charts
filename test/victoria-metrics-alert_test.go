package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// TestVictoriaMetricsAlertInstallDefault tests that the victoria-metrics-alert chart can be installed with default values.
func TestVictoriaMetricsAlertInstallDefault(t *testing.T) {
	t.Parallel()

	const helmChartPath = "../charts/victoria-metrics-alert"

	namespaceName := fmt.Sprintf("vmalert-%s", strings.ToLower(random.UniqueId()))
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    kubectlOptions,
		SetValues: map[string]string{
			"server.datasource.url": "http://example.com",
		},
	}

	// Install the chart and verify no errors occurred.
	releaseName := fmt.Sprintf("vmalert-%s", strings.ToLower(random.UniqueId()))
	defer helm.Delete(t, options, releaseName, true)
	helm.Install(t, options, helmChartPath, releaseName)

	// Verify the Deployment was created and is ready using manual polling
	vmAlertName := fmt.Sprintf("%s-victoria-metrics-alert-server", releaseName)
	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, vmAlertName, retries, pollingInterval)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, vmAlertName, retries, pollingInterval)
}
