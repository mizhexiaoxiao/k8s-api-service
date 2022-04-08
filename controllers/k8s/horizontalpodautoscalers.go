package k8s

import (
	"context"
	"errors"
	"github.com/mizhexiaoxiao/k8s-api-service/models/metadata"
	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HorizontalPodAutoScalerInterface interface {
	Create(ctx context.Context, scaler *v1.HorizontalPodAutoscaler) (*v1.HorizontalPodAutoscaler, error)
	List(ctx context.Context, queryParam metadata.CommonQueryParameter) ([]v1.HorizontalPodAutoscaler, error)
	Get(ctx context.Context, namespace, name string) (*v1.HorizontalPodAutoscaler, error)
	Update(ctx context.Context, namespace, name string, scaler *v1.HorizontalPodAutoscaler) (*v1.HorizontalPodAutoscaler, error)
	Delete(ctx context.Context, namespace, name string) error
}

type HorizontalPodAutoScalerOperation struct {
	clientSet *kubernetes.Clientset
}

func NewHorizontalPodAutoScalerOperation(client *kubernetes.Clientset) HorizontalPodAutoScalerInterface {
	return &HorizontalPodAutoScalerOperation{
		clientSet: client,
	}
}

func (o *HorizontalPodAutoScalerOperation) Create(ctx context.Context, scaler *v1.HorizontalPodAutoscaler) (*v1.HorizontalPodAutoscaler, error) {
	autoscalers := o.clientSet.AutoscalingV1().HorizontalPodAutoscalers(scaler.Namespace)
	return autoscalers.Create(ctx, scaler, metav1.CreateOptions{})
}

func (o *HorizontalPodAutoScalerOperation) List(ctx context.Context, queryParam metadata.CommonQueryParameter) ([]v1.HorizontalPodAutoscaler, error) {
	autoscalers := o.clientSet.AutoscalingV1().HorizontalPodAutoscalers(queryParam.NameSpace)
	option := metav1.ListOptions{LabelSelector: queryParam.LabelSelector}
	result, err := autoscalers.List(ctx, option)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (o *HorizontalPodAutoScalerOperation) Get(ctx context.Context, namespace, name string) (*v1.HorizontalPodAutoscaler, error) {
	autoscalers := o.clientSet.AutoscalingV1().HorizontalPodAutoscalers(namespace)
	return autoscalers.Get(ctx, name, metav1.GetOptions{})
}

func (o *HorizontalPodAutoScalerOperation) Update(ctx context.Context, namespace, name string, scaler *v1.HorizontalPodAutoscaler) (*v1.HorizontalPodAutoscaler, error) {
	originScaler, err := o.Get(ctx, namespace, name)
	if err != nil {
		return nil, err
	}
	if originScaler != nil && originScaler.Name == name && originScaler.Namespace == namespace {
		return o.clientSet.AutoscalingV1().HorizontalPodAutoscalers(scaler.Namespace).Update(ctx, scaler, metav1.UpdateOptions{})
	}
	return nil, errors.New("originScaler not match update object")
}

func (o *HorizontalPodAutoScalerOperation) Delete(ctx context.Context, namespace, name string) error {
	originScaler, err := o.Get(ctx, namespace, name)
	if err != nil {
		return err
	}
	if originScaler != nil {
		return o.clientSet.AutoscalingV1().HorizontalPodAutoscalers(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	}
	return errors.New("originScaler not match delete object")
}
