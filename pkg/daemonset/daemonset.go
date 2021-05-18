package daemonset

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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/util"
)

// Instance type refers to the Deployment object
type Instance struct {
	Name       string
	Namespace  string
	LabelKey   string
	LabelValue string
	Pod        struct {
		Name                  string
		Image                 string
		ContainerPortName     string
		ContainerPortProtocol string // "TCP" or "UDP"
		ContainerPort         int32
	}
}

// Create will create a daemonset
//
// Args:
//	- Client struct from client module
//	- daemonset from this module
//
// Returns:
//	- error
func Create(c *client.Client, d *Instance) error {

	podProtocol, err := util.DetectContainerPortProtocol(d.Pod.ContainerPortProtocol)
	if err != nil {
		return err
	}

	label := map[string]string{d.LabelKey: d.LabelValue}

	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.Name,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: label,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  d.Pod.Name,
							Image: d.Pod.Image,
							Ports: []v1.ContainerPort{
								{
									Name:          d.Pod.ContainerPortName,
									Protocol:      podProtocol,
									ContainerPort: d.Pod.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	// Create Daemonset
	_, err = c.Clientset.AppsV1().DaemonSets(d.Namespace).Create(
		context.TODO(),
		daemonset,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
