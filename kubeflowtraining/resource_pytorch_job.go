package kubeflowtraining

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/client"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/pytorch_job"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
	"k8s.io/apimachinery/pkg/api/errors"
)

func resourceKubeFlowPyTorchJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeFlowPyTorchJobCreate,
		Read:   resourceKubeFlowPyTorchJobRead,
		Update: resourceKubeFlowPyTorchJobUpdate,
		Delete: resourceKubeFlowPyTorchJobDelete,
		Exists: resourceKubeFlowPyTorchJobExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: pytorch_job.PyTorchJobFields(),
	}
}

func resourceKubeFlowPyTorchJobCreate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	ptj, err := pytorch_job.FromResourceData(resourceData)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating new PyTorchJob: %#v", ptj)
	if err := cli.CreatePyTorchJob(ptj); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new PyTorchJob: %#v", ptj)
	if err := pytorch_job.ToResourceData(*ptj, resourceData); err != nil {
		return err
	}
	resourceData.SetId(utils.BuildId(ptj.ObjectMeta))

	// Wait for PyTorchJob instance's status phase to be succeeded:
	name := ptj.ObjectMeta.Name
	namespace := ptj.ObjectMeta.Namespace

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Creating"},
		Target:  []string{"Succeeded"},
		Timeout: resourceData.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			var err error
			ptj, err = cli.GetPyTorchJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("[DEBUG] PyTorchJob %s is not created yet", name)
					return ptj, "Creating", nil
				}
				return ptj, "", err
			}

			if err = kubeflowv1.ValidateV1PyTorchJob(ptj); err != nil {
				log.Printf("[DEBUG] PyTorchJob %s is not valid yet: %s", name, err)
				return ptj, "Creating", nil
			}

			for _, c := range ptj.Status.Conditions {
				if c.Type == commonv1.JobSucceeded && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PyTorchJob %s is succeeded", name)
					return ptj, "Succeeded", nil
				}

				if c.Type == commonv1.JobFailed && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PyTorchJob %s is failed", name)
					return ptj, "Failed", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PyTorchJob %s is running", name)
					return ptj, "Running", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] PyTorchJob %s is pending", name)
					return ptj, "Pending", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PyTorchJob %s is created", name)
					return ptj, "Created", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] PyTorchJob %s is creating", name)
					return ptj, "Creating", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PyTorchJob %s is restarting", name)
					return ptj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] PyTorchJob %s is restarting", name)
					return ptj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionUnknown {
					log.Printf("[DEBUG] PyTorchJob %s is restarting", name)
					return ptj, "Restarting", nil
				}
			}

			if ptj.Status.StartTime == nil {
				log.Printf("[DEBUG] PyTorchJob %s is not started yet", name)
				return ptj, "Creating", nil
			}

			log.Printf("[DEBUG] PyTorchJob %s is being created", name)
			return ptj, "Creating", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}
	return pytorch_job.ToResourceData(*ptj, resourceData)
}

func resourceKubeFlowPyTorchJobRead(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PyTorchJob %s", name)

	ptj, err := cli.GetPyTorchJob(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received PyTorchJob: %#v", ptj)

	return pytorch_job.ToResourceData(*ptj, resourceData)
}

func resourceKubeFlowPyTorchJobUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	ops := pytorch_job.AppendPatchOps("", "", resourceData, []patch.PatchOperation{})
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating PyTorchJob: %s", ops)
	out := &kubeflowv1.PyTorchJob{}
	if err := cli.UpdatePyTorchJob(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated PyTorchJob: %#v", out)

	return resourceKubeFlowPyTorchJobRead(resourceData, meta)
}

func resourceKubeFlowPyTorchJobDelete(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PyTorchJob: %#v", name)
	if err := cli.DeletePyTorchJob(namespace, name); err != nil {
		return err
	}

	// Wait for  instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: resourceData.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			ptj, err := cli.GetPyTorchJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil, "", nil
				}
				return ptj, "", err
			}

			log.Printf("[DEBUG] PyTorchJob %s is being deleted", ptj.GetName())
			return ptj, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("[INFO] PyTorchJob %s deleted", name)

	resourceData.SetId("")
	return nil
}

func resourceKubeFlowPyTorchJobExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking PyTorchJob %s", name)
	if _, err := cli.GetPyTorchJob(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
