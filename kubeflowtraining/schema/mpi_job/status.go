package mpi_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
)

func mpiJobStatusFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"conditions":       mpiJobConditionsSchema(),
		"replica_statuses": mpiJobReplicaStatusesSchema(),
	}
}

func mpiJobReplicaStatusesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("ReplicaStatuses is map of ReplicaType and ReplicaStatus, specifies the status of each replica."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: mpiJobReplicaStatusesFields(),
		},
	}
}

func mpiJobReplicaStatusesFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"active": &schema.Schema{
			Type:        schema.TypeInt,
			Description: fmt.Sprintf("The number of actively running pods."),
			Optional:    true,
		},
		"succeeded": &schema.Schema{
			Type:        schema.TypeInt,
			Description: fmt.Sprintf("The number of pods which reached phase Succeeded."),
			Optional:    true,
		},
		"failed": &schema.Schema{
			Type:        schema.TypeInt,
			Description: fmt.Sprintf("The number of pods which reached phase Failed."),
			Optional:    true,
		},
		"label_selector": &schema.Schema{
			Type:        schema.TypeString,
			Description: fmt.Sprintf("Deprecated: Use Selector instead"),
			Optional:    true,
		},
		"selector": &schema.Schema{
			Type:        schema.TypeString,
			Description: fmt.Sprintf("A Selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty Selector matches all objects. A null Selector matches no objects."),
			Optional:    true,
		},
	}
}

func mpiJobStatusSchema() *schema.Schema {
	fields := mpiJobStatusFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("MPIJobStatus represents the status returned by the controller to describe how the MPIJob is doing."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandMPIJobStatus(mpiJobStatus []interface{}) (commonv1.JobStatus, error) {
	result := commonv1.JobStatus{}

	if len(mpiJobStatus) == 0 || mpiJobStatus[0] == nil {
		return result, nil
	}

	in := mpiJobStatus[0].(map[string]interface{})

	if v, ok := in["conditions"].([]interface{}); ok {
		conditions, err := expandMPIJobConditions(v)
		if err != nil {
			return result, err
		}
		result.Conditions = conditions
	}

	return result, nil
}

func flattenMPIJobStatus(in commonv1.JobStatus) []interface{} {
	att := make(map[string]interface{})

	return []interface{}{att}
}
