package mpi_job

import (
	// "github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/k8s"
	// "github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/utils/patch"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/kubernetes"
)

func mpiJobReplicaSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"launcher": mpiJobReplicaSpecSchema(),
		"worker":   mpiJobReplicaSpecSchema(),
	}
}

func mpiJobReplicaSpecTemplateFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"replicas": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"template": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: kubernetes.PodTemplateFields("mpijob"),
			},
			Optional: true,
		},
		"restart_policy": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "Never",
		},
	}
}

func mpiJobReplicaSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: mpiJobReplicaSpecTemplateFields(),
		},
		Optional: true,
	}
}

func expandMPIReplicaSpec(l []interface{}) (map[commonv1.ReplicaType]*commonv1.ReplicaSpec, error) {
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}

	m := make(map[commonv1.ReplicaType]*commonv1.ReplicaSpec)
	for k, v := range l[0].(map[string]interface{}) {
		if !strings.EqualFold(k, "launcher") && !strings.EqualFold(k, "worker") {
			continue
		}
		if k == "launcher" {
			k = "Launcher"
		} else {
			k = "Worker"
		}
		replicaType := commonv1.ReplicaType(k)
		replicaSpec, err := expandMPIJobReplicaSpec(v.([]interface{}))
		if err != nil {
			return nil, err
		}

		m[replicaType] = replicaSpec
	}
	return m, nil
}

func flattenMPIReplicaSpec(in map[commonv1.ReplicaType]*commonv1.ReplicaSpec) ([]interface{}, error) {
	if in == nil {
		return []interface{}{nil}, nil
	}

	m := make(map[string]interface{})
	for k, v := range in {
		replicaType := string(k)
		replicaSpec, err := flattenMPIJobReplicaSpec(v)
		if err != nil {
			return nil, err
		}
		if replicaType == "Launcher" {
			replicaType = "launcher"
		} else {
			replicaType = "worker"
		}
		m[replicaType] = replicaSpec
	}
	return []interface{}{m}, nil
}

func expandMPIJobReplicaSpec(l []interface{}) (*commonv1.ReplicaSpec, error) {
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	m := l[0].(map[string]interface{})

	replicas := m["replicas"].(*int32)
	template, err := kubernetes.ExpandPodTemplate(m["template"].([]interface{}))
	if err != nil {
		return nil, err
	}
	restartPolicy := m["restart_policy"].(string)

	return &commonv1.ReplicaSpec{
		Replicas:      replicas,
		Template:      *template,
		RestartPolicy: commonv1.RestartPolicy(restartPolicy),
	}, nil
}

func flattenMPIJobReplicaSpec(in *commonv1.ReplicaSpec) ([]interface{}, error) {
	if in == nil {
		return []interface{}{nil}, nil
	}

	replicas := in.Replicas
	template, err := kubernetes.FlattenPodTemplateSpec(in.Template)
	if err != nil {
		return nil, err
	}
	restartPolicy := in.RestartPolicy

	return []interface{}{map[string]interface{}{
		"replicas":       replicas,
		"template":       template,
		"restart_policy": restartPolicy,
	}}, nil
}
