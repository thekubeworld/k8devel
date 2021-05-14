package cronjob

/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/util"
)

// Instance type refers to the Deployment object
type Instance struct {
	Name                       string
	Namespace                  string
	Schedule                   string
	ConcurrencyPolicy          string // Allow, Forbid and Replace
	RestartPolicy              string
	Command                    []string
	Parallelism                int32 // 0 means unlimited parallelism.
	BackoffLimit               int32 // default is 6
	Completions                int32
	SuccessfulJobsHistoryLimit int32
	FailedJobsHistoryLimit     int32

	Pod struct {
		Name    string
		Image   string
		Command []string
	}
}

// Create will create a deployment
//
// Args:
//	- Client struct from client module
//	- Deployment from this module
//
// Returns:
//	- error
func Create(c *client.Client, i *Instance) error {

	restartPolicy, err := util.DetectContainerRestartPolicy(i.RestartPolicy)
	concurrencyPolicy, err := util.DetectConcurrencyPolicy(i.ConcurrencyPolicy)
	if err != nil {
		return err
	}

	job := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: i.Name,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          i.Schedule,
			ConcurrencyPolicy: concurrencyPolicy,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Parallelism:  &i.Parallelism,
					Completions:  &i.Completions,
					BackoffLimit: &i.BackoffLimit,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							RestartPolicy: restartPolicy,
							Containers: []v1.Container{
								{
									Name:    i.Pod.Name,
									Image:   i.Pod.Image,
									Command: i.Pod.Command,
								},
							},
						},
					},
				},
			},
		},
	}
	job.Spec.SuccessfulJobsHistoryLimit = &i.SuccessfulJobsHistoryLimit
	job.Spec.FailedJobsHistoryLimit = &i.FailedJobsHistoryLimit
	job.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Command = i.Command

	_, err = c.Clientset.BatchV1().CronJobs(i.Namespace).Create(
		context.TODO(),
		job,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
