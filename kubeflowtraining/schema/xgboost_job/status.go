package xgboost_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
)

func xgboostJobStatusFields() map[string]*schema.Schema {
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
		"conditions": xgboostJobConditionsSchema(),
	}
}

func xgboostJobStatusSchema() *schema.Schema {
	fields := xgboostJobStatusFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("XGBoostJobStatus represents the status returned by the controller to describe how the XGBoostJob is doing."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandXGBoostJobStatus(xgboostJobStatus []interface{}) (commonv1.JobStatus, error) {
	result := commonv1.JobStatus{}

	if len(xgboostJobStatus) == 0 || xgboostJobStatus[0] == nil {
		return result, nil
	}

	// in := xgboostJobStatus[0].(map[string]interface{})

	// if v, ok := in["created"].(bool); ok {
	// 	result.Created = v
	// }
	// if v, ok := in["ready"].(bool); ok {
	// 	result.Ready = v
	// }
	// if v, ok := in["conditions"].([]interface{}); ok {
	// 	conditions, err := expandXGBoostJobConditions(v)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result.Conditions = conditions
	// }
	// if v, ok := in["state_change_requests"].([]interface{}); ok {
	// 	result.StateChangeRequests = expandXGBoostJobStateChangeRequests(v)
	// }

	return result, nil
}

func flattenXGBoostJobStatus(in commonv1.JobStatus) []interface{} {
	att := make(map[string]interface{})

	// att["created"] = in.Created
	// att["ready"] = in.Ready
	// att["conditions"] = flattenXGBoostJobConditions(in.Conditions)
	// att["state_change_requests"] = flattenXGBoostJobStateChangeRequests(in.StateChangeRequests)

	return []interface{}{att}
}
