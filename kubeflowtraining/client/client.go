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
	CreatePyTorchJob(job *kubeflowv1.PyTorchJob) error
	GetPyTorchJob(namespace string, name string) (*kubeflowv1.PyTorchJob, error)
	UpdatePyTorchJob(namespace string, name string, job *kubeflowv1.PyTorchJob, data []byte) error
	DeletePyTorchJob(namespace string, name string) error

	// generate TFJob, MPIJob, XGBoostJob, PaddleJob CRUD
	CreateTFJob(job *kubeflowv1.TFJob) error
	GetTFJob(namespace string, name string) (*kubeflowv1.TFJob, error)
	UpdateTFJob(namespace string, name string, job *kubeflowv1.TFJob, data []byte) error
	DeleteTFJob(namespace string, name string) error

	CreateMPIJob(job *kubeflowv1.MPIJob) error
	GetMPIJob(namespace string, name string) (*kubeflowv1.MPIJob, error)
	UpdateMPIJob(namespace string, name string, job *kubeflowv1.MPIJob, data []byte) error
	DeleteMPIJob(namespace string, name string) error

	CreateXGBoostJob(job *kubeflowv1.XGBoostJob) error
	GetXGBoostJob(namespace string, name string) (*kubeflowv1.XGBoostJob, error)
	UpdateXGBoostJob(namespace string, name string, job *kubeflowv1.XGBoostJob, data []byte) error
	DeleteXGBoostJob(namespace string, name string) error

	CreatePaddleJob(job *kubeflowv1.PaddleJob) error
	GetPaddleJob(namespace string, name string) (*kubeflowv1.PaddleJob, error)
	UpdatePaddleJob(namespace string, name string, job *kubeflowv1.PaddleJob, data []byte) error
	DeletePaddleJob(namespace string, name string) error
}

type client struct {
	dynamicClient dynamic.Interface
}

// CreateMPIJob implements Client
func (c *client) CreateMPIJob(mpij *kubeflowv1.MPIJob) error {
	mpijUpdateTypeMeta(mpij)
	return c.createResource(mpij, mpij.Namespace, mpijRes())
}

// CreatePaddleJob implements Client
func (c *client) CreatePaddleJob(pj *kubeflowv1.PaddleJob) error {
	pjUpdateTypeMeta(pj)
	return c.createResource(pj, pj.Namespace, pjRes())
}

// CreateTFJob implements Client
func (c *client) CreateTFJob(tfj *kubeflowv1.TFJob) error {
	tfjUpdateTypeMeta(tfj)
	return c.createResource(tfj, tfj.Namespace, tfjRes())
}

// CreateXGBoostJob implements Client
func (c *client) CreateXGBoostJob(xgbj *kubeflowv1.XGBoostJob) error {
	xgbjUpdateTypeMeta(xgbj)
	return c.createResource(xgbj, xgbj.Namespace, xgbjRes())
}

// DeleteMPIJob implements Client
func (c *client) DeleteMPIJob(namespace string, name string) error {
	return c.deleteResource(namespace, name, mpijRes())
}

// DeletePaddleJob implements Client
func (c *client) DeletePaddleJob(namespace string, name string) error {
	return c.deleteResource(namespace, name, pjRes())
}

// DeleteTFJob implements Client
func (c *client) DeleteTFJob(namespace string, name string) error {
	return c.deleteResource(namespace, name, tfjRes())
}

// DeleteXGBoostJob implements Client
func (c *client) DeleteXGBoostJob(namespace string, name string) error {
	return c.deleteResource(namespace, name, xgbjRes())
}

// GetMPIJob implements Client
func (c *client) GetMPIJob(namespace string, name string) (*kubeflowv1.MPIJob, error) {
	var mpij kubeflowv1.MPIJob
	resp, err := c.getResource(namespace, name, mpijRes())
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("[Warning] MPIJob %s not found (namespace=%s)", name, namespace)
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get MPIJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, &mpij); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to MPIJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	return &mpij, nil
}

// GetPaddleJob implements Client
func (c *client) GetPaddleJob(namespace string, name string) (*kubeflowv1.PaddleJob, error) {
	var pj kubeflowv1.PaddleJob
	resp, err := c.getResource(namespace, name, pjRes())
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("[Warning] PaddleJob %s not found (namespace=%s)", name, namespace)
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get PaddleJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, &pj); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to PaddleJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	return &pj, nil
}

// GetTFJob implements Client
func (c *client) GetTFJob(namespace string, name string) (*kubeflowv1.TFJob, error) {
	var tfj kubeflowv1.TFJob
	resp, err := c.getResource(namespace, name, tfjRes())
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("[Warning] TFJob %s not found (namespace=%s)", name, namespace)
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get TFJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, &tfj); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to TFJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	return &tfj, nil
}

// GetXGBoostJob implements Client
func (c *client) GetXGBoostJob(namespace string, name string) (*kubeflowv1.XGBoostJob, error) {
	var xgbj kubeflowv1.XGBoostJob
	resp, err := c.getResource(namespace, name, xgbjRes())
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("[Warning] XGBoostJob %s not found (namespace=%s)", name, namespace)
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get XGBoostJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, &xgbj); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to XGBoostJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	return &xgbj, nil
}

// UpdateMPIJob implements Client
func (c *client) UpdateMPIJob(namespace string, name string, job *kubeflowv1.MPIJob, data []byte) error {
	mpijUpdateTypeMeta(job)
	return c.updateResource(namespace, name, mpijRes(), job, data)
}

// UpdatePaddleJob implements Client
func (c *client) UpdatePaddleJob(namespace string, name string, job *kubeflowv1.PaddleJob, data []byte) error {
	pjUpdateTypeMeta(job)
	return c.updateResource(namespace, name, pjRes(), job, data)
}

// UpdateTFJob implements Client
func (c *client) UpdateTFJob(namespace string, name string, job *kubeflowv1.TFJob, data []byte) error {
	tfjUpdateTypeMeta(job)
	return c.updateResource(namespace, name, tfjRes(), job, data)
}

// UpdateXGBoostJob implements Client
func (c *client) UpdateXGBoostJob(namespace string, name string, job *kubeflowv1.XGBoostJob, data []byte) error {
	xgbjUpdateTypeMeta(job)
	return c.updateResource(namespace, name, xgbjRes(), job, data)
}

// CreatePyTorchJob implements Client
func (c *client) CreatePyTorchJob(ptj *kubeflowv1.PyTorchJob) error {
	ptjUpdateTypeMeta(ptj)
	return c.createResource(ptj, ptj.Namespace, ptjRes())
}

// DeletePyTorchJob implements Client
func (c *client) DeletePyTorchJob(namespace string, name string) error {
	return c.deleteResource(namespace, name, ptjRes())
}

// GetPyTorchJob implements Client
func (c *client) GetPyTorchJob(namespace string, name string) (*kubeflowv1.PyTorchJob, error) {
	var ptj kubeflowv1.PyTorchJob
	resp, err := c.getResource(namespace, name, ptjRes())
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("[Warning] PyTorchJob %s not found (namespace=%s)", name, namespace)
			return nil, err
		}
		msg := fmt.Sprintf("Failed to get PyTorchJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	unstructured := resp.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, &ptj); err != nil {
		msg := fmt.Sprintf("Failed to translate unstructed to PyTorchJob, with error: %v", err)
		log.Printf("[Error] %s", msg)
		return nil, fmt.Errorf(msg)
	}
	return &ptj, nil
}

// UpdatePyTorchJob implements Client
func (c *client) UpdatePyTorchJob(namespace string, name string, job *kubeflowv1.PyTorchJob, data []byte) error {
	ptjUpdateTypeMeta(job)
	return c.updateResource(namespace, name, ptjRes(), job, data)
}

func ptjUpdateTypeMeta(job *kubeflowv1.PyTorchJob) {
	job.TypeMeta = metav1.TypeMeta{
		Kind:       "PyTorchJob",
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

func tfjUpdateTypeMeta(job *kubeflowv1.TFJob) {
	job.TypeMeta = metav1.TypeMeta{
		Kind:       "TFJob",
		APIVersion: kubeflowv1.GroupVersion.String(),
	}
}

func tfjRes() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    kubeflowv1.GroupVersion.Group,
		Version:  kubeflowv1.GroupVersion.Version,
		Resource: "tfjobs",
	}
}

func mpijUpdateTypeMeta(job *kubeflowv1.MPIJob) {
	job.TypeMeta = metav1.TypeMeta{
		Kind:       "MPIJob",
		APIVersion: kubeflowv1.GroupVersion.String(),
	}
}

func mpijRes() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    kubeflowv1.GroupVersion.Group,
		Version:  kubeflowv1.GroupVersion.Version,
		Resource: "mpijobs",
	}
}

func xgbjUpdateTypeMeta(job *kubeflowv1.XGBoostJob) {
	job.TypeMeta = metav1.TypeMeta{
		Kind:       "XGBoostJob",
		APIVersion: kubeflowv1.GroupVersion.String(),
	}
}

func xgbjRes() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    kubeflowv1.GroupVersion.Group,
		Version:  kubeflowv1.GroupVersion.Version,
		Resource: "xgboostjobs",
	}
}

func pjUpdateTypeMeta(job *kubeflowv1.PaddleJob) {
	job.TypeMeta = metav1.TypeMeta{
		Kind:       "PaddleJob",
		APIVersion: kubeflowv1.GroupVersion.String(),
	}
}

func pjRes() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    kubeflowv1.GroupVersion.Group,
		Version:  kubeflowv1.GroupVersion.Version,
		Resource: "paddlejobs",
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
