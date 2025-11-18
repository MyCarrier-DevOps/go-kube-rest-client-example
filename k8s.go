package main

import (
	"encoding/base64"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// decodeBase64 safely decodes a base64 encoded string.
// It handles empty input strings by returning nil data and nil error.
// If the input string is not empty but fails decoding, it returns an error
// indicating the failure.
//
// Parameters:
//
//	encodedData: The base64 encoded string to decode.
//
// Returns:
//
//	A byte slice containing the decoded data, or nil if the input was empty.
//	An error if decoding fails, otherwise nil.
func decodeBase64(encodedData string) ([]byte, error) {
	if encodedData == "" {
		return nil, nil
	}

	// Use := for declaration and assignment
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return decodedData, nil
}

// CreateExternalClusterKubeRestClient creates a Kubernetes clientset configured to connect
// to a cluster from outside the cluster network (e.g., from a developer machine).
// It uses the provided K8sConfig which contains the API server host URL and
// base64 encoded TLS credentials (client certificate, client key, CA certificate).
//
// This function first decodes the base64 encoded certificate data from the K8sConfig.
// It requires all three data fields (CertData, KeyData, CAData) to be present and valid.
// If any data is missing or fails decoding, it returns an error.
//
// After decoding, it constructs a rest.Config object using the host URL and TLS
// configuration. This config is then used to create a kubernetes.Clientset.
//
// Finally, it performs a test query (fetching the server version) to verify the
// connection to the cluster. If the connection is successful, it prints a success
// message and returns the clientset. If the connection fails, it returns an error.
//
// Parameters:
//
//	k8sconfig: A K8sConfig struct containing the connection details and credentials
//	           for the target Kubernetes cluster.
//
// Returns:
//
//	A pointer to a configured kubernetes.Clientset ready for interacting with the cluster.
//	An error if any step fails (decoding credentials, creating config, creating clientset,
//	or connecting to the cluster).
func CreateExternalClusterKubeRestClient(k8sconfig K8sConfig) (*kubernetes.Clientset, error) {
	var certData, keyData, caData []byte
	var err error

	// Only attempt to decode if data is present
	if k8sconfig.Config.CertData != "" {
		certData, err = decodeBase64(k8sconfig.Config.CertData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode certificate data for cluster %s: %w", k8sconfig.Name, err)
		}
	} else {
		return nil, fmt.Errorf("no certificate data provided for cluster %s", k8sconfig.Name)
	}

	if k8sconfig.Config.KeyData != "" {
		keyData, err = decodeBase64(k8sconfig.Config.KeyData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode key data for cluster %s: %w", k8sconfig.Name, err)
		}
	} else {
		return nil, fmt.Errorf("no key data provided for cluster %s", k8sconfig.Name)
	}

	if k8sconfig.Config.CAData != "" {
		caData, err = decodeBase64(k8sconfig.Config.CAData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode CA data for cluster %s: %w", k8sconfig.Name, err)
		}
	} else {
		return nil, fmt.Errorf("no ca certificate data provided for cluster %s", k8sconfig.Name)
	}

	// Directly create REST config from K8sConfig fields
	restConfig := &rest.Config{
		Host: k8sconfig.Host,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: k8sconfig.Config.Insecure,
			CertData: certData,
			KeyData:  keyData,
			CAData:   caData,
		},
	}

	// Create a Kubernetes clientset using the REST config
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		// Updated error message
		return nil, fmt.Errorf("failed to create Kubernetes clientset for cluster %s: %w", k8sconfig.Name, err)
	}

	// Run a test query to ensure the clientset is working
	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kubernetes cluster %s: %w", k8sconfig.Name, err)
	} else {
		fmt.Printf("Successfully connected to Kubernetes cluster %s\n", k8sconfig.Name)
	}

	return clientset, nil
}

// CreateInClusterKubeRestClient creates a Kubernetes clientset configured to run
// from within a Kubernetes cluster (e.g., inside a pod).
// It automatically uses the service account token and CA certificate mounted
// into the pod by Kubernetes, requiring no explicit configuration parameters.
//
// This function calls rest.InClusterConfig() to load the configuration provided
// by the Kubernetes environment (service account token, API server host/port from
// environment variables, and the cluster's CA certificate).
//
// It then uses this configuration to create a kubernetes.Clientset.
//
// Similar to CreateExternalClusterKubeRestClient, it performs a test query (fetching the
// server version) to verify the connection. If successful, it prints a success
// message and returns the clientset. If any step fails (loading in-cluster config,
// creating clientset, or connecting), it returns an error.
//
// Returns:
//
//	A pointer to a configured kubernetes.Clientset ready for interacting with the cluster.
//	An error if it fails to load the in-cluster configuration, create the clientset,
//	or connect to the cluster API server.
func CreateInClusterKubeRestClient() (*kubernetes.Clientset, error) {
	// Create a Kubernetes client using in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	// Verify the connection to the Kubernetes cluster
	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kubernetes cluster: %w", err)
	} else {
		fmt.Printf("Successfully connected to Kubernetes cluster")
	}

	// Return the clientset
	return clientset, nil
}
