package kubeflowtraining

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/client"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/paddle_job"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func resourceKubeFlowPaddleJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeFlowPaddleJobCreate,
		Read:   resourceKubeFlowPaddleJobRead,
		Update: resourceKubeFlowPaddleJobUpdate,
		Delete: resourceKubeFlowPaddleJobDelete,
		Exists: resourceKubeFlowPaddleJobExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: paddle_job.PaddleJobFields(),
	}
}

func resourceKubeFlowPaddleJobCreate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	pj, err := paddle_job.FromResourceData(resourceData)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating new PaddleJob: %#v", pj)
	if err := cli.CreatePaddleJob(pj); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new PaddleJob: %#v", pj)
	if err := paddle_job.ToResourceData(*pj, resourceData); err != nil {
		return err
	}
	resourceData.SetId(utils.BuildId(pj.ObjectMeta))

	// Wait for PaddleJob instance's status phase to be succeeded:
	name := pj.ObjectMeta.Name
	namespace := pj.ObjectMeta.Namespace

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Creating"},
		Target:  []string{"Succeeded"},
		Timeout: resourceData.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			var err error
			pj, err = cli.GetPaddleJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("[DEBUG] PaddleJob %s is not created yet", name)
					return pj, "Creating", nil
				}
				return pj, "", err
			}

			for _, c := range pj.Status.Conditions {
				if c.Type == commonv1.JobSucceeded && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PadddleJob %s is succeeded", name)
					return pj, "Succeeded", nil
				}

				if c.Type == commonv1.JobFailed && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PadddleJob %s is failed", name)
					return pj, "Failed", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PadddleJob %s is running", name)
					return pj, "Running", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] PadddleJob %s is pending", name)
					return pj, "Pending", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PadddleJob %s is created", name)
					return pj, "Created", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] PadddleJob %s is creating", name)
					return pj, "Creating", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] PadddleJob %s is restarting", name)
					return pj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] PadddleJob %s is restarting", name)
					return pj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionUnknown {
					log.Printf("[DEBUG] PadddleJob %s is restarting", name)
					return pj, "Restarting", nil
				}
			}

			if pj.Status.StartTime == nil {
				log.Printf("[DEBUG] PadddleJob %s is not started yet", name)
				return pj, "Creating", nil
			}

			log.Printf("[DEBUG] PaddleJob %s is being created", name)
			return pj, "Creating", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}
	return paddle_job.ToResourceData(*pj, resourceData)
}

func resourceKubeFlowPaddleJobRead(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading PaddleJob %s", name)

	pj, err := cli.GetPaddleJob(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received PaddleJob: %#v", pj)

	return paddle_job.ToResourceData(*pj, resourceData)
}

func resourceKubeFlowPaddleJobUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	ops := paddle_job.AppendPatchOps("", "", resourceData, []patch.PatchOperation{})
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating PaddleJob: %s", ops)
	out := &kubeflowv1.PaddleJob{}
	if err := cli.UpdatePaddleJob(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated PaddleJob: %#v", out)

	return resourceKubeFlowPaddleJobRead(resourceData, meta)
}

func resourceKubeFlowPaddleJobDelete(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PaddleJob: %#v", name)
	if err := cli.DeletePaddleJob(namespace, name); err != nil {
		return err
	}

	// Wait for PaddleJob instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: resourceData.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			pj, err := cli.GetPaddleJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil, "", nil
				}
				return pj, "", err
			}

			log.Printf("[DEBUG] PaddleJob %s is being deleted", pj.GetName())
			return pj, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("[INFO] PaddleJob %s deleted", name)

	resourceData.SetId("")
	return nil
}

func resourceKubeFlowPaddleJobExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking PaddleJob %s", name)
	if _, err := cli.GetPaddleJob(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
