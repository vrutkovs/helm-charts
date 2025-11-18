package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// TestVictoriaMetricsAuthInstallDefault tests that the victoria-metrics-auth chart can be installed with default values.
func TestVictoriaMetricsAuthInstallDefault(t *testing.T) {
	t.Parallel()

	const helmChartPath = "../charts/victoria-metrics-auth"

	namespaceName := fmt.Sprintf("vmauth-%s", strings.ToLower(random.UniqueId()))
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    kubectlOptions,
	}

	// Install the chart and verify no errors occurred.
	releaseName := fmt.Sprintf("vmauth-%s", strings.ToLower(random.UniqueId()))
	defer helm.Delete(t, options, releaseName, true)
	helm.Install(t, options, helmChartPath, releaseName)

	// Verify the Deployment was created and is ready
	vmAuthName := fmt.Sprintf("%s-victoria-metrics-auth", releaseName)
	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, vmAuthName, retries, pollingInterval)

	// Verify the Service was created and is available
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, vmAuthName, retries, pollingInterval)
}
