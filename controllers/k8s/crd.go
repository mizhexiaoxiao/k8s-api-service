package k8s

import (
	"context"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type CRDInterface interface {
	Create(ctx context.Context, gvk schema.GroupVersionResource, data map[string]interface{}) (*unstructured.Unstructured, error)
	Delete(ctx context.Context, gvk schema.GroupVersionResource, namespace, name string) error
	Get(ctx context.Context, gvk schema.GroupVersionResource, namespace, name string) (*unstructured.Unstructured, error)
	Update(ctx context.Context, gvk schema.GroupVersionResource, namespace, name string, data map[string]interface{}) (*unstructured.Unstructured, error)
}

type CRDOperation struct {
	dyn dynamic.Interface
}

func NewCRDOperation(dyn dynamic.Interface) CRDInterface {
	return &CRDOperation{dyn: dyn}
}

func (o *CRDOperation) Get(ctx context.Context, gvk schema.GroupVersionResource, namespace, name string) (*unstructured.Unstructured, error) {
	return o.dyn.Resource(gvk).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (o *CRDOperation) Create(ctx context.Context, gvk schema.GroupVersionResource, data map[string]interface{}) (*unstructured.Unstructured, error) {
	metadata, ok := data["metadata"].(map[string]interface{})
	if !ok {
		return nil, errors.New("converting data metadata to map failed")
	}
	namespace, ok := metadata["namespace"].(string)
	if !ok {
		return nil, errors.New("converting data namespace to string failed")
	}
	obj := unstructured.Unstructured{Object: data}
	return o.dyn.Resource(gvk).Namespace(namespace).Create(ctx, &obj, metav1.CreateOptions{})
}

func (o *CRDOperation) Update(ctx context.Context, gvk schema.GroupVersionResource, namespace, name string, data map[string]interface{}) (*unstructured.Unstructured, error) {
	metadata, ok := data["metadata"].(map[string]interface{})
	if !ok {
		return nil, errors.New("converting data metadata to map failed")
	}
	metadata["namespace"] = namespace
	metadata["name"] = name
	obj := unstructured.Unstructured{Object: data}
	return o.dyn.Resource(gvk).Namespace(namespace).Update(ctx, &obj, metav1.UpdateOptions{})
}

func (o *CRDOperation) Delete(ctx context.Context, gvk schema.GroupVersionResource, namespace, name string) error {
	return o.dyn.Resource(gvk).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
