package paddle_job

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/kubernetes"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
)

func PaddleJobFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": kubernetes.NamespacedMetadataSchema("PaddleJob", false),
		"spec":     paddleJobSpecSchema(),
		"status":   paddleJobStatusSchema(),
	}
}

func ExpandPaddleJob(paddleJobs []interface{}) (*kubeflowv1.PaddleJob, error) {
	result := &kubeflowv1.PaddleJob{}

	if len(paddleJobs) == 0 || paddleJobs[0] == nil {
		return result, nil
	}

	in := paddleJobs[0].(map[string]interface{})

	if v, ok := in["metadata"].([]interface{}); ok {
		result.ObjectMeta = kubernetes.ExpandMetadata(v)
	}
	if v, ok := in["spec"].([]interface{}); ok {
		spec, err := expandPaddleJobSpec(v)
		if err != nil {
			return result, err
		}
		result.Spec = spec
	}
	if v, ok := in["status"].([]interface{}); ok {
		status, err := expandPaddleJobStatus(v)
		if err != nil {
			return result, err
		}
		result.Status = status
	}

	return result, nil
}

func FlattenPaddleJob(in kubeflowv1.PaddleJob) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = kubernetes.FlattenMetadata(in.ObjectMeta)
	att["spec"] = flattenPaddleJobSpec(in.Spec)
	att["status"] = flattenPaddleJobStatus(in.Status)

	return []interface{}{att}
}

func FromResourceData(resourceData *schema.ResourceData) (*kubeflowv1.PaddleJob, error) {
	result := &kubeflowv1.PaddleJob{}

	result.ObjectMeta = kubernetes.ExpandMetadata(resourceData.Get("metadata").([]interface{}))
	spec, err := expandPaddleJobSpec(resourceData.Get("spec").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Spec = spec
	status, err := expandPaddleJobStatus(resourceData.Get("status").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Status = status

	return result, nil
}

func ToResourceData(vm kubeflowv1.PaddleJob, resourceData *schema.ResourceData) error {
	if err := resourceData.Set("metadata", kubernetes.FlattenMetadata(vm.ObjectMeta)); err != nil {
		return err
	}
	if err := resourceData.Set("spec", flattenPaddleJobSpec(vm.Spec)); err != nil {
		return err
	}
	if err := resourceData.Set("status", flattenPaddleJobStatus(vm.Status)); err != nil {
		return err
	}

	return nil
}

func AppendPatchOps(keyPrefix, pathPrefix string, resourceData *schema.ResourceData, ops []patch.PatchOperation) patch.PatchOperations {
	return kubernetes.AppendPatchOps(keyPrefix+"metadata.0.", pathPrefix+"/metadata/", resourceData, ops)
}
