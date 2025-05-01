# go-kube-rest-client-example

This repository provides examples of how to create Kubernetes clientset instances in Go using the `k8s.io/client-go` library. It demonstrates two common scenarios:

1.  **In-Cluster Client**: Connecting to the Kubernetes API server from within a pod running inside the cluster.
2.  **External Client**: Connecting to the Kubernetes API server from outside the cluster.

## Intent

The primary goal of this example is to demonstrate how to configure a Kubernetes client using environment variables instead of relying on a `kubeconfig` file. This approach can simplify deployment scenarios, particularly in containerized environments or CI/CD pipelines where managing file mounts and volumes for configuration files can add complexity. By fetching configuration directly from environment variables, the application becomes more self-contained regarding its connection setup.

## Code Overview

*   `main.go`: Contains example usage for both in-cluster and external clients. It attempts to create both clients and list resources (pods for in-cluster, service accounts for external) to demonstrate functionality.
*   `k8s.go`: Defines the functions `CreateInClusterKubeRestClient` and `CreateExternalClusterKubeRestClient` responsible for creating the respective clientsets. It also includes a helper function `decodeBase64`.
*   `config.go`: Defines the configuration structures (`K8sConfig`, `TLSClientConfig`) and the `GetK8sConfigs` function, which reads external cluster configuration from environment variables.

## Client Types

### In-Cluster Client (`CreateInClusterKubeRestClient`)

This client is intended to be used when your application is running inside a Kubernetes pod. It automatically uses the service account token and CA certificate mounted into the pod by Kubernetes. No explicit configuration is needed, provided the service account the application pod runs with has the necessary RBAC permissions.

### External Client (`CreateExternalClusterKubeRestClient`)

This client is used when connecting from outside the Kubernetes cluster. It requires explicit configuration provided via environment variables.

## Configuration for External Client

To use the external client, you need to set the following environment variables:

1.  **`K8S_HOST`**: The full URL of the Kubernetes API server.
    *   Example: `https://<your-cluster-api-server-ip-or-dns>:6443`

2.  **`K8S_CONFIG`**: A JSON string containing the TLS client configuration. This JSON object must have a `tlsClientConfig` key, which holds the necessary certificate data (base64 encoded) and insecurity flag.

    *   **Structure:**
        ```json
        {
          "tlsClientConfig": {
            "insecure": false, // Set to true to skip TLS verification (not recommended for production)
            "certData": "BASE64_ENCODED_CLIENT_CERTIFICATE_DATA",
            "keyData": "BASE64_ENCODED_CLIENT_PRIVATE_KEY_DATA",
            "caData": "BASE64_ENCODED_CA_CERTIFICATE_DATA"
          }
        }
        ```
    *   **Example `K8S_CONFIG` value (replace placeholders with actual base64 data):**
        ```sh
        export K8S_CONFIG='{"tlsClientConfig":{"insecure":false,"certData":"LS0t...<snip>...LS0tLQo=","keyData":"LS0t...<snip>...LS0tLQo=","caData":"LS0t...<snip>...LS0tLQo="}}'
        ```
    *   **How to get certificate data:** You can typically find this data in your `~/.kube/config` file if you have `kubectl` configured to access the cluster. Look for the `cluster` and `user` sections corresponding to your target cluster. The `certificate-authority-data`, `client-certificate-data`, and `client-key-data` fields contain the required base64 encoded strings.

## Running the Example

Ensure you have Go installed.

*   **For External Client:** Set the `K8S_HOST` and `K8S_CONFIG` environment variables as described above.
*   **For In-Cluster Client:** Build a container image and run it as a pod within your Kubernetes cluster. Ensure the pod's service account has permissions to list pods (or other resources you intend to access).

Then, run the code:

```bash
go run .
```

The program will attempt to connect using both methods (if applicable) and print the names of resources it lists or any errors encountered.