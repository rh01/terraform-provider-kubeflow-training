package paddle_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
)

func paddleJobConditionsFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "VirtualMachineConditionType represent the type of the VM as concluded from its VMi status.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Failure",
				"Ready",
				"Paused",
				"RenameOperation",
			}, false),
		},
		"status": {
			Type:        schema.TypeString,
			Description: "ConditionStatus represents the status of this VM condition, if the VM currently in the condition.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"True",
				"False",
				"Unknown",
			}, false),
		},

		"reason": {
			Type:        schema.TypeString,
			Description: "Condition reason.",
			Optional:    true,
		},
		"message": {
			Type:        schema.TypeString,
			Description: "Condition message.",
			Optional:    true,
		},
	}
}

func paddleJobConditionsSchema() *schema.Schema {
	fields := paddleJobConditionsFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Hold the state information of the VirtualMachine and its VirtualMachineInstance."),
		Required:    true,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandPaddleJobConditions(conditions []interface{}) ([]commonv1.JobCondition, error) {
	result := make([]commonv1.JobCondition, len(conditions))

	if len(conditions) == 0 || conditions[0] == nil {
		return result, nil
	}

	// for i, condition := range conditions {
	// 	// in := condition.(map[string]interface{})

	// }

	return result, nil
}

func flattenPaddleJobConditions(in []commonv1.JobCondition) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})
		c["type"] = string(v.Type)
		c["status"] = string(v.Status)
		c["reason"] = v.Reason
		c["message"] = v.Message

		att[i] = c
	}

	return att
}
