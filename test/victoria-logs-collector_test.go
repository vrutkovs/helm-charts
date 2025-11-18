package test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// TestVictoriaLogsCollectorInstallDefault tests that the victoria-logs-collector chart can be installed with default values.
func TestVictoriaLogsCollectorInstallDefault(t *testing.T) {
	t.Parallel()

	const helmChartPath = "../charts/victoria-logs-collector"

	namespaceName := fmt.Sprintf("vlogcollector-%s", strings.ToLower(random.UniqueId()))
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    kubectlOptions,
		SetValues: map[string]string{
			"remoteWrite[0].url": "http://victoria-logs-1:9428",
		},
	}

	// Install the chart and verify no errors occurred.
	releaseName := fmt.Sprintf("vlog-collector-%s", strings.ToLower(random.UniqueId()))
	defer helm.Delete(t, options, releaseName, true)
	helm.Install(t, options, helmChartPath, releaseName)

	k8sClient, err := k8s.GetKubernetesClientFromOptionsE(t, kubectlOptions)
	require.NoError(t, err)

	// Verify the DaemonSet was created and is ready using manual polling
	daemonSetName := fmt.Sprintf("%s-victoria-logs-collector", releaseName)
	var daemonset *appsv1.DaemonSet
	err = wait.PollUntilContextTimeout(context.Background(), pollingInterval, pollingTimeout, true, func(ctx context.Context) (done bool, err error) {
		daemonset, err = k8sClient.AppsV1().DaemonSets(namespaceName).Get(ctx, daemonSetName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return daemonset.Status.CurrentNumberScheduled == daemonset.Status.DesiredNumberScheduled &&
			daemonset.Status.NumberReady == daemonset.Status.DesiredNumberScheduled, nil
	})
	require.NoError(t, err)
	require.NotNil(t, daemonset)
}
