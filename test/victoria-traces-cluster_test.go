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

// TestVictoriaTracesClusterInstallDefault tests that the victoria-traces-cluster chart can be installed with default values.
func TestVictoriaTracesClusterInstallDefault(t *testing.T) {
	t.Parallel()

	const helmChartPath = "../charts/victoria-traces-cluster"

	namespaceName := fmt.Sprintf("vtracecluster-%s", strings.ToLower(random.UniqueId()))
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    kubectlOptions,
	}

	// Install the chart and verify no errors occurred.
	releaseName := fmt.Sprintf("vtcluster-%s", strings.ToLower(random.UniqueId()))
	defer helm.Delete(t, options, releaseName, true)
	helm.Install(t, options, helmChartPath, releaseName)

	k8sClient, err := k8s.GetKubernetesClientFromOptionsE(t, kubectlOptions)
	require.NoError(t, err)

	// Verify trace-ingester StatefulSet was created and is ready using manual polling
	vtStorageName := fmt.Sprintf("%s-vt-cluster-vtstorage", releaseName)
	var vtStorageStatefulSet *appsv1.StatefulSet
	err = wait.PollUntilContextTimeout(context.Background(), pollingInterval, pollingTimeout, true, func(ctx context.Context) (done bool, err error) {
		vtStorageStatefulSet, err = k8sClient.AppsV1().StatefulSets(namespaceName).Get(ctx, vtStorageName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		// Ensure all replicas are ready
		return vtStorageStatefulSet.Status.ReadyReplicas == *vtStorageStatefulSet.Spec.Replicas && *vtStorageStatefulSet.Spec.Replicas > 0, nil
	})
	require.NoError(t, err)
	require.NotNil(t, vtStorageStatefulSet)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, vtStorageName, retries, pollingInterval)

	// Verify vtinsert Deployment was created and is ready using manual polling
	vtInsertName := fmt.Sprintf("%s-vt-cluster-vtinsert", releaseName)
	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, vtInsertName, retries, pollingInterval)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, vtInsertName, retries, pollingInterval)

	// Verify vtselect Deployment was created and is ready using manual polling
	vtSelectName := fmt.Sprintf("%s-vt-cluster-vtselect", releaseName)
	k8s.WaitUntilDeploymentAvailable(t, kubectlOptions, vtSelectName, retries, pollingInterval)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, vtSelectName, retries, pollingInterval)
}
