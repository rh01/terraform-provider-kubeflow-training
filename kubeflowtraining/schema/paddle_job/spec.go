package paddle_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func paddleJobSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"run_policy": {
			Type:        schema.TypeString,
			Description: "RunPolicy encapsulates various runtime policies of the distributed training job, for example how to clean up resources and how long the job can stay active.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"AutoDelete",
				"LongRunning",
			}, false),
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
		"pytorch_replica_specs": {
			Type:        schema.TypeMap,
			Description: "A map of PyTorchReplicaType (type) to ReplicaSpec (value). Specifies the PyTorch cluster configuration.",
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
