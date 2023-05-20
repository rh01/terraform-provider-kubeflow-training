package paddle_job

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining/schema/kubernetes"
)

func paddleJobReplicaSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"master": paddleJobReplicaSpecSchema(),
		"worker": paddleJobReplicaSpecSchema(),
	}
}

func paddleJobReplicaSpecTemplateFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"replicas": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"template": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: kubernetes.PodTemplateFields("paddlejob"),
			},
			Optional: true,
		},
		"restart_policy": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "Never",
		},
	}
}

func paddleJobReplicaSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: paddleJobReplicaSpecTemplateFields(),
		},
		Optional: true,
	}
}

func expandPaddleJobReplicaSpec(l []interface{}) (*commonv1.ReplicaSpec, error) {
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

func flattenPaddleJobReplicaSpec(in *commonv1.ReplicaSpec) ([]interface{}, error) {
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
