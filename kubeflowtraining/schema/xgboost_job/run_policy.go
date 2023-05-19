package xgboost_job

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
)

func runPolicyFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"clean_pod_policy": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Running",
			Description: "CleanPodPolicy defines the policy to kill pods after the job completes.",
		},
		"ttl_seconds_after_finished": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "TTLSecondsAfterFinished is the TTL to clean up jobs.",
		},
		"active_deadline_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Specifies the duration in seconds relative to the startTime that the job may be active before the system tries to terminate it; value must be positive integer.",
		},
		"backoff_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Optional number of retries before marking this job failed.",
		},
		"scheduling_policy": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "SchedulingPolicy encapsulates various scheduling policies of the distributed training job, for example `minAvailable` for gang-scheduling.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"min_available": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "MinAvailable is the minimum number of workers available for scheduling.",
					},
					"queue": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Queue is the name of the queue to schedule the job to.",
					},
					"min_resources": {
						Type:        schema.TypeMap,
						Optional:    true,
						Description: "MinResources is the minimum resources required for scheduling.",
					},
					"priority_class": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "PriorityClass is the name of the priority class to schedule the job to.",
					},
					"schedule_timeout_seconds": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "ScheduleTimeoutSeconds is the timeout for scheduling the job.",
					},
				},
			},
		},
	}
}

func expandRunPolicy(runPolicy interface{}) (*commonv1.RunPolicy, error) {
	if runPolicy == nil {
		return nil, nil
	}
	rp := &commonv1.RunPolicy{}
	m := runPolicy.(map[string]interface{})
	if v, ok := m["clean_pod_policy"]; ok {
		*rp.CleanPodPolicy = commonv1.CleanPodPolicy(v.(string))
	}
	if v, ok := m["ttl_seconds_after_finished"]; ok {
		*rp.TTLSecondsAfterFinished = int32(v.(int))
	}
	if v, ok := m["active_deadline_seconds"]; ok {
		*rp.ActiveDeadlineSeconds = int64(v.(int))
	}
	if v, ok := m["backoff_limit"]; ok {
		*rp.BackoffLimit = int32(v.(int))
	}
	if v, ok := m["scheduling_policy"]; ok {
		rp.SchedulingPolicy = expandSchedulingPolicy(v)
	}
	return rp, nil
}

func expandSchedulingPolicy(v interface{}) *commonv1.SchedulingPolicy {
	if v == nil {
		return nil
	}
	m := v.([]interface{})[0].(map[string]interface{})
	sp := &commonv1.SchedulingPolicy{}
	if v, ok := m["min_available"]; ok {
		*sp.MinAvailable = int32(v.(int))
	}
	if v, ok := m["queue"]; ok {
		sp.Queue = v.(string)
	}
	// if v, ok := m["min_resources"]; ok {
	// 	*sp.MinResources = expandResources(v)
	// }
	if v, ok := m["priority_class"]; ok {
		sp.PriorityClass = v.(string)
	}
	if v, ok := m["schedule_timeout_seconds"]; ok {
		*sp.ScheduleTimeoutSeconds = int32(v.(int))
	}
	return sp
}

func flattenSchedulingPolicy(sp *commonv1.SchedulingPolicy) []interface{} {
	if sp == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{}
	m["min_available"] = int(*sp.MinAvailable)
	m["queue"] = sp.Queue
	// m["min_resources"] = flattenResources(sp.MinResources)
	m["priority_class"] = sp.PriorityClass
	m["schedule_timeout_seconds"] = int(*sp.ScheduleTimeoutSeconds)
	return []interface{}{m}
}

func flattenRunPolicy(rp *commonv1.RunPolicy) map[string]interface{} {
	m := map[string]interface{}{}
	if rp == nil {
		return m
	}
	m["clean_pod_policy"] = string(*rp.CleanPodPolicy)
	m["ttl_seconds_after_finished"] = int(*rp.TTLSecondsAfterFinished)
	m["active_deadline_seconds"] = int(*rp.ActiveDeadlineSeconds)
	m["backoff_limit"] = int(*rp.BackoffLimit)
	m["scheduling_policy"] = flattenSchedulingPolicy(rp.SchedulingPolicy)
	return m
}
