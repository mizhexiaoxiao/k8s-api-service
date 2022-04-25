package k8s

import (
	"context"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentInterface interface {
	Get(ctx context.Context, namespace, name string) (*appsv1.Deployment, error)
	Create(ctx context.Context, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(ctx context.Context, namespace, name string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
}

type DeploymentOperation struct {
	clientSet *kubernetes.Clientset
}

func NewDeploymentOperation(client *kubernetes.Clientset) DeploymentInterface {
	return DeploymentOperation{
		clientSet: client,
	}
}

func (o DeploymentOperation) Get(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	return o.clientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (o DeploymentOperation) Create(ctx context.Context, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	deploymentsClient := o.clientSet.AppsV1().Deployments(deployment.Namespace)
	return deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
}

func (o DeploymentOperation) Update(ctx context.Context, namespace, name string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	_, err := o.Get(ctx, namespace, name)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("DeploymentOperation of Get deployment failed, err: %s", err))
	}
	deployment, err = o.clientSet.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("DeploymentOperation of Update deployment failed, err: %s", err))
	}
	return deployment, nil
}
