package xgboost_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func xgboostJobSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"run_policy": {
			Type:        schema.TypeList,
			Description: "RunPolicy is a policy for how to run a job.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: runPolicyFields(),
			},
		},
		"elastic_policy": {
			Type:        schema.TypeList,
			Description: "ElasticPolicy is a policy for elastic distributed training.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: elasticPolicyFields(),
			},
		},
		"xgboost_replica_specs": {
			Type:        schema.TypeMap,
			Description: "A map of XGBoostReplicaType (type) to ReplicaSpec (value). Specifies the XGBoost cluster configuration.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: xgboostJobReplicaSpecFields(),
			},
		},
	}
}

func xgboostJobSpecSchema() *schema.Schema {
	fields := xgboostJobSpecFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("VirtualMachineSpec describes how the proper VirtualMachine should look like."),
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandXGBoostJobSpec(xgboostJob []interface{}) (kubeflowv1.XGBoostJobSpec, error) {
	result := kubeflowv1.XGBoostJobSpec{}

	if len(xgboostJob) == 0 || xgboostJob[0] == nil {
		return result, nil
	}

	return result, nil
}

func flattenXGBoostJobSpec(in kubeflowv1.XGBoostJobSpec) []interface{} {
	att := make(map[string]interface{})

	return []interface{}{att}
}
