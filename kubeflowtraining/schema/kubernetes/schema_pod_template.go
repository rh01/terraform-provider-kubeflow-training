package kubernetes

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func PodTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"spec": {
			Type:        schema.TypeList,
			Description: "Specification of the desired behavior of the pod. More info: " + "" + "https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status" + "" + "",
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

	spec, err := flattenPodSpec(t.Spec)
	if err != nil {
		return []interface{}{template}, err
	}
	template["spec"] = spec

	return []interface{}{template}, nil
}
