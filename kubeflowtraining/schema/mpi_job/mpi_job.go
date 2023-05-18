package mpi_job

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/kubernetes"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
)

func MPIJobFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": kubernetes.NamespacedMetadataSchema("MPIJob", false),
		"spec":     mpiJobSpecSchema(),
		"status":   mpiJobStatusSchema(),
	}
}

func ExpandMPIJob(mpiJobs []interface{}) (*kubeflowv1.MPIJob, error) {
	result := &kubeflowv1.MPIJob{}

	if len(mpiJobs) == 0 || mpiJobs[0] == nil {
		return result, nil
	}

	in := mpiJobs[0].(map[string]interface{})

	if v, ok := in["metadata"].([]interface{}); ok {
		result.ObjectMeta = kubernetes.ExpandMetadata(v)
	}
	if v, ok := in["spec"].([]interface{}); ok {
		spec, err := expandMPIJobSpec(v)
		if err != nil {
			return result, err
		}
		result.Spec = spec
	}
	if v, ok := in["status"].([]interface{}); ok {
		status, err := expandMPIJobStatus(v)
		if err != nil {
			return result, err
		}
		result.Status = status
	}

	return result, nil
}

func FlattenMPIJob(in kubeflowv1.MPIJob) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = kubernetes.FlattenMetadata(in.ObjectMeta)
	att["spec"] = flattenMPIJobSpec(in.Spec)
	att["status"] = flattenVirtualMachineStatus(in.Status)

	return []interface{}{att}
}

func FromResourceData(resourceData *schema.ResourceData) (*kubeflowv1.MPIJob, error) {
	result := &kubeflowv1.MPIJob{}

	result.ObjectMeta = kubernetes.ExpandMetadata(resourceData.Get("metadata").([]interface{}))
	spec, err := expandMPIJobSpec(resourceData.Get("spec").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Spec = spec
	status, err := expandMPIJobStatus(resourceData.Get("status").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Status = status

	return result, nil
}

func ToResourceData(vm kubeflowv1.MPIJob, resourceData *schema.ResourceData) error {
	if err := resourceData.Set("metadata", kubernetes.FlattenMetadata(vm.ObjectMeta)); err != nil {
		return err
	}
	if err := resourceData.Set("spec", flattenMPIJobSpec(vm.Spec)); err != nil {
		return err
	}
	if err := resourceData.Set("status", flattenVirtualMachineStatus(vm.Status)); err != nil {
		return err
	}

	return nil
}

func AppendPatchOps(keyPrefix, pathPrefix string, resourceData *schema.ResourceData, ops []patch.PatchOperation) patch.PatchOperations {
	return kubernetes.AppendPatchOps(keyPrefix+"metadata.0.", pathPrefix+"/metadata/", resourceData, ops)
}
