package xgboost_job

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/kubernetes"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
)

func XGBoostJobFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": kubernetes.NamespacedMetadataSchema("XGBoostJob", false),
		"spec":     xgboostJobSpecSchema(),
		"status":   xgboostJobStatusSchema(),
	}
}

func ExpandXGBoostJob(xgboostJobs []interface{}) (*kubeflowv1.XGBoostJob, error) {
	result := &kubeflowv1.XGBoostJob{}

	if len(xgboostJobs) == 0 || xgboostJobs[0] == nil {
		return result, nil
	}

	in := xgboostJobs[0].(map[string]interface{})

	if v, ok := in["metadata"].([]interface{}); ok {
		result.ObjectMeta = kubernetes.ExpandMetadata(v)
	}
	if v, ok := in["spec"].([]interface{}); ok {
		spec, err := expandXGBoostJobSpec(v)
		if err != nil {
			return result, err
		}
		result.Spec = spec
	}
	if v, ok := in["status"].([]interface{}); ok {
		status, err := expandXGBoostJobStatus(v)
		if err != nil {
			return result, err
		}
		result.Status = status
	}

	return result, nil
}

func FlattenXGBoostJob(in kubeflowv1.XGBoostJob) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = kubernetes.FlattenMetadata(in.ObjectMeta)
	att["spec"] = flattenXGBoostJobSpec(in.Spec)
	att["status"] = flattenXGBoostJobStatus(in.Status)

	return []interface{}{att}
}

func FromResourceData(resourceData *schema.ResourceData) (*kubeflowv1.XGBoostJob, error) {
	result := &kubeflowv1.XGBoostJob{}

	result.ObjectMeta = kubernetes.ExpandMetadata(resourceData.Get("metadata").([]interface{}))
	spec, err := expandXGBoostJobSpec(resourceData.Get("spec").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Spec = spec
	status, err := expandXGBoostJobStatus(resourceData.Get("status").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Status = status

	return result, nil
}

func ToResourceData(vm kubeflowv1.XGBoostJob, resourceData *schema.ResourceData) error {
	if err := resourceData.Set("metadata", kubernetes.FlattenMetadata(vm.ObjectMeta)); err != nil {
		return err
	}
	if err := resourceData.Set("spec", flattenXGBoostJobSpec(vm.Spec)); err != nil {
		return err
	}
	if err := resourceData.Set("status", flattenXGBoostJobStatus(vm.Status)); err != nil {
		return err
	}

	return nil
}

func AppendPatchOps(keyPrefix, pathPrefix string, resourceData *schema.ResourceData, ops []patch.PatchOperation) patch.PatchOperations {
	return kubernetes.AppendPatchOps(keyPrefix+"metadata.0.", pathPrefix+"/metadata/", resourceData, ops)
}
