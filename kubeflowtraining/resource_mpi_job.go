package kubeflowtraining

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/client"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/mpi_job"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
	"k8s.io/apimachinery/pkg/api/errors"
)

func resourceKubeFlowMPIJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubeFlowMPIJobCreate,
		Read:   resourceKubeFlowMPIJobRead,
		Update: resourceKubeFlowMPIJobUpdate,
		Delete: resourceKubeFlowMPIJobDelete,
		Exists: resourceKubeFlowMPIJobExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: mpi_job.MPIJobFields(),
	}
}

func resourceKubeFlowMPIJobCreate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	dv, err := mpi_job.FromResourceData(resourceData)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating new data volume: %#v", dv)
	if err := cli.CreateMPIJob(dv); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new data volume: %#v", dv)
	if err := mpi_job.ToResourceData(*dv, resourceData); err != nil {
		return err
	}
	resourceData.SetId(utils.BuildId(dv.ObjectMeta))

	// Wait for data volume instance's status phase to be succeeded:
	name := dv.ObjectMeta.Name
	namespace := dv.ObjectMeta.Namespace

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Creating"},
		Target:  []string{"Succeeded"},
		Timeout: resourceData.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			var err error
			dv, err = cli.GetMPIJob(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("[DEBUG] data volume %s is not created yet", name)
					return dv, "Creating", nil
				}
				return dv, "", err
			}

			// switch dv.Status.Phase {
			// case cdiv1.Succeeded:
			// 	return dv, "Succeeded", nil
			// case cdiv1.Failed:
			// 	return dv, "", fmt.Errorf("data volume failed to be created, finished with phase=\"failed\"")
			// }

			log.Printf("[DEBUG] data volume %s is being created", name)
			return dv, "Creating", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}
	return mpi_job.ToResourceData(*dv, resourceData)
}

func resourceKubeFlowMPIJobRead(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading data volume %s", name)

	dv, err := cli.GetMPIJob(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received data volume: %#v", dv)

	return mpi_job.ToResourceData(*dv, resourceData)
}

func resourceKubeFlowMPIJobUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	ops := mpi_job.AppendPatchOps("", "", resourceData, []patch.PatchOperation{})
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating data volume: %s", ops)
	out := &kubeflowv1.MPIJob{}
	if err := cli.UpdateMPIJob(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated data volume: %#v", out)

	return resourceKubeFlowMPIJobRead(resourceData, meta)
}

func resourceKubeFlowMPIJobDelete(resourceData *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting data volume: %#v", name)
	if err := cli.DeleteMPIJob(namespace, name); err != nil {
		return err
	}

	// Wait for data volume instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: resourceData.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			dv, err := cli.GetMPIJob(namespace, name)
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

func resourceKubeFlowMPIJobExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(resourceData.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking data volume %s", name)
	if _, err := cli.GetMPIJob(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
