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
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/xgboost_job"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func resourceKubeFlowXGBoostJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeFlowXGBoostJobCreate,
		Read:   resourceKubeFlowXGBoostJobRead,
		Update: resourceKubeFlowXGBoostJobUpdate,
		Delete: resourceKubeFlowXGBoostJobDelete,
		Exists: resourceKubeFlowXGBoostJobExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: xgboost_job.XGBoostJobFields(),
	}
}

func resourceKubeFlowXGBoostJobCreate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	xgbj, err := xgboost_job.FromResourceData(resourceData)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating new XGBoostJob: %#v", xgbj)
	if err := cli.CreateXGBoostJob(xgbj); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new XGBoostJob: %#v", xgbj)
	if err := xgboost_job.ToResourceData(*xgbj, resourceData); err != nil {
		return err
	}
	resourceData.SetId(utils.BuildId(xgbj.ObjectMeta))

	// Wait for XGBoostJob instance's status phase to be succeeded:
	name := xgbj.ObjectMeta.Name
	namespace := xgbj.ObjectMeta.Namespace

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Creating"},
		Target:  []string{"Succeeded"},
		Timeout: resourceData.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			var err error
			xgbj, err = cli.GetXGBoostJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("[DEBUG] XGBoostJob %s is not created yet", name)
					return xgbj, "Creating", nil
				}
				return xgbj, "", err
			}

			for _, c := range xgbj.Status.Conditions {
				if c.Type == commonv1.JobSucceeded && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] XGBoostJob %s is succeeded", name)
					return xgbj, "Succeeded", nil
				}

				if c.Type == commonv1.JobFailed && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] XGBoostJob %s is failed", name)
					return xgbj, "Failed", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] XGBoostJob %s is running", name)
					return xgbj, "Running", nil
				}

				if c.Type == commonv1.JobRunning && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] XGBoostJob %s is pending", name)
					return xgbj, "Pending", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] XGBoostJob %s is created", name)
					return xgbj, "Created", nil
				}

				if c.Type == commonv1.JobCreated && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] XGBoostJob %s is creating", name)
					return xgbj, "Creating", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionTrue {
					log.Printf("[DEBUG] XGBoostJob %s is restarting", name)
					return xgbj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionFalse {
					log.Printf("[DEBUG] XGBoostJob %s is restarting", name)
					return xgbj, "Restarting", nil
				}

				if c.Type == commonv1.JobRestarting && c.Status == corev1.ConditionUnknown {
					log.Printf("[DEBUG] XGBoostJob %s is restarting", name)
					return xgbj, "Restarting", nil
				}
			}

			if xgbj.Status.StartTime == nil {
				log.Printf("[DEBUG] XGBoostJob %s is not started yet", name)
				return xgbj, "Creating", nil
			}

			log.Printf("[DEBUG] XGBoostJob %s is being created", name)
			return xgbj, "Creating", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}
	return xgboost_job.ToResourceData(*xgbj, resourceData)
}

func resourceKubeFlowXGBoostJobRead(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading XGBoostJob %s", name)

	xgbj, err := cli.GetXGBoostJob(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received XGBoostJob: %#v", xgbj)

	return xgboost_job.ToResourceData(*xgbj, resourceData)
}

func resourceKubeFlowXGBoostJobUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	ops := xgboost_job.AppendPatchOps("", "", resourceData, []patch.PatchOperation{})
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating XGBoostJob: %s", ops)
	out := &kubeflowv1.XGBoostJob{}
	if err := cli.UpdateXGBoostJob(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated XGBoostJob: %#v", out)

	return resourceKubeFlowXGBoostJobRead(resourceData, meta)
}

func resourceKubeFlowXGBoostJobDelete(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting XGBoostJob: %#v", name)
	if err := cli.DeleteXGBoostJob(namespace, name); err != nil {
		return err
	}

	// Wait for XGBoostJob instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: resourceData.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			xgbj, err := cli.GetXGBoostJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil, "", nil
				}
				return xgbj, "", err
			}

			log.Printf("[DEBUG] XGBoostJob %s is being deleted", xgbj.GetName())
			return xgbj, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("[INFO] XGBoostJob %s deleted", name)

	resourceData.SetId("")
	return nil
}

func resourceKubeFlowXGBoostJobExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking XGBoostJob %s", name)
	if _, err := cli.GetXGBoostJob(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
