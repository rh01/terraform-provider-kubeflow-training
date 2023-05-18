package tensorflow_job

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/kubernetes"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
)

func TFJobFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": kubernetes.NamespacedMetadataSchema("TFJob", false),
		"spec":     tfJobSpecSchema(),
		"status":   tfJobStatusSchema(),
	}
}

func ExpandTFJob(tfJobs []interface{}) (*kubeflowv1.TFJob, error) {
	result := &kubeflowv1.TFJob{}

	if len(tfJobs) == 0 || tfJobs[0] == nil {
		return result, nil
	}

	in := tfJobs[0].(map[string]interface{})

	if v, ok := in["metadata"].([]interface{}); ok {
		result.ObjectMeta = kubernetes.ExpandMetadata(v)
	}
	if v, ok := in["spec"].([]interface{}); ok {
		spec, err := expandTFJobSpec(v)
		if err != nil {
			return result, err
		}
		result.Spec = spec
	}
	if v, ok := in["status"].([]interface{}); ok {
		status, err := expandTFJobStatus(v)
		if err != nil {
			return result, err
		}
		result.Status = status
	}

	return result, nil
}

func FlattenTFJob(in kubeflowv1.TFJob) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = kubernetes.FlattenMetadata(in.ObjectMeta)
	att["spec"] = flattenTFJobSpec(in.Spec)
	att["status"] = flattenVirtualMachineStatus(in.Status)

	return []interface{}{att}
}

func FromResourceData(resourceData *schema.ResourceData) (*kubeflowv1.TFJob, error) {
	result := &kubeflowv1.TFJob{}

	result.ObjectMeta = kubernetes.ExpandMetadata(resourceData.Get("metadata").([]interface{}))
	spec, err := expandTFJobSpec(resourceData.Get("spec").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Spec = spec
	status, err := expandTFJobStatus(resourceData.Get("status").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Status = status

	return result, nil
}

func ToResourceData(vm kubeflowv1.TFJob, resourceData *schema.ResourceData) error {
	if err := resourceData.Set("metadata", kubernetes.FlattenMetadata(vm.ObjectMeta)); err != nil {
		return err
	}
	if err := resourceData.Set("spec", flattenTFJobSpec(vm.Spec)); err != nil {
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
