package job

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

// Instance type refers to the Job object
type Instance struct {
	Name          string
	Namespace     string
	Command       string
	RestartPolicy string // never, always, onfailure
	BackoffLimit  int32  // default is 6

	Pod struct {
		Name    string
		Image   string
		Command []string
	}
}

// Create will create a job
//
// Args:
//	- Client struct from client module
//	- Instance from this module
//
// Returns:
//	- error
func Create(c *client.Client, i *Instance) error {

	restartPolicy, err := util.DetectContainerRestartPolicy(i.RestartPolicy)
	if err != nil {
		return err
	}

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.Name,
			Namespace: i.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    i.Pod.Name,
							Image:   i.Pod.Image,
							Command: i.Pod.Command,
						},
					},
					RestartPolicy: restartPolicy,
				},
			},
			BackoffLimit: &i.BackoffLimit,
		},
	}

	_, err = c.Clientset.BatchV1().Jobs(i.Namespace).Create(
		context.TODO(),
		jobSpec, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
