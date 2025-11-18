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

func TestVictoriaMetricsK8sStackBasic(t *testing.T) {
	const retries = 60
	const sleepBetweenRetries = time.Second
	const helmChartPath = "../charts/victoria-metrics-k8s-stack"

	t.Parallel()

	namespaceName := fmt.Sprintf("vmstack-%s", strings.ToLower(random.UniqueId()))

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    kubectlOptions,
	}

	releaseName := fmt.Sprintf("vm-%s", strings.ToLower(random.UniqueId()))
	defer helm.Delete(t, options, releaseName, true)
	helm.Install(t, options, helmChartPath, releaseName)

	// Check that service
	grafanaName := fmt.Sprintf("%s-grafana", releaseName)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, grafanaName, retries, sleepBetweenRetries)
	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, grafanaName, retries, sleepBetweenRetries)
	k8s.WaitUntilSecretAvailable(t, kubectlOptions, grafanaName, retries, sleepBetweenRetries)
	// TODO: check that secret has necessary keys
	k8s.WaitUntilConfigMapAvailable(t, kubectlOptions, grafanaName, retries, sleepBetweenRetries)
	// TODO: check that configmap has necessary keys
	k8s.WaitUntilConfigMapAvailable(t, kubectlOptions, fmt.Sprintf("%s-grafana-config-dashboards", releaseName), retries, sleepBetweenRetries)

	kubeStateMetricsName := fmt.Sprintf("%s-kube-state-metrics", releaseName)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, kubeStateMetricsName, retries, sleepBetweenRetries)
	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, kubeStateMetricsName, retries, sleepBetweenRetries)

	promNodeExporterName := fmt.Sprintf("%s-prometheus-node-exporter", releaseName)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, promNodeExporterName, retries, sleepBetweenRetries)

	operatorName := fmt.Sprintf("%s-victoria-metrics-operator", releaseName)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, operatorName, retries, sleepBetweenRetries)
	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, operatorName, retries, sleepBetweenRetries)

	// Other services are created by operator
}
