# Kubernetes Admission Webhook Controller

This project implements both validating and mutating admission webhooks for Kubernetes using Go. The webhooks can intercept and optionally modify resource requests to the Kubernetes API server before persistence of the object.

## Project Structure

```
.
├── .dockerignore           # Docker build exclusions
├── certs/                  # TLS certificates for webhook HTTPS
├── deploy-webhook.yaml     # Kubernetes deployment configuration
├── Dockerfile              # Container image build instructions
├── go.mod                  # Go module dependencies
├── go.sum                  # Go module checksums
├── main.go                 # Main application entry point
├── mutatingcontroller.go   # Mutating webhook logic
├── pod.json                # JSON patch template for mutations
├── validatingcontroller.go # Validating webhook logic
└── sample-pod.yaml         # Example pod configuration
```

## Prerequisites

- Go 1.23 or later
- Docker
- Kubernetes cluster with admin access
- `kubectl` configured to access your cluster
- `yq` command-line tool (for YAML processing)

## Setup and Installation

### 1. Generate TLS Certificates

The webhook requires TLS certificates since Kubernetes only allows HTTPS webhooks. Generate certificates in the `certs` directory:

There is a sample certificate that you can use if you are comfortable with it. If not, read how to do do [from this link](https://taskmasterernest.github.io/posts/011-generating-self-signed-certs/).

Then you will have to plug in the generated key and crt values in the `deploy-webhook.yaml` file.

### 2. Build the Docker Image

```bash
docker build -t your-registry/webhook-controller:latest .
docker push your-registry/webhook-controller:latest
```

### 3. Deploy to Kubernetes

Before deploying, update the `deploy-webhook.yaml` with:
- Your image repository
- The base64-encoded CA bundle from your generated certificates
- Any specific namespace configurations

```bash
# Deploy the webhook
kubectl apply -f deploy-webhook.yaml
```

## Testing the Webhook

### 1. Convert Sample Pod to JSON

The project includes a sample pod configuration that can be converted to JSON for testing:

```bash
yq eval '. | to_json' sample-pod.yaml > pod.json
```

### 2. Test Pod Creation

```bash
# Create a pod using the sample configuration
kubectl apply -f sample-pod.yaml
```

## Webhook Behavior

### Validating Webhook

The validating webhook (`validatingcontroller.go`) performs the following checks:
- Validates pod specifications against defined rules
- Returns allow/deny decisions based on validation results
- Cannot modify the incoming request

### Mutating Webhook

The mutating webhook (`mutatingcontroller.go`) can:
- Modify incoming resources using JSON patch operations
- Add/modify labels, annotations, or other fields
- Uses `pod.json` as a template for modifications

## Configuration

### Webhook Configuration

The webhook configuration in `deploy-webhook.yaml` includes:
- Webhook service deployment
- Service account and RBAC rules
- ValidatingWebhookConfiguration
- MutatingWebhookConfiguration

### Patch Template

The `pod.json` file defines the mutation operations to be applied. Modify this file to change the mutation behavior.


## Troubleshooting

### Common Issues

1. Certificate Issues
   - Ensure certificates are properly generated and mounted
   - Verify CA bundle in webhook configuration

2. Webhook Not Triggering
   - Check webhook configuration namespaceSelector
   - Verify service endpoints are correct
   - Check webhook server logs

### Debugging

```bash
# View webhook logs
kubectl logs -l app=webhook-controller

# Check webhook configuration
kubectl get validatingwebhookconfigurations
kubectl get mutatingwebhookconfigurations

# Verify service
kubectl get svc webhook-server
```