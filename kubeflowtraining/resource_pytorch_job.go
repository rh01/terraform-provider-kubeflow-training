package kubeflowtraining

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
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

	log.Printf("[INFO] Creating new data volume: %#v", ptj)
	if err := cli.CreatePyTorchJob(ptj); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new data volume: %#v", ptj)
	if err := pytorch_job.ToResourceData(*ptj, resourceData); err != nil {
		return err
	}
	resourceData.SetId(utils.BuildId(ptj.ObjectMeta))

	// Wait for data volume instance's status phase to be succeeded:
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

			

			// ptj.Status.ReplicaStatuses

			// switch dv.Status.Phase {
			// // case cdiv1.Succeeded:
			// // 	return dv, "Succeeded", nil
			// // case cdiv1.Failed:
			// // 	return dv, "", fmt.Errorf("data volume failed to be created, finished with phase=\"failed\"")
			// }

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

	log.Printf("[INFO] Reading data volume %s", name)

	dv, err := cli.GetPyTorchJob(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received data volume: %#v", dv)

	return pytorch_job.ToResourceData(*dv, resourceData)
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
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating data volume: %s", ops)
	out := &kubeflowv1.PyTorchJob{}
	if err := cli.UpdatePyTorchJob(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated data volume: %#v", out)

	return resourceKubeFlowPyTorchJobRead(resourceData, meta)
}

func resourceKubeFlowPyTorchJobDelete(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting data volume: %#v", name)
	if err := cli.DeletePyTorchJob(namespace, name); err != nil {
		return err
	}

	// Wait for  instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: resourceData.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			dv, err := cli.GetPyTorchJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil, "", nil
				}
				return dv, "", err
			}

			log.Printf("[DEBUG] data volume %s is being deleted", dv.GetName())
			return dv, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("[INFO] data volume %s deleted", name)

	resourceData.SetId("")
	return nil
}

func resourceKubeFlowPyTorchJobExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking data volume %s", name)
	if _, err := cli.GetPyTorchJob(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
