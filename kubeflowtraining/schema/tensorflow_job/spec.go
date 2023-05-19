package tensorflow_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func tfJobSpecFields() map[string]*schema.Schema {
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
		"tensorflow_replica_specs": {
			Type:        schema.TypeMap,
			Description: "A map of TensorflowReplicaType (type) to ReplicaSpec (value). Specifies the Tensorflow cluster configuration.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: tfJobReplicaSpecFields(),
			},
		},
	}
}

func tfJobSpecSchema() *schema.Schema {
	fields := tfJobSpecFields()

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

func expandTFJobSpec(tfJob []interface{}) (kubeflowv1.TFJobSpec, error) {
	result := kubeflowv1.TFJobSpec{}

	if len(tfJob) == 0 || tfJob[0] == nil {
		return result, nil
	}

	return result, nil
}

func flattenTFJobSpec(in kubeflowv1.TFJobSpec) []interface{} {
	att := make(map[string]interface{})

	return []interface{}{att}
}
