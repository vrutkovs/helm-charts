package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// TestVictoriaMetricsGatewayInstallDefault tests that the victoria-metrics-gateway chart can be installed with default values.
func TestVictoriaMetricsGatewayInstallDefault(t *testing.T) {
	t.Parallel()

	const helmChartPath = "../charts/victoria-metrics-gateway"

	namespaceName := fmt.Sprintf("vmgateway-%s", strings.ToLower(random.UniqueId()))
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	// TODO: needs license
	// options := &helm.Options{
	// 	BuildDependencies: true,
	// 	KubectlOptions:    kubectlOptions,
	// }
}
