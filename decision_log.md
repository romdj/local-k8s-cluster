# Decision Log

This document tracks architectural and technical decisions made for the local Kubernetes cluster project.

## 2025-08-04

### Decision: K3s over Standard Kubernetes
**Status**: Accepted  
**Context**: Need to choose between standard Kubernetes (K8s) and K3s for local cluster setup.

**Decision**: Selected K3s as the Kubernetes distribution.

**Rationale**:
- Resource efficiency: <100MB vs 1.5GB+ for standard K8s
- Simplified setup: Single binary with built-in SQLite (no etcd complexity)
- Perfect for local development and learning
- 100% Kubernetes API compatible and CNCF certified
- Built-in components (Traefik, storage provider, load balancer)
- Can migrate to full K8s later if production needs require it

**Alternatives Considered**:
- Standard Kubernetes: Too resource-heavy for local development
- OpenShift Local: More complex setup, heavier resource usage
- Kind: Good for testing but less production-like than K3s

### Decision: GitHub Actions for CI/CD
**Status**: Accepted  
**Context**: Need CI/CD solution for local Kubernetes cluster deployments.

**Decision**: Use GitHub Actions as primary CI/CD platform.

**Rationale**:
- Native integration with GitHub repositories
- No external CI/CD service dependencies
- Built-in container registry (ghcr.io)
- Cost-effective for local development
- Extensive marketplace of Kubernetes-focused actions
- Can run self-hosted runners on local cluster if needed

**Alternatives Considered**:
- Jenkins: Too complex for local setup
- GitLab CI: Would require GitLab migration
- ArgoCD: Will consider for GitOps layer later

### Decision: Technical Stack Components
**Status**: Accepted  
**Context**: Define supporting tools and technologies.

**Decisions**:
- **Container Registry**: GitHub Container Registry (ghcr.io)
- **Infrastructure as Code**: Kubernetes YAML manifests + Helm charts
- **Configuration Management**: Kustomize for environment-specific configs
- **Monitoring**: Prometheus + Grafana (Phase 2)

**Implementation Phases**:
1. **Phase 1**: Basic K3s + GitHub Actions + simple deployments
2. **Phase 2**: Add Helm, monitoring, GitOps with ArgoCD
3. **Phase 3**: Advanced security, multi-environment support

### Decision: K3d as K3s Implementation
**Status**: Accepted  
**Context**: Need to choose specific K3s implementation method for macOS development.

**Decision**: Use k3d (K3s in Docker) as the primary K3s implementation.

**Rationale**:
- Native macOS support (K3s needs Linux, k3d runs in Docker)
- Fast cluster creation/destruction (seconds vs minutes)
- Multi-node cluster support for testing
- Port mapping for ingress/load balancer access
- Easier GitHub Actions integration
- Docker ecosystem familiarity

**Target Scale**: 5-25 applications running simultaneously on local cluster.

**Alternatives Considered**:
- Multipass + K3s: More resource overhead, VM management complexity
- Rancher Desktop: GUI overhead, less automation-friendly
- Native K3s: Requires Linux VM on macOS

### Decision: Skip Traditional Configuration Management Tools
**Status**: Accepted  
**Context**: Evaluated Chef, Puppet, Ansible for managing 5-25 applications on Kubernetes.

**Decision**: Do not use Chef, Puppet, or Ansible for Kubernetes application management.

**Rationale**:
- Traditional tools designed for server configuration management (pre-container era)
- Kubernetes provides declarative configuration through YAML manifests
- Modern cloud-native tools better suited for container orchestration
- Immutable container paradigm conflicts with mutable server configuration
- 2025 ecosystem has shifted to GitOps and cloud-native approaches

**Modern Toolstack Selection**:
- **Application Packaging**: Helm charts
- **GitOps Deployments**: ArgoCD  
- **Environment Configs**: Kustomize
- **Traffic Management**: Traefik (built into K3s)
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: GitHub Actions

**Architecture Pattern**: 
- Git-based declarative configurations
- Namespace-based environment separation (dev/staging/production)
- Per-application Helm charts for consistency
- Centralized ingress and shared services

**Note**: Ansible may still be useful for cluster provisioning and host setup, but not for application deployment within Kubernetes.