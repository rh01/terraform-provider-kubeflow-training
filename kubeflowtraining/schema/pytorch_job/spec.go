package pytorch_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func pyTorchJobSpecFields() map[string]*schema.Schema {
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
				Schema: pyTorchJobReplicaSpecFields(),
			},
		},
	}
}

func pyTorchJobSpecSchema() *schema.Schema {
	fields := pyTorchJobSpecFields()

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

func expandPyTorchJobSpec(pyTorchJob []interface{}) (kubeflowv1.PyTorchJobSpec, error) {
	result := kubeflowv1.PyTorchJobSpec{}

	if len(pyTorchJob) == 0 || pyTorchJob[0] == nil {
		return result, nil
	}

	in := pyTorchJob[0].(map[string]interface{})
	if v, ok := in["elastic_policy"]; ok {
		result.ElasticPolicy, _ = expandElasticPolicy(v.([]interface{}))
	}

	return result, nil
}

func flattenPyTorchJobSpec(in kubeflowv1.PyTorchJobSpec) []interface{} {
	att := make(map[string]interface{})

	// if in.Running != nil {
	// 	att["running"] = strconv.FormatBool(*in.Running)
	// }
	if in.ElasticPolicy != nil {
		att["elastic_policy"] = flattenElasticPolicy(in.ElasticPolicy)
	}
	// if in.Template != nil {
	// 	att["template"] = virtualmachineinstance.FlattenVirtualMachineInstanceTemplateSpec(*in.Template)
	// }
	// att["data_volume_templates"] = datavolume.FlattenDataVolumeTemplates(in.DataVolumeTemplates)

	return []interface{}{att}
}
