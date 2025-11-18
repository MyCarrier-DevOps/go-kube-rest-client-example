package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// example usage of CreateInClusterKubeRestClient
	// no need to pass configuration as it uses the service account token automatically
	// mounted in the pod by kubernetes. make sure this service account has
	// the necessary permissions to access the resources you want to query.
	inClusterClientSet, err := CreateInClusterKubeRestClient()
	if err != nil {
		panic(err)
	}

	pods, err := inClusterClientSet.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		println("In-Cluster Pod Name:", pod.Name)
	}

	// example usage of CreateExternalClusterKubeRestClient
	// get external Kubernetes configuration from environment variables
	k8sConfig, err := GetK8sConfigs()
	if err != nil {
		panic(err)
	}

	externalClusterClientSet, err := CreateExternalClusterKubeRestClient(k8sConfig)
	if err != nil {
		panic(err)
	}

	// example usage of the clientset to list service accounts in the "default" namespace and print their names
	serviceAccounts, err := externalClusterClientSet.CoreV1().
		ServiceAccounts("default").
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, sa := range serviceAccounts.Items {
		println("External Cluster Service Account Name:", sa.Name)
	}
}
