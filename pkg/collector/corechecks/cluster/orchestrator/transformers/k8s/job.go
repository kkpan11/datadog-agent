// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build orchestrator

package k8s

import (
	model "github.com/DataDog/agent-payload/v5/process"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/cluster/orchestrator/processors"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/cluster/orchestrator/transformers"
	batchv1 "k8s.io/api/batch/v1"
)

// ExtractJob returns the protobuf model corresponding to a Kubernetes Job
// resource.
func ExtractJob(ctx processors.ProcessorContext, j *batchv1.Job) *model.Job {
	job := model.Job{
		Metadata: extractMetadata(&j.ObjectMeta),
		Spec:     &model.JobSpec{},
		Status: &model.JobStatus{
			Active:           j.Status.Active,
			ConditionMessage: extractJobConditionMessage(j.Status.Conditions),
			Failed:           j.Status.Failed,
			Succeeded:        j.Status.Succeeded,
		},
	}

	if j.Spec.ActiveDeadlineSeconds != nil {
		job.Spec.ActiveDeadlineSeconds = *j.Spec.ActiveDeadlineSeconds
	}
	if j.Spec.BackoffLimit != nil {
		job.Spec.BackoffLimit = *j.Spec.BackoffLimit
	}
	if j.Spec.Completions != nil {
		job.Spec.Completions = *j.Spec.Completions
	}
	if j.Spec.ManualSelector != nil {
		job.Spec.ManualSelector = *j.Spec.ManualSelector
	}
	if j.Spec.Parallelism != nil {
		job.Spec.Parallelism = *j.Spec.Parallelism
	}
	if j.Spec.Selector != nil {
		job.Spec.Selectors = extractLabelSelector(j.Spec.Selector)
	}

	if j.Status.StartTime != nil {
		job.Status.StartTime = j.Status.StartTime.Unix()
	}
	if j.Status.CompletionTime != nil {
		job.Status.CompletionTime = j.Status.CompletionTime.Unix()
	}

	if len(j.Status.Conditions) > 0 {
		jConditions, conditionTags := extractJobConditions(j)
		job.Conditions = jConditions
		job.Tags = append(job.Tags, conditionTags...)
	}

	job.Spec.ResourceRequirements = ExtractPodTemplateResourceRequirements(j.Spec.Template)

	pctx := ctx.(*processors.K8sProcessorContext)
	job.Tags = append(job.Tags, transformers.RetrieveUnifiedServiceTags(j.ObjectMeta.Labels)...)
	job.Tags = append(job.Tags, transformers.RetrieveMetadataTags(j.ObjectMeta.Labels, j.ObjectMeta.Annotations, pctx.LabelsAsTags, pctx.AnnotationsAsTags)...)

	return &job
}

func extractJobConditionMessage(conditions []batchv1.JobCondition) string {
	for _, c := range conditions {
		if c.Type == batchv1.JobFailed && c.Message != "" {
			return c.Message
		}
	}
	return ""
}

// extractJobConditions iterates over job conditions and returns:
// - the payload representation of those conditions
// - the list of tags that will enable pod filtering by condition
func extractJobConditions(p *batchv1.Job) ([]*model.JobCondition, []string) {
	conditions := make([]*model.JobCondition, 0, len(p.Status.Conditions))
	conditionTags := make([]string, 0, len(p.Status.Conditions))

	for _, condition := range p.Status.Conditions {
		c := &model.JobCondition{
			Message: condition.Message,
			Reason:  condition.Reason,
			Status:  string(condition.Status),
			Type:    string(condition.Type),
		}
		if !condition.LastTransitionTime.IsZero() {
			c.LastTransitionTime = condition.LastTransitionTime.Unix()
		}

		if !condition.LastProbeTime.IsZero() {
			c.LastProbeTime = condition.LastProbeTime.Unix()
		}

		conditions = append(conditions, c)

		conditionTag := createConditionTag(string(condition.Type), string(condition.Status))
		conditionTags = append(conditionTags, conditionTag)
	}

	return conditions, conditionTags
}
