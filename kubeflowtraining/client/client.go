/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"fmt"
	"log"

	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	pkgApi "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
)

//go:generate mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

type Client interface {
	// PyTorchJob CRUD operations
	CreatePytorchJob(vm *kubeflowv1.PyTorchJob) error
	GetPytorchJob(namespace string, name string) (*kubeflowv1.PyTorchJob, error)
	UpdatePytorchJob(namespace string, name string, vm *kubeflowv1.PyTorchJob, data []byte) error
	DeletePytorchJob(namespace string, name string) error
}

type client struct {
	dynamicClient dynamic.Interface
}

// CreatePytorchJob implements Client
func (c *client) CreatePytorchJob(ptj *kubeflowv1.PyTorchJob) error {
	ptjUpdateTypeMeta(ptj)
	return c.createResource(ptj, ptj.Namespace, ptjRes())
}

// DeletePytorchJob implements Client
func (c *client) DeletePytorchJob(namespace string, name string) error {
	return c.deleteResource(namespace, name, ptjRes())
}

// GetPytorchJob implements Client
func (c *client) GetPytorchJob(namespace string, name string) (*kubeflowv1.PyTorchJob, error) {
	var ptj kubeflowv1.PyTorchJob
	resp, err := c.getResource(namespace, name, ptjRes())
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("[Warning] VirtualMachine %s not found (namespace=%s)", name, namespace)
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get VirtualMachine, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, &ptj); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to VirtualMachine, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	return &ptj, nil
}

// UpdatePytorchJob implements Client
func (c *client) UpdatePytorchJob(namespace string, name string, vm *kubeflowv1.PyTorchJob, data []byte) error {
	ptjUpdateTypeMeta(vm)
	return c.updateResource(namespace, name, ptjRes(), vm, data)
}

func ptjUpdateTypeMeta(vm *kubeflowv1.PyTorchJob) {
	vm.TypeMeta = metav1.TypeMeta{
		Kind:       "PytorchJob",
		APIVersion: kubeflowv1.GroupVersion.String(),
	}
}

func ptjRes() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    kubeflowv1.GroupVersion.Group,
		Version:  kubeflowv1.GroupVersion.Version,
		Resource: "pytorchjobs",
	}

}

// New creates our client wrapper object for the actual kubeVirt and kubernetes clients we use.
func NewClient(cfg *restclient.Config) (Client, error) {
	result := &client{}
	c, err := dynamic.NewForConfig(cfg)
	if err != nil {
		msg := fmt.Sprintf("Failed to create client, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	result.dynamicClient = c
	return result, nil
}

// Generic Resource CRUD operations
func (c *client) createResource(obj interface{}, namespace string, resource schema.GroupVersionResource) error {
	resultMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		msg := fmt.Sprintf("Failed to translate %s to Unstructed (for create operation), with error: %v", resource.Resource, err)
		log.Printf("[Error] %s", msg)
		return fmt.Errorf(msg)
	}
	input := unstructured.Unstructured{}
	input.SetUnstructuredContent(resultMap)
	resp, err := c.dynamicClient.Resource(resource).Namespace(namespace).Create(context.Background(), &input, metav1.CreateOptions{})
	if err != nil {
		msg := fmt.Sprintf("Failed to create %s, with error: %v", resource.Resource, err)
		log.Printf("[Error] %s", msg)
		return fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	return runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, obj)
}

func (c *client) getResource(namespace string, name string, resource schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	return c.dynamicClient.Resource(resource).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (c *client) updateResource(namespace string, name string, resource schema.GroupVersionResource, obj interface{}, data []byte) error {
	resp, err := c.dynamicClient.Resource(resource).Namespace(namespace).Patch(context.Background(), name, pkgApi.JSONPatchType, data, metav1.PatchOptions{})
	if err != nil {
		msg := fmt.Sprintf("Failed to update %s, with error: %v", resource.Resource, err)
		log.Printf("[Error] %s", msg)
		return fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	return runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, obj)
}

func (c *client) deleteResource(namespace string, name string, resource schema.GroupVersionResource) error {
	return c.dynamicClient.Resource(resource).Namespace(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}
