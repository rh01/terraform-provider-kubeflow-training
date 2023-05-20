package tensorflow_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
)

func tfJobStatusFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created": &schema.Schema{
			Type:        schema.TypeBool,
			Description: "Created indicates if the virtual machine is created in the cluster.",
			Optional:    true,
		},
		"ready": &schema.Schema{
			Type:        schema.TypeBool,
			Description: "Ready indicates if the virtual machine is running and ready.",
			Optional:    true,
		},
		"conditions": tfJobConditionsSchema(),
	}
}

func tfJobStatusSchema() *schema.Schema {
	fields := tfJobStatusFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("TFJobStatus represents the status returned by the controller to describe how the TFJob is doing."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandTFJobStatus(tfJobStatus []interface{}) (commonv1.JobStatus, error) {
	result := commonv1.JobStatus{}

	if len(tfJobStatus) == 0 || tfJobStatus[0] == nil {
		return result, nil
	}

	// in := tfJobStatus[0].(map[string]interface{})

	// if v, ok := in["created"].(bool); ok {
	// 	result.Created = v
	// }
	// if v, ok := in["ready"].(bool); ok {
	// 	result.Ready = v
	// }
	// if v, ok := in["conditions"].([]interface{}); ok {
	// 	conditions, err := expandTFJobConditions(v)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result.Conditions = conditions
	// }
	// if v, ok := in["state_change_requests"].([]interface{}); ok {
	// 	result.StateChangeRequests = expandTFJobStateChangeRequests(v)
	// }

	return result, nil
}

func flattenTFJobStatus(in commonv1.JobStatus) []interface{} {
	att := make(map[string]interface{})

	// att["created"] = in.Created
	// att["ready"] = in.Ready
	// att["conditions"] = flattenTFJobConditions(in.Conditions)
	// att["state_change_requests"] = flattenTFJobStateChangeRequests(in.StateChangeRequests)

	return []interface{}{att}
}
