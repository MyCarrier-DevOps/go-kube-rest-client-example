package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// K8sConfig represents the configuration for a single Kubernetes cluster connection.
// It encapsulates all necessary parameters for authenticating and connecting to a specific cluster,
// including the cluster's name, TLS client configuration, and the API server host address.
type K8sConfig struct {
	// Name is a user-defined identifier for the Kubernetes cluster configuration.
	// This helps in managing configurations for multiple clusters, although currently
	// only a single cluster configuration ("default") is supported by GetK8sConfigs.
	Name string `mapstructure:"name"`

	// Config holds the TLS client configuration required for secure communication
	// with the Kubernetes API server. This includes certificate data and security settings.
	Config TLSClientConfig `mapstructure:"config"`

	// Host is the URL of the Kubernetes API server for the cluster.
	// Example: "https://192.168.1.100:6443"
	Host string `mapstructure:"host"`
}

// TLSClientConfig contains the TLS certificate data required for authenticating
// with a Kubernetes cluster's API server using client certificates.
// All certificate data fields (CertData, KeyData, CAData) are expected to be
// base64 encoded strings.
type TLSClientConfig struct {
	// Insecure determines whether the client should skip TLS verification when
	// connecting to the Kubernetes API server. Setting this to true is generally
	// discouraged in production environments due to security risks, but can be
	// useful for development or testing with self-signed certificates.
	Insecure bool `json:"insecure"`

	// CertData contains the base64 encoded client certificate data. This certificate
	// is used by the client to authenticate itself to the Kubernetes API server.
	CertData string `json:"certData"`

	// KeyData contains the base64 encoded client private key data. This key corresponds
	// to the client certificate provided in CertData.
	KeyData string `json:"keyData"`

	// CAData contains the base64 encoded certificate authority (CA) data. This CA
	// certificate is used by the client to verify the identity of the Kubernetes
	// API server.
	CAData string `json:"caData"`
}

// KubeConfig represents the structure expected within the K8S_CONFIG environment
// variable when it contains JSON formatted configuration. Specifically, it looks
// for a 'tlsClientConfig' key holding the TLS configuration details.
type KubeConfig struct {
	// TLSClientConfig embeds the TLS configuration details (certificates, keys, CA)
	// needed for establishing a secure connection.
	TLSClientConfig TLSClientConfig `json:"tlsClientConfig"`
}

// GetK8sConfigs retrieves Kubernetes cluster configuration from environment variables.
// It expects the TLS client configuration (certificates, keys, CA) to be provided
// as a JSON string within the 'K8S_CONFIG' environment variable, and the API server
// host URL within the 'K8S_HOST' environment variable.
//
// The 'K8S_CONFIG' environment variable should contain a JSON object with a
// 'tlsClientConfig' key, which in turn contains 'insecure', 'certData', 'keyData',
// and 'caData' fields. All certificate data must be base64 encoded.
// Example K8S_CONFIG value:
// '{"tlsClientConfig":{"insecure":false,"certData":"LS0t...","keyData":"LS0t...","caData":"LS0t..."}}'
//
// The 'K8S_HOST' environment variable should contain the full URL of the Kubernetes
// API server.
// Example K8S_HOST value:
// 'https://my-kube-api.example.com:6443'
//
// It returns a K8sConfig struct populated with the retrieved configuration data
// and a default name "default". If either environment variable is missing or if
// the JSON in K8S_CONFIG cannot be unmarshalled, it returns an error.
func GetK8sConfigs() (K8sConfig, error) {
	viper.AutomaticEnv() // Automatically read environment variables

	config := os.Getenv("K8S_CONFIG")

	if config == "" {
		return K8sConfig{}, fmt.Errorf("K8S_CONFIG environment variable is not set")
	}

	var tlsConfig TLSClientConfig
	var kubeConfig KubeConfig

	// Try unmarshalling the JSON configuration from the environment variable
	if err := json.Unmarshal([]byte(config), &kubeConfig); err == nil {
		tlsConfig = kubeConfig.TLSClientConfig
	} else {
		return K8sConfig{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	k8sConfig := K8sConfig{
		Name:   "default",
		Config: tlsConfig,
		Host:   os.Getenv("K8S_HOST"),
	}
	if k8sConfig.Host == "" {
		return K8sConfig{}, fmt.Errorf("K8S_HOST environment variable is not set")
	}

	return k8sConfig, nil
}
