package paddle_job

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func elasticPolicyFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"min_replicas": {
			Type:        schema.TypeInt,
			Description: "minReplicas is the lower limit for the number of replicas to which the training job can scale down.  It defaults to null.",
			Optional:    true,
		},
		"max_replicas": {
			Type:        schema.TypeInt,
			Description: "upper limit for the number of pods that can be set by the autoscaler; cannot be smaller than MinReplicas, defaults to null.",
			Optional:    true,
		},
		"rdzv_backend": {
			Type:        schema.TypeString,
			Description: "RDZVBackend is the rendezvous backend to use.",
			Optional:    true,
		},
		"rdzv_port": {
			Type:        schema.TypeInt,
			Description: "RDZVPort is the port to use for rendezvous.",
			Optional:    true,
		},
		"rdzv_host": {
			Type:        schema.TypeString,
			Description: "RDZVHost is the host to use for rendezvous.",
			Optional:    true,
		},
		"rdzv_id": {
			Type:        schema.TypeString,
			Description: "RDZVID is the ID to use for rendezvous.",
			Optional:    true,
		},

		"rdzv_conf": {
			Type:        schema.TypeList,
			Description: "RDZVConf contains additional rendezvous configuration (<key1>=<value1>,<key2>=<value2>,...).",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeMap,
			},
		},
		"standalone": {
			Type:        schema.TypeBool,
			Description: "Start a local standalone rendezvous backend that is represented by a C10d TCP store on port 29400. Useful when launching single-node, multi-worker job. If specified --rdzv_backend, --rdzv_endpoint, --rdzv_id are auto-assigned; any explicitly set values are ignored.",
			Optional:    true,
		},
		"nproc_per_node": {
			Type:        schema.TypeInt,
			Description: "Number of workers per node; supported values: [auto, cpu, gpu, int].",
			Optional:    true,
		},
		"max_restarts": {
			Type:        schema.TypeInt,
			Description: "MaxRestarts is the maximum number of times a single pod can be restarted.",
			Optional:    true,
		},
	}

}

func expandElasticPolicy(elasticPolicy interface{}) (*kubeflowv1.ElasticPolicy, error) {
	if elasticPolicy == nil {
		return nil, nil
	}
	elasticPolicyMap := elasticPolicy.(map[string]interface{})
	elasticPolicyObj := &kubeflowv1.ElasticPolicy{}
	if v, ok := elasticPolicyMap["min_replicas"].(int32); ok {
		elasticPolicyObj.MinReplicas = &v
	}

	if v, ok := elasticPolicyMap["max_replicas"].(int32); ok {
		elasticPolicyObj.MaxReplicas = &v
	}

	if v, ok := elasticPolicyMap["rdzv_backend"].(string); ok {

		switch v {
		case "c10d":
			*elasticPolicyObj.RDZVBackend = kubeflowv1.BackendC10D
		case "etcd":
			*elasticPolicyObj.RDZVBackend = kubeflowv1.BackendETCD
		case "etcd-v2":
			*elasticPolicyObj.RDZVBackend = kubeflowv1.BackendETCDV2
		default:
			return nil, fmt.Errorf("invalid rdzv_backend %s", v)
		}

	}

	if v, ok := elasticPolicyMap["rdzv_port"].(int32); ok {
		elasticPolicyObj.RDZVPort = &v
	}

	if v, ok := elasticPolicyMap["rdzv_host"].(string); ok {
		elasticPolicyObj.RDZVHost = &v
	}

	if v, ok := elasticPolicyMap["rdzv_id"].(string); ok {
		elasticPolicyObj.RDZVID = &v
	}

	if v, ok := elasticPolicyMap["rdzv_conf"].([]interface{}); ok {
		elasticPolicyObj.RDZVConf = make([]kubeflowv1.RDZVConf, len(v))
		for i, vv := range v {
			// vv should be map[string]interface{}
			vvv := vv.(map[string]interface{})
			for k, vvvv := range vvv {
				elasticPolicyObj.RDZVConf[i].Key = k
				elasticPolicyObj.RDZVConf[i].Value = vvvv.(string)
			}
		}
	}

	if v, ok := elasticPolicyMap["standalone"].(bool); ok {
		elasticPolicyObj.Standalone = &v
	}

	if v, ok := elasticPolicyMap["nproc_per_node"].(int32); ok {
		elasticPolicyObj.NProcPerNode = &v
	}

	if v, ok := elasticPolicyMap["max_restarts"].(int32); ok {
		elasticPolicyObj.MaxRestarts = &v
	}

	return elasticPolicyObj, nil

}

func flattenElasticPolicy(elasticPolicy *kubeflowv1.ElasticPolicy) interface{} {
	if elasticPolicy == nil {
		return nil
	}
	elasticPolicyMap := make(map[string]interface{})
	if elasticPolicy.MinReplicas != nil {
		elasticPolicyMap["min_replicas"] = *elasticPolicy.MinReplicas
	}
	if elasticPolicy.MaxReplicas != nil {
		elasticPolicyMap["max_replicas"] = *elasticPolicy.MaxReplicas
	}
	if elasticPolicy.RDZVBackend != nil {
		elasticPolicyMap["rdzv_backend"] = *elasticPolicy.RDZVBackend
	}
	if elasticPolicy.RDZVPort != nil {
		elasticPolicyMap["rdzv_port"] = *elasticPolicy.RDZVPort
	}
	if elasticPolicy.RDZVHost != nil {
		elasticPolicyMap["rdzv_host"] = *elasticPolicy.RDZVHost
	}
	if elasticPolicy.RDZVID != nil {
		elasticPolicyMap["rdzv_id"] = *elasticPolicy.RDZVID
	}
	if elasticPolicy.RDZVConf != nil {
		rdzvConf := make([]interface{}, len(elasticPolicy.RDZVConf))
		for i, v := range elasticPolicy.RDZVConf {
			rdzvConf[i] = map[string]interface{}{
				v.Key: v.Value,
			}
		}
		elasticPolicyMap["rdzv_conf"] = rdzvConf
	}
	if elasticPolicy.Standalone != nil {
		elasticPolicyMap["standalone"] = *elasticPolicy.Standalone
	}
	if elasticPolicy.NProcPerNode != nil {
		elasticPolicyMap["nproc_per_node"] = *elasticPolicy.NProcPerNode
	}
	if elasticPolicy.MaxRestarts != nil {
		elasticPolicyMap["max_restarts"] = *elasticPolicy.MaxRestarts
	}
	return elasticPolicyMap
}
