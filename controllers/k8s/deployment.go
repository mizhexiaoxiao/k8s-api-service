package k8s

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentInterface interface {
	Create(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
}

type DeploymentOperation struct {
	clientSet *kubernetes.Clientset
}

func NewDeploymentOperation(client *kubernetes.Clientset) DeploymentInterface {
	return DeploymentOperation{
		clientSet: client,
	}
}

func (o DeploymentOperation) Create(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	deploymentsClient := o.clientSet.AppsV1().Deployments(namespace)
	return deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
}
