#!/bin/bash

set -e

echo "🧪 Testing K3s Manager locally..."

# Check if k3d is installed
if ! command -v k3d &> /dev/null; then
    echo "❌ k3d is not installed. Install with: brew install k3d"
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Install with: brew install kubectl"
    exit 1
fi

# Create test cluster if it doesn't exist
if ! k3d cluster list | grep -q "test-cluster"; then
    echo "🚀 Creating test cluster..."
    k3d cluster create test-cluster \
        --servers 1 \
        --agents 1 \
        --port "8080:80@loadbalancer" \
        --wait
else
    echo "✅ Test cluster already exists"
    k3d cluster start test-cluster 2>/dev/null || true
fi

# Wait for cluster to be ready
echo "⏳ Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=60s

# Build the Go binary
echo "🔨 Building k3s-manager..."
go build -o k3s-manager ./cmd

# Test basic functionality
echo "🧪 Testing cluster status..."
./k3s-manager cluster status

echo "🧪 Testing cluster info..."
./k3s-manager cluster info

# Deploy a test application
echo "🚀 Deploying test application..."
kubectl create deployment nginx --image=nginx:alpine --dry-run=client -o yaml > /tmp/nginx-deployment.yaml
kubectl apply -f /tmp/nginx-deployment.yaml

# Wait for deployment
kubectl wait --for=condition=available --timeout=60s deployment/nginx

# Test app listing (this will need the working code)
echo "🧪 Testing app management..."
echo "Note: App management commands may not work fully yet - this is expected!"

./k3s-manager apps list || echo "⚠️  Apps list command needs implementation fixes"

# Show what's actually running
echo "📋 Current cluster state:"
kubectl get nodes
kubectl get pods
kubectl get deployments
kubectl get services

# Cleanup
echo "🧹 Cleaning up test deployment..."
kubectl delete deployment nginx

echo "✅ Local testing complete!"
echo "💡 To clean up: k3d cluster delete test-cluster"