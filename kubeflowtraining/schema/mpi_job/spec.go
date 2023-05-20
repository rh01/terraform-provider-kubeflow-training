package mpi_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func mpiJobSpecFields() map[string]*schema.Schema {
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
		"mpi_replica_specs": {
			Type:        schema.TypeMap,
			Description: "A map of MPIReplicaType (type) to ReplicaSpec (value). Specifies the MPI cluster configuration.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: mpiJobReplicaSpecFields(),
			},
		},
		"main_container": {
			Type:        schema.TypeString,
			Description: "MainContainer specifies name of the main container which executes the MPI code.",
			Optional:    true,
		},
		"clean_pod_policy": {
			Type:        schema.TypeString,
			Description: "CleanPodPolicy defines the policy that whether to kill pods after the job completes.",
			Optional:    true,
		},
		"slots_per_worker": {
			Type:        schema.TypeInt,
			Description: "Specifies the number of slots per worker used in hostfile.",
			Optional:    true,
		},
	}
}

func mpiJobSpecSchema() *schema.Schema {
	fields := mpiJobSpecFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("MPIJobSpec describes how the proper MPIJob should look like."),
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandMPIJobSpec(mpiJob []interface{}) (kubeflowv1.MPIJobSpec, error) {
	result := kubeflowv1.MPIJobSpec{}

	if len(mpiJob) == 0 || mpiJob[0] == nil {
		return result, nil
	}

	mpiJobMap := mpiJob[0].(map[string]interface{})
	if mpiJobMap == nil {
		return result, nil
	}

	if v, ok := mpiJobMap["run_policy"]; ok {
		runPolicy, err := expandRunPolicy(v.([]interface{}))
		if err != nil {
			return result, err
		}
		result.RunPolicy = *runPolicy
	}

	if v, ok := mpiJobMap["mpi_replica_specs"]; ok {
		mpiReplicaSpecs, err := expandMPIReplicaSpec(v.([]interface{}))
		if err != nil {
			return result, err
		}
		result.MPIReplicaSpecs = mpiReplicaSpecs
	}

	if v, ok := mpiJobMap["main_container"]; ok {
		result.MainContainer = v.(string)
	}

	if v, ok := mpiJobMap["clean_pod_policy"]; ok {
		*result.CleanPodPolicy = commonv1.CleanPodPolicy(v.(string))
	}

	if v, ok := mpiJobMap["slots_per_worker"]; ok {
		*result.SlotsPerWorker = int32(v.(int))
	}

	return result, nil
}

func flattenMPIJobSpec(in kubeflowv1.MPIJobSpec) []interface{} {
	att := make(map[string]interface{})

	att["run_policy"] = flattenRunPolicy(in.RunPolicy)

	if in.MPIReplicaSpecs != nil {
		att["mpi_replica_specs"], _ = flattenMPIReplicaSpec(in.MPIReplicaSpecs)
	}

	att["main_container"] = in.MainContainer

	if in.CleanPodPolicy != nil {
		att["clean_pod_policy"] = string(*in.CleanPodPolicy)
	}

	if in.SlotsPerWorker != nil {
		att["slots_per_worker"] = int(*in.SlotsPerWorker)
	}

	return []interface{}{att}
}
