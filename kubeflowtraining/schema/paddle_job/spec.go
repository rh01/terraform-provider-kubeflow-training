package paddle_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func paddleJobSpecFields() map[string]*schema.Schema {
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
		"paddle_replica_specs": {
			Type:        schema.TypeMap,
			Description: "A map of PaddleReplicaType (type) to ReplicaSpec (value). Specifies the Paddle cluster configuration.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: paddleJobReplicaSpecFields(),
			},
		},
	}
}

func paddleJobSpecSchema() *schema.Schema {
	fields := paddleJobSpecFields()

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

func expandPaddleJobSpec(paddleJob []interface{}) (kubeflowv1.PaddleJobSpec, error) {
	result := kubeflowv1.PaddleJobSpec{}

	if len(paddleJob) == 0 || paddleJob[0] == nil {
		return result, nil
	}

	return result, nil
}

func flattenPaddleJobSpec(in kubeflowv1.PaddleJobSpec) []interface{} {
	att := make(map[string]interface{})

	return []interface{}{att}
}
