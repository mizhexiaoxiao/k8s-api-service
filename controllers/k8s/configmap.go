package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/mizhexiaoxiao/k8s-api-service/models/metadata"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigmapInterface interface {
	Create(ctx context.Context, confMap *v1.ConfigMap) (*v1.ConfigMap, error)
	List(ctx context.Context, queryParam metadata.CommonQueryParameter) ([]v1.ConfigMap, error)
	Delete(ctx context.Context, namespace, name string) error
	Get(ctx context.Context, namespace, name string) (*v1.ConfigMap, error)
	Update(ctx context.Context, namespace, name string, configMap *v1.ConfigMap) (*v1.ConfigMap, error)
}

type ConfigmapOperation struct {
	clientSet *kubernetes.Clientset
}

func NewConfigmapOperation(client *kubernetes.Clientset) ConfigmapInterface {
	return &ConfigmapOperation{
		clientSet: client,
	}
}

func (c ConfigmapOperation) Create(ctx context.Context, confMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	return c.clientSet.CoreV1().ConfigMaps(confMap.Namespace).Create(ctx, confMap, metav1.CreateOptions{})
}

func (c ConfigmapOperation) List(ctx context.Context, queryParam metadata.CommonQueryParameter) ([]v1.ConfigMap, error) {
	configMaps := c.clientSet.CoreV1().ConfigMaps(queryParam.NameSpace)
	option := metav1.ListOptions{LabelSelector: queryParam.LabelSelector}
	result, err := configMaps.List(ctx, option)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("List() configmap failed, err: %s", err))
	}
	return result.Items, nil
}

func (c ConfigmapOperation) Delete(ctx context.Context, namespace, name string) error {
	configMap, err := c.Get(ctx, namespace, name)
	if err != nil {
		return errors.New(fmt.Sprintf("Get() configmap failed, err: %s", err))
	}
	if configMap == nil {
		return errors.New(fmt.Sprintf("configmap with namespace: %s,name: %s not found", namespace, name))
	}
	return c.clientSet.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c ConfigmapOperation) Get(ctx context.Context, namespace, name string) (*v1.ConfigMap, error) {
	return c.clientSet.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c ConfigmapOperation) Update(ctx context.Context, namespace, name string, configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	oldConfigMap, err := c.Get(ctx, namespace, name)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Get() configmap failed, err: %s", err))
	}
	if oldConfigMap == nil {
		return nil, errors.New(fmt.Sprintf("configmap with namespace: %s,name: %s not found", namespace, name))
	}
	configMap.Namespace = namespace
	configMap.Name = name
	return c.clientSet.CoreV1().ConfigMaps(namespace).Update(ctx, configMap, metav1.UpdateOptions{})
}
