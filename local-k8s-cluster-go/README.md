# K3s Cluster Manager (Go)

A modern, type-safe Go application for managing K3s clusters and application deployments.

## Features

- **Cluster Management**: Monitor cluster health, nodes, and resources
- **Application Deployment**: Deploy and manage applications with Kubernetes manifests
- **Infrastructure Setup**: Automated K3s installation and platform services setup
- **Type Safety**: Leverage Go's type system for reliable infrastructure management
- **CLI Interface**: User-friendly command-line interface powered by Cobra

## Project Structure

```
local-k8s-cluster-go/
├── cmd/                    # Application entrypoint
│   └── main.go
├── internal/              # Internal packages
│   ├── cli/              # CLI commands and flags
│   │   ├── root.go       # Root command and configuration
│   │   ├── cluster.go    # Cluster management commands
│   │   ├── apps.go       # Application management commands
│   │   └── setup.go      # Infrastructure setup commands
│   ├── k8s/              # Kubernetes client wrapper
│   │   └── client.go     # K8s client with helper methods
│   ├── apps/             # Application management logic
│   │   └── manager.go    # Deploy, list, and manage apps
│   └── setup/            # Infrastructure setup logic
│       └── installer.go  # K3s and platform service installation
├── pkg/                  # Public packages (if any)
├── manifests/            # Kubernetes manifests for applications
├── scripts/              # Helper scripts
├── go.mod               # Go module definition
└── README.md           # This file
```

## Installation

1. **Prerequisites**:
   - Go 1.21 or later
   - Access to a Kubernetes cluster (or ability to install K3s)

2. **Build the application**:
   ```bash
   go build -o k3s-manager ./cmd
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

## Usage

### Global Options
```bash
k3s-manager --help
k3s-manager --verbose           # Enable verbose logging
k3s-manager --config /path/to/config.yaml
```

### Cluster Management
```bash
# Check cluster status and health
k3s-manager cluster status

# Get detailed cluster information
k3s-manager cluster info
```

### Application Management
```bash
# Deploy an application
k3s-manager apps deploy my-app --namespace production

# List all applications
k3s-manager apps list --namespace production

# Get application status
k3s-manager apps status my-app --namespace production

# Deploy with dry-run
k3s-manager apps deploy my-app --dry-run
```

### Infrastructure Setup
```bash
# Setup K3s server node
k3s-manager setup server

# Setup K3s worker node
k3s-manager setup worker --server-ip 192.168.1.100 --node-token <token>

# Setup platform services (ingress, ArgoCD, monitoring)
k3s-manager setup platform
```

## Configuration

Create a configuration file at `~/.k3s-manager.yaml`:

```yaml
verbose: false
cluster:
  kubeconfig: ~/.kube/config
apps:
  default-namespace: default
  manifest-directory: ./manifests
setup:
  install-k3s: true
  setup-ingress: true
  setup-argocd: true
```

## Application Manifests

Place your Kubernetes manifests in the `manifests/` directory:

```
manifests/
├── my-app/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
└── another-app/
    └── all-in-one.yaml
```

## Key Go Dependencies

- **Cobra**: CLI framework for command structure and flags
- **Viper**: Configuration management with YAML/ENV support  
- **client-go**: Official Kubernetes Go client library
- **k8s.io/api**: Kubernetes API types and structures
- **k8s.io/apimachinery**: Kubernetes API machinery and utilities

## Development

### Running locally
```bash
go run ./cmd --help
```

### Running tests
```bash
go test ./...
```

### Building for different platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o k3s-manager-linux ./cmd

# macOS
GOOS=darwin GOARCH=amd64 go build -o k3s-manager-darwin ./cmd

# Windows
GOOS=windows GOARCH=amd64 go build -o k3s-manager.exe ./cmd
```

## Advantages of Go for Infrastructure

1. **Single Binary**: No runtime dependencies, easy deployment
2. **Type Safety**: Compile-time error checking prevents runtime issues
3. **Performance**: Fast execution, low memory footprint
4. **Concurrency**: Built-in goroutines for parallel operations
5. **Ecosystem**: Rich Kubernetes ecosystem with official client libraries
6. **Cross-Platform**: Build for any target platform from any development machine

## Extending the Tool

To add new commands:

1. Create a new file in `internal/cli/`
2. Define your command using Cobra patterns
3. Add business logic in appropriate `internal/` packages
4. Register the command in the CLI structure

Example command structure:
```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description of my command",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation here
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

This Go-based approach provides a solid foundation for managing your K3s infrastructure with type safety, performance, and maintainability.