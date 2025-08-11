# Getting Started with K3s Cluster Manager

Welcome to the K3s Cluster Manager! This guide will help you get up and running with modern, type-safe Kubernetes cluster management using our Go-based tool.

## 🎯 What You'll Learn

- How to set up a local K3s development environment
- Install and use the k3s-manager CLI tool
- Deploy and manage applications on your cluster
- Monitor cluster health and performance
- Set up automated workflows

## 📋 Prerequisites

### Required
- **macOS, Linux, or Windows** (amd64 or arm64)
- **Docker** installed and running
- **kubectl** CLI tool

### Optional (for development)
- **Go 1.21+** (if building from source)
- **Git** (for contributing)

## 🚀 Quick Start (5 minutes)

### 1. Create a Local K3s Cluster

Using **k3d** (recommended for local development):

```bash
# Install k3d
brew install k3d                    # macOS
# or
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

# Create cluster
k3d cluster create my-cluster \
  --servers 1 \
  --agents 2 \
  --port "8080:80@loadbalancer" \
  --port "8443:443@loadbalancer"

# Verify cluster is running
kubectl get nodes
```

### 2. Install K3s Manager

#### Option A: Download Pre-built Binary (Recommended)

```bash
# Download latest release for your platform
# macOS Apple Silicon
curl -L -o k3s-manager.tar.gz https://github.com/romdj/local-k8s-cluster/releases/latest/download/k3s-manager-darwin-arm64.tar.gz

# macOS Intel
curl -L -o k3s-manager.tar.gz https://github.com/romdj/local-k8s-cluster/releases/latest/download/k3s-manager-darwin-amd64.tar.gz

# Linux
curl -L -o k3s-manager.tar.gz https://github.com/romdj/local-k8s-cluster/releases/latest/download/k3s-manager-linux-amd64.tar.gz

# Extract and install
tar -xzf k3s-manager.tar.gz
sudo mv k3s-manager-* /usr/local/bin/k3s-manager
chmod +x /usr/local/bin/k3s-manager

# Verify installation
k3s-manager --version
```

#### Option B: Build from Source

```bash
# Clone repository
git clone https://github.com/romdj/local-k8s-cluster.git
cd local-k8s-cluster/local-k8s-cluster-go

# Build binary
go build -o k3s-manager ./cmd

# Install to PATH
sudo mv k3s-manager /usr/local/bin/
```

### 3. Verify Everything Works

```bash
# Check cluster status
k3s-manager cluster status

# Get detailed cluster info
k3s-manager cluster info

# List applications (should be empty initially)
k3s-manager apps list
```

Expected output:
```
$ k3s-manager cluster status
Cluster Status: Healthy
Nodes: 3 ready, 3 total
Pods: 8 running, 8 total
Namespaces: 4

$ k3s-manager cluster info
Kubernetes Version: v1.28.2+k3s1
Platform: linux/amd64
API Server: https://0.0.0.0:35637
Nodes:
  - k3d-my-cluster-server-0 (control-plane): Ready
  - k3d-my-cluster-agent-0 (worker): Ready
  - k3d-my-cluster-agent-1 (worker): Ready
```

## 🛠 Basic Usage

### Cluster Management

```bash
# Monitor cluster health
k3s-manager cluster status

# Get comprehensive cluster information
k3s-manager cluster info

# Enable verbose output for troubleshooting
k3s-manager --verbose cluster status
```

### Application Management

#### Deploy a Sample Application

Create a simple nginx deployment:

```bash
# Create manifests directory
mkdir -p manifests/nginx

# Create deployment manifest
cat > manifests/nginx/deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
  namespace: default
spec:
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
EOF

# Deploy using kubectl (k3s-manager app deployment coming soon)
kubectl apply -f manifests/nginx/

# Check deployment status
kubectl get pods,services
```

#### Access Your Application

```bash
# With k3d port mapping, access via localhost
curl http://localhost:8080

# Or use port-forwarding
kubectl port-forward service/nginx 8090:80
curl http://localhost:8090
```

#### Monitor Applications

```bash
# List all applications in default namespace
k3s-manager apps list

# Check specific application status
k3s-manager apps status nginx

# Monitor with watch
watch k3s-manager cluster status
```

## 🏗 Production Server Setup

For setting up K3s on a production Linux server, see our comprehensive guide:

### Quick Server Setup

```bash
# On Ubuntu/Debian server
curl -sfL https://get.k3s.io | sh -s - \
  --write-kubeconfig-mode 644 \
  --node-name production-server

# Copy kubeconfig for remote access
scp user@server:/etc/rancher/k3s/k3s.yaml ~/.kube/config-production
sed -i 's/127.0.0.1/YOUR_SERVER_IP/g' ~/.kube/config-production

# Use remote cluster
export KUBECONFIG=~/.kube/config-production
k3s-manager cluster status
```

See [`server_init.md`](server_init.md) for detailed production setup instructions.

## 📊 Common Workflows

### Development Workflow

```bash
# 1. Start local cluster
k3d cluster create dev

# 2. Check cluster health
k3s-manager cluster status

# 3. Deploy your application
kubectl apply -f your-app/

# 4. Monitor and iterate
watch k3s-manager cluster status

# 5. Clean up when done
k3d cluster delete dev
```

### Multi-Environment Setup

```bash
# Development cluster
k3d cluster create dev --port "8080:80@loadbalancer"

# Staging cluster  
k3d cluster create staging --port "8081:80@loadbalancer"

# Switch between clusters
kubectl config use-context k3d-dev
k3s-manager cluster status

kubectl config use-context k3d-staging
k3s-manager cluster status
```

### Application Testing

```bash
# Deploy test version
kubectl apply -f manifests/app/

# Check status
k3s-manager apps list
k3s-manager apps status myapp

# Load test
curl -X GET http://localhost:8080/health

# Rolling update
kubectl set image deployment/myapp container=myapp:v2

# Monitor rollout
kubectl rollout status deployment/myapp
```

## 🔧 Configuration

### CLI Configuration

Create `~/.k3s-manager.yaml`:

```yaml
verbose: false
cluster:
  kubeconfig: ~/.kube/config
  context: ""
apps:
  default-namespace: default
  manifest-directory: ./manifests
output:
  format: table  # table, json, yaml
```

### Environment Variables

```bash
# Override kubeconfig location
export KUBECONFIG=~/.kube/my-cluster-config

# Enable debug logging
export K3S_MANAGER_DEBUG=true

# Set default namespace
export K3S_MANAGER_NAMESPACE=production
```

## 🏷 Working with Multiple Clusters

### Cluster Contexts

```bash
# List available contexts
kubectl config get-contexts

# Switch context
kubectl config use-context k3d-production

# Check current cluster
k3s-manager cluster info

# Use specific config file
k3s-manager --config ~/.kube/staging-config cluster status
```

### Namespace Management

```bash
# Create namespace
kubectl create namespace myapp-prod

# Deploy to specific namespace
k3s-manager apps list --namespace myapp-prod

# Set default namespace
kubectl config set-context --current --namespace=myapp-prod
```

## 📈 Monitoring and Observability

### Basic Monitoring

```bash
# Cluster overview
k3s-manager cluster status

# Resource usage
kubectl top nodes
kubectl top pods --all-namespaces

# Event monitoring
kubectl get events --sort-by=.metadata.creationTimestamp
```

### Health Checks

```bash
# Check cluster components
kubectl get componentstatuses

# Verify DNS
kubectl run test-dns --image=busybox --rm -it -- nslookup kubernetes.default

# Check ingress
curl -H "Host: myapp.local" http://localhost:8080
```

## 🔍 Troubleshooting

### Common Issues

#### Cluster Not Starting
```bash
# Check Docker is running
docker ps

# Recreate cluster
k3d cluster delete my-cluster
k3d cluster create my-cluster --wait

# Check cluster logs
k3d cluster list
kubectl get events
```

#### kubectl Connection Issues
```bash
# Verify kubeconfig
kubectl config current-context
kubectl cluster-info

# Reset kubeconfig
k3d kubeconfig write my-cluster --overwrite

# Test connection
kubectl get nodes
```

#### Application Not Accessible
```bash
# Check service endpoints
kubectl get endpoints

# Verify port mappings
k3d cluster list

# Check ingress controller
kubectl get pods -n kube-system
```

#### Performance Issues
```bash
# Check resource usage
kubectl top nodes
kubectl describe node

# Check for failed pods
kubectl get pods --all-namespaces --field-selector=status.phase=Failed

# Review resource limits
kubectl describe pod myapp-pod
```

### Getting Help

```bash
# Built-in help
k3s-manager --help
k3s-manager cluster --help
k3s-manager apps --help

# Version information
k3s-manager --version

# Verbose output for debugging
k3s-manager --verbose cluster status
```

## 🚀 Next Steps

### Advanced Features

1. **Set up GitOps** with ArgoCD
2. **Add monitoring** with Prometheus and Grafana
3. **Configure ingress** with custom domains
4. **Set up CI/CD** with GitHub Actions
5. **Scale to production** with our server setup guide

### Learning Resources

- [Platform Architecture](platform_architecture.md) - Deep dive into our design
- [Server Setup](server_init.md) - Production deployment guide
- [Contributing](CONTRIBUTING.md) - Join the development
- [Decision Log](decision_log.md) - Understand our choices

### Community

- **Issues**: [GitHub Issues](https://github.com/romdj/local-k8s-cluster/issues)
- **Discussions**: [GitHub Discussions](https://github.com/romdj/local-k8s-cluster/discussions)
- **Releases**: [Latest Releases](https://github.com/romdj/local-k8s-cluster/releases)

## 🎉 Success!

You now have a fully functional K3s cluster with modern management tools. You can:

- ✅ Monitor cluster health and performance
- ✅ Deploy and manage applications  
- ✅ Scale from development to production
- ✅ Use type-safe, reliable tooling

Happy clustering! 🚀

---

**Questions?** Check our [troubleshooting section](#🔍-troubleshooting) or [open an issue](https://github.com/romdj/local-k8s-cluster/issues/new).