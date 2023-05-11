package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func PodTemplateFields(owner string) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"metadata": metadataSchema(owner, true),
		"spec": {
			Type:        schema.TypeList,
			Description: fmt.Sprintf("Spec of the pods owned by the %s", owner),
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSpecFields(true, false, false),
			},
		},
	}
	return s
}

func ExpandPodTemplate(l []interface{}) (*corev1.PodTemplateSpec, error) {
	obj := &corev1.PodTemplateSpec{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})

	obj.ObjectMeta = expandMetadata(in["metadata"].([]interface{}))

	if v, ok := in["spec"].([]interface{}); ok && len(v) > 0 {
		podSpec, err := expandPodSpec(in["spec"].([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.Spec = *podSpec
	}
	return obj, nil
}

func FlattenPodTemplateSpec(t corev1.PodTemplateSpec, prefix ...string) ([]interface{}, error) {
	template := make(map[string]interface{})

	// metaPrefix := "spec.0.template.0."
	// if len(prefix) > 0 {
	// 	metaPrefix = prefix[0]
	// }
	// template["metadata"] = flattenMetadata(t.ObjectMeta, d, metaPrefix)
	spec, err := flattenPodSpec(t.Spec)
	if err != nil {
		return []interface{}{template}, err
	}
	template["spec"] = spec

	return []interface{}{template}, nil
}
