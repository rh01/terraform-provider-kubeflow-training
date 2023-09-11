package mpi_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	corev1 "k8s.io/api/core/v1"

	mpiv2beta1 "github.com/kubeflow/mpi-operator/pkg/apis/kubeflow/v2beta1"
)

func mpiJobConditionsFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "MPIJobConditionType represent the type of the VM as concluded from its VMi status.",
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

func mpiJobConditionsSchema() *schema.Schema {
	fields := mpiJobConditionsFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Hold the state information of the MPIJob and its MPIJobInstance."),
		Required:    true,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandMPIJobConditions(conditions []interface{}) ([]mpiv2beta1.JobCondition, error) {
	result := make([]mpiv2beta1.JobCondition, len(conditions))

	if len(conditions) == 0 || conditions[0] == nil {
		return result, nil
	}

	for i, v := range conditions {
		c := v.(map[string]interface{})
		result[i] = mpiv2beta1.JobCondition{
			Type:    mpiv2beta1.JobConditionType(c["type"].(string)),
			Status:  corev1.ConditionStatus(c["status"].(string)),
			Reason:  c["reason"].(string),
			Message: c["message"].(string),
		}
	}

	return result, nil
}

func flattenMPIJobConditions(in []mpiv2beta1.JobCondition) []interface{} {
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
