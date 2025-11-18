package test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestVictoriaMetricsSingleBasic(t *testing.T) {
	const retries = 60
	const sleepBetweenRetries = time.Second
	const helmChartPath = "../charts/victoria-metrics-single"

	t.Parallel()

	namespaceName := fmt.Sprintf("vmsingle-%s", strings.ToLower(random.UniqueId()))

	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	options := &helm.Options{
		BuildDependencies: true,
		KubectlOptions:    kubectlOptions,
	}

	releaseName := fmt.Sprintf("vmsingle-%s", strings.ToLower(random.UniqueId()))
	defer helm.Delete(t, options, releaseName, true)
	helm.Install(t, options, helmChartPath, releaseName)

	// Check that service was created
	serviceName := fmt.Sprintf("%s-victoria-metrics-single-server", releaseName)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, serviceName, retries, sleepBetweenRetries)

	// No statefulset support in terratest yet, so k8sClient is used
	k8sClient, err := k8s.GetKubernetesClientFromOptionsE(t, kubectlOptions)
	assert.NoError(t, err)

	var sts *appsv1.StatefulSet
	stsName := fmt.Sprintf("%s-victoria-metrics-single-server", releaseName)
	err = wait.PollUntilContextTimeout(context.Background(), time.Second, 10*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		sts, err = k8sClient.AppsV1().StatefulSets(namespaceName).Get(ctx, stsName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return sts.Status.ReadyReplicas > 0, nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, int32(1), sts.Status.Replicas)
	assert.Equal(t, int32(1), sts.Status.ReadyReplicas)
}
