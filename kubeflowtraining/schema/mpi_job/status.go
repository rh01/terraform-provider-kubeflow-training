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
func expandMPIJobReplicaStatuses(in []interface{}) (map[commonv1.ReplicaType]*commonv1.ReplicaStatus, error) {
	result := make(map[commonv1.ReplicaType]*commonv1.ReplicaStatus)

	if len(in) == 0 || in[0] == nil {
		return result, nil
	}

	for _, v := range in {
		replicaStatus := &commonv1.ReplicaStatus{}
		if err := expandMPIJobReplicaStatus(v, replicaStatus); err != nil {
			return result, err
		}
	}

	return result, nil
}

func flattenMPIJobStatus(in commonv1.JobStatus) []interface{} {
	att := make(map[string]interface{})

	att["conditions"] = flattenMPIJobConditions(in.Conditions)

	att["replica_statuses"] = flattenMPIJobReplicaStatuses(in.ReplicaStatuses)

	return []interface{}{att}

}

func flattenMPIJobReplicaStatuses(in map[commonv1.ReplicaType]*commonv1.ReplicaStatus) []interface{} {
	result := make([]interface{}, 0)

	for _, v := range in {
		result = append(result, flattenMPIJobReplicaStatus(v))
	}

	return result
}

func mpiJobReplicaStatusesFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"active": {
			Type:        schema.TypeInt,
			Description: fmt.Sprintf("The number of actively running pods."),
			Optional:    true,
		},
		"succeeded": {
			Type:        schema.TypeInt,
			Description: fmt.Sprintf("The number of pods which reached phase Succeeded."),
			Optional:    true,
		},
		"failed": {
			Type:        schema.TypeInt,
			Description: fmt.Sprintf("The number of pods which reached phase Failed."),
			Optional:    true,
		},
		"label_selector": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("Deprecated: Use Selector instead"),
			Optional:    true,
		},
		"selector": {
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

func flattenMPIJobReplicaStatus(in *commonv1.ReplicaStatus) map[string]interface{} {
	att := make(map[string]interface{})

	att["active"] = in.Active

	att["succeeded"] = in.Succeeded

	att["failed"] = in.Failed

	att["selector"] = in.Selector

	return att
}

func expandMPIJobReplicaStatus(in interface{}, out *commonv1.ReplicaStatus) error {

	if in == nil {
		return nil
	}

	replicaStatus := in.(map[string]interface{})

	if v, ok := replicaStatus["active"].(int); ok {
		out.Active = int32(v)
	}

	if v, ok := replicaStatus["succeeded"].(int); ok {
		out.Succeeded = int32(v)
	}

	if v, ok := replicaStatus["failed"].(int); ok {
		out.Failed = int32(v)
	}

	if v, ok := replicaStatus["selector"].(string); ok {
		out.Selector = v
	}

	return nil
}
