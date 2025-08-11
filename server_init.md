# Server Initialization Guide

Complete setup guide for a vanilla Linux server to production-ready K3s cluster hosting 5-25 applications.

## Distribution Recommendation

### **Ubuntu Server 24.04 LTS** (Recommended)
- **Best balance**: Ease of use + production stability
- **K3s optimized**: Built-in vxlan support, no additional modules needed
- **Community support**: Extensive documentation and troubleshooting resources
- **Long-term support**: 5 years of security updates

### Alternative Options:
- **Rocky Linux 9**: For RHEL compatibility + enterprise features
- **Debian 12**: Maximum stability, conservative updates
- **Alpine Linux**: Ultra-lightweight (advanced users only)

## Hardware Requirements

### Minimum Specifications
```
CPU: 4 cores (8+ recommended)
RAM: 8GB (16GB+ recommended for 15+ apps)
Storage: 100GB SSD (200GB+ recommended)
Network: 1Gbps ethernet
```

### Optimal Specifications
```
CPU: 8-16 cores
RAM: 16-32GB
Storage: 500GB NVMe SSD
Network: 1Gbps+ with static IP
```

## Initial Server Setup

### 1. Base Ubuntu Server Installation
```bash
# During installation, select:
# - OpenSSH server
# - No snap packages initially
# - LVM for disk management
# - Entire disk for simplicity
```

### 2. Initial System Configuration
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install essential packages
sudo apt install -y \
  curl \
  wget \
  git \
  vim \
  htop \
  unzip \
  software-properties-common \
  apt-transport-https \
  ca-certificates \
  gnupg \
  lsb-release

# Set timezone
sudo timedatectl set-timezone UTC

# Configure hostname (replace 'k3s-server' with your preference)
sudo hostnamectl set-hostname k3s-server
echo "127.0.1.1 k3s-server" | sudo tee -a /etc/hosts
```

### 3. User and SSH Configuration
```bash
# Create non-root user for K3s management
sudo adduser k3suser
sudo usermod -aG sudo k3suser

# SSH key setup (run from your local machine)
ssh-copy-id k3suser@YOUR_SERVER_IP

# SSH hardening
sudo vim /etc/ssh/sshd_config
# Add/modify these lines:
# PermitRootLogin no
# PasswordAuthentication no
# PubkeyAuthentication yes
# Port 22 (or change to non-standard port)

sudo systemctl restart sshd
```

### 4. Firewall Configuration
```bash
# Enable UFW firewall
sudo ufw enable

# Allow SSH
sudo ufw allow 22/tcp

# Allow K3s API server
sudo ufw allow 6443/tcp

# Allow K3s worker nodes (if multi-node)
sudo ufw allow from 10.42.0.0/16  # Pod CIDR
sudo ufw allow from 10.43.0.0/16  # Service CIDR

# Allow HTTP/HTTPS for applications
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Check status
sudo ufw status verbose
```

## K3s Installation

### 1. Install K3s Server (Single Node)
```bash
# Install K3s with recommended settings
curl -sfL https://get.k3s.io | sh -s - \
  --write-kubeconfig-mode 644 \
  --disable traefik \
  --disable servicelb \
  --node-name k3s-master

# Verify installation
sudo k3s kubectl get nodes
sudo k3s kubectl get pods -A
```

### 2. Configure kubectl Access
```bash
# Copy kubeconfig for regular user
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Test kubectl access
kubectl get nodes
kubectl get namespaces
```

### 3. Multi-Node Setup (Optional)
```bash
# On master node, get join token
sudo cat /var/lib/rancher/k3s/server/node-token

# On worker nodes, join cluster
curl -sfL https://get.k3s.io | K3S_URL=https://MASTER_IP:6443 \
  K3S_TOKEN=NODE_TOKEN sh -

# Verify nodes joined
kubectl get nodes
```

## Essential Tool Installation

### 1. Helm Package Manager
```bash
# Install Helm
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt update
sudo apt install helm

# Verify installation
helm version
```

### 2. Docker (for local development)
```bash
# Install Docker (optional, for local builds)
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
```

### 3. Additional Utilities
```bash
# Install useful tools
sudo apt install -y \
  jq \
  yq \
  tree \
  ncdu \
  iotop \
  nethogs \
  fail2ban

# Install kubectx/kubens for context switching
sudo git clone https://github.com/ahmetb/kubectx /opt/kubectx
sudo ln -s /opt/kubectx/kubectx /usr/local/bin/kubectx
sudo ln -s /opt/kubectx/kubens /usr/local/bin/kubens
```

## Platform Services Installation

### 1. Ingress Controller (NGINX)
```bash
# Install NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml

# Wait for deployment
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

# Verify
kubectl get pods -n ingress-nginx
```

### 2. Cert-Manager (SSL Certificates)
```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.2/cert-manager.yaml

# Wait for deployment
kubectl wait --namespace cert-manager \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/instance=cert-manager \
  --timeout=120s

# Verify
kubectl get pods -n cert-manager
```

### 3. ArgoCD (GitOps)
```bash
# Create argocd namespace
kubectl create namespace argocd

# Install ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for deployment
kubectl wait --namespace argocd \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/name=argocd-server \
  --timeout=300s

# Get initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo

# Port forward to access UI (run in separate terminal)
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

## Storage Configuration

### 1. Local Storage Setup
```bash
# Create storage directories
sudo mkdir -p /opt/k3s-storage/{databases,uploads,logs}
sudo chown -R k3suser:k3suser /opt/k3s-storage
sudo chmod -R 755 /opt/k3s-storage

# Create StorageClass for local storage
cat <<EOF | kubectl apply -f -
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-storage
provisioner: rancher.io/local-path
parameters:
  nodePath: /opt/k3s-storage
reclaimPolicy: Retain
allowVolumeExpansion: true
EOF
```

## Security Hardening

### 1. System Security
```bash
# Configure fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Set up automatic security updates
echo 'Unattended-Upgrade::Automatic-Reboot "false";' | sudo tee -a /etc/apt/apt.conf.d/50unattended-upgrades
sudo systemctl enable unattended-upgrades

# Configure log rotation
sudo vim /etc/logrotate.d/k3s
# Add:
# /var/lib/rancher/k3s/server/logs/*.log {
#     daily
#     missingok
#     rotate 7
#     compress
#     notifempty
#     create 0644 root root
# }
```

### 2. K3s Security
```bash
# Create restricted pod security policy
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: secure-apps
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
EOF

# Configure network policies (example)
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: default
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
EOF
```

## Monitoring Setup (Basic)

### 1. Node Exporter
```bash
# Install node exporter for metrics
wget https://github.com/prometheus/node_exporter/releases/download/v1.6.1/node_exporter-1.6.1.linux-amd64.tar.gz
tar xvfz node_exporter-1.6.1.linux-amd64.tar.gz
sudo mv node_exporter-1.6.1.linux-amd64/node_exporter /usr/local/bin/
sudo chown root:root /usr/local/bin/node_exporter

# Create systemd service
sudo tee /etc/systemd/system/node_exporter.service > /dev/null <<EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF

# Create user and start service
sudo useradd --no-create-home --shell /bin/false node_exporter
sudo systemctl daemon-reload
sudo systemctl enable node_exporter
sudo systemctl start node_exporter
```

## Backup Configuration

### 1. K3s Backup Script
```bash
# Create backup directory
sudo mkdir -p /opt/k3s-backups

# Create backup script
sudo tee /opt/k3s-backups/backup-k3s.sh > /dev/null <<'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/k3s-backups"

# Backup K3s data
sudo cp -r /var/lib/rancher/k3s/server/db "$BACKUP_DIR/db_$DATE"

# Backup configuration
kubectl get all --all-namespaces -o yaml > "$BACKUP_DIR/cluster_state_$DATE.yaml"

# Remove backups older than 7 days
find "$BACKUP_DIR" -name "db_*" -mtime +7 -exec rm -rf {} \;
find "$BACKUP_DIR" -name "cluster_state_*" -mtime +7 -delete

echo "Backup completed: $DATE"
EOF

sudo chmod +x /opt/k3s-backups/backup-k3s.sh

# Add to crontab (daily backup at 2 AM)
echo "0 2 * * * /opt/k3s-backups/backup-k3s.sh >> /var/log/k3s-backup.log 2>&1" | sudo crontab -
```

## Verification and Testing

### 1. Cluster Health Check
```bash
# Check node status
kubectl get nodes -o wide

# Check system pods
kubectl get pods -A

# Check storage
kubectl get storageclass

# Check ingress
kubectl get ingressclass

# Test DNS resolution
kubectl run test-dns --image=busybox --rm -it --restart=Never -- nslookup kubernetes.default
```

### 2. Deploy Test Application
```bash
# Deploy simple test app
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
      - name: test-app
        image: nginx:alpine
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: test-app
  namespace: default
spec:
  selector:
    app: test-app
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
EOF

# Verify deployment
kubectl get pods
kubectl get services

# Clean up test
kubectl delete deployment test-app
kubectl delete service test-app
```

## Post-Installation Checklist

- [ ] K3s cluster running and accessible
- [ ] kubectl configured for regular user
- [ ] Firewall configured with appropriate ports
- [ ] Ingress controller installed and ready
- [ ] ArgoCD installed and accessible
- [ ] Storage class configured
- [ ] Security hardening applied
- [ ] Backup script configured
- [ ] Monitoring basics in place
- [ ] Test application deployed successfully

## Next Steps

1. **Set up development workflow** with GitHub Actions
2. **Install monitoring stack** (Prometheus + Grafana)
3. **Configure DNS** for application access
4. **Deploy first production application**
5. **Set up log aggregation** (optional)

Your K3s server is now ready for production application deployments!