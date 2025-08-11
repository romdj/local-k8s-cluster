package k8s

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Client wraps the K3s/Kubernetes client with additional functionality
type Client struct {
	clientset kubernetes.Interface
	config    *rest.Config
}

// ClusterStatus represents the current state of the cluster
type ClusterStatus struct {
	Status        string
	ReadyNodes    int
	TotalNodes    int
	RunningPods   int
	TotalPods     int
	Namespaces    int
	UnhealthyPods []PodInfo
}

// ClusterInfo contains detailed cluster information
type ClusterInfo struct {
	Version   string
	Platform  string
	APIServer string
	Nodes     []NodeInfo
}

// NodeInfo represents node details
type NodeInfo struct {
	Name   string
	Role   string
	Status string
}

// PodInfo represents pod details
type PodInfo struct {
	Name      string
	Namespace string
	Phase     string
}

// NewClient creates a new K3s cluster client
func NewClient() (*Client, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Client{
		clientset: clientset,
		config:    config,
	}, nil
}

// GetClusterStatus returns comprehensive cluster status
func (c *Client) GetClusterStatus(ctx context.Context) (*ClusterStatus, error) {
	status := &ClusterStatus{}

	// Get nodes
	nodes, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	status.TotalNodes = len(nodes.Items)
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
				status.ReadyNodes++
				break
			}
		}
	}

	// Get all pods
	pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	status.TotalPods = len(pods.Items)
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			status.RunningPods++
		} else if pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodPending {
			status.UnhealthyPods = append(status.UnhealthyPods, PodInfo{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Phase:     string(pod.Status.Phase),
			})
		}
	}

	// Get namespaces
	namespaces, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	status.Namespaces = len(namespaces.Items)

	// Determine overall status
	if status.ReadyNodes == status.TotalNodes && len(status.UnhealthyPods) == 0 {
		status.Status = "Healthy"
	} else if status.ReadyNodes > 0 {
		status.Status = "Degraded"
	} else {
		status.Status = "Unhealthy"
	}

	return status, nil
}

// GetClusterInfo returns detailed cluster information
func (c *Client) GetClusterInfo(ctx context.Context) (*ClusterInfo, error) {
	info := &ClusterInfo{}

	// Get server version
	version, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}
	info.Version = version.GitVersion
	info.Platform = version.Platform

	// API Server endpoint
	info.APIServer = c.config.Host

	// Get nodes with details
	nodes, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	for _, node := range nodes.Items {
		role := "worker"
		if _, exists := node.Labels["node-role.kubernetes.io/control-plane"]; exists {
			role = "control-plane"
		} else if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
			role = "master"
		}

		status := "NotReady"
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
				status = "Ready"
				break
			}
		}

		info.Nodes = append(info.Nodes, NodeInfo{
			Name:   node.Name,
			Role:   role,
			Status: status,
		})
	}

	return info, nil
}

// WaitForPod waits for a pod to be ready
func (c *Client) WaitForPod(ctx context.Context, namespace, labelSelector string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for pod with selector %s", labelSelector)
		default:
			pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
				LabelSelector: labelSelector,
			})
			if err != nil {
				return fmt.Errorf("failed to list pods: %w", err)
			}

			if len(pods.Items) > 0 {
				pod := pods.Items[0]
				if pod.Status.Phase == v1.PodRunning {
					// Check if all containers are ready
					ready := true
					for _, status := range pod.Status.ContainerStatuses {
						if !status.Ready {
							ready = false
							break
						}
					}
					if ready {
						return nil
					}
				}
			}

			time.Sleep(2 * time.Second)
		}
	}
}