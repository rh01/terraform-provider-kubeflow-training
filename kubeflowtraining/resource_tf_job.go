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
	tf_job "github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/tensorflow_job"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func resourceKubeFlowTFJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeFlowTFJobCreate,
		Read:   resourceKubeFlowTFJobRead,
		Update: resourceKubeFlowTFJobUpdate,
		Delete: resourceKubeFlowTFJobDelete,
		Exists: resourceKubeFlowTFJobExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: tf_job.TFJobFields(),
	}
}

func resourceKubeFlowTFJobCreate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	tfj, err := tf_job.FromResourceData(resourceData)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating new TFJob: %#v", tfj)
	if err := cli.CreateTFJob(tfj); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new TFJob: %#v", tfj)
	if err := tf_job.ToResourceData(*tfj, resourceData); err != nil {
		return err
	}
	resourceData.SetId(utils.BuildId(tfj.ObjectMeta))

	// Wait for TFJob instance's status phase to be succeeded:
	name := tfj.ObjectMeta.Name
	namespace := tfj.ObjectMeta.Namespace

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Creating"},
		Target:  []string{"Succeeded"},
		Timeout: resourceData.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			var err error
			tfj, err = cli.GetTFJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("[DEBUG] TFJob %s is not created yet", name)
					return tfj, "Creating", nil
				}
				return tfj, "", err
			}

			for _, c := range tfj.Status.Conditions {
				if c.Type == commonv1.JobSucceeded && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] TFJob %s is succeeded", name)
					return tfj, "Succeeded", nil
				}

				if c.Type == commonv1.JobFailed && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] TFJob %s is failed", name)
					return tfj, "Failed", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] TFJob %s is running", name)
					return tfj, "Running", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] TFJob %s is pending", name)
					return tfj, "Pending", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] TFJob %s is created", name)
					return tfj, "Created", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] TFJob %s is creating", name)
					return tfj, "Creating", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] TFJob %s is restarting", name)
					return tfj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] TFJob %s is restarting", name)
					return tfj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionUnknown {
					log.Printf("[DEBUG] TFJob %s is restarting", name)
					return tfj, "Restarting", nil
				}
			}

			if tfj.Status.StartTime == nil {
				log.Printf("[DEBUG] TFJob %s is not started yet", name)
				return tfj, "Creating", nil
			}

			log.Printf("[DEBUG] TFJob %s is being created", name)
			return tfj, "Creating", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}
	return tf_job.ToResourceData(*tfj, resourceData)
}

func resourceKubeFlowTFJobRead(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading TFJob %s", name)

	tfj, err := cli.GetTFJob(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received TFJob: %#v", tfj)

	return tf_job.ToResourceData(*tfj, resourceData)
}

func resourceKubeFlowTFJobUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	ops := tf_job.AppendPatchOps("", "", resourceData, []patch.PatchOperation{})
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating TFJob: %s", ops)
	out := &kubeflowv1.TFJob{}
	if err := cli.UpdateTFJob(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated TFJob: %#v", out)

	return resourceKubeFlowTFJobRead(resourceData, meta)
}

func resourceKubeFlowTFJobDelete(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting TFJob: %#v", name)
	if err := cli.DeleteTFJob(namespace, name); err != nil {
		return err
	}

	// Wait for TFJob instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: resourceData.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			tfj, err := cli.GetTFJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil, "", nil
				}
				return tfj, "", err
			}

			log.Printf("[DEBUG] TFJob %s is being deleted", tfj.GetName())
			return tfj, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("[INFO] TFJob %s deleted", name)

	resourceData.SetId("")
	return nil
}

func resourceKubeFlowTFJobExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking TFJob %s", name)
	if _, err := cli.GetTFJob(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
