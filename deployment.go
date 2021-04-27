package k8devel

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
	"time"
	v1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sirupsen/logrus"
)

// Deployment type refers to the Deployment object
type Deployment struct {
	Name string
	Namespace string
	Replicas int32
	LabelKey string
	LabelValue string
	Pod struct {
		Name string
		Image string
		ContainerPortName string
		ContainerPortProtocol string  // "TCP" or "UDP"
		ContainerPort int32
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
func CreateDeployment(c *Client, d *Deployment) error {

        deployClient := c.Clientset.AppsV1().Deployments(d.Namespace)
	podProtocol, err := DetectContainerPortProtocol(d.Pod.ContainerPortProtocol)
	if err != nil {
		logrus.Fatal(err)
	}

        deployment := &apiv1.Deployment {
                ObjectMeta: metav1.ObjectMeta{
                        Name: d.Name,
                },
                Spec: apiv1.DeploymentSpec{
                        Replicas: &d.Replicas,
                        Selector: &metav1.LabelSelector{
                                MatchLabels: map[string]string{
                                        d.LabelKey: d.LabelValue,
                                },
                        },
                        Template: v1.PodTemplateSpec{
                                ObjectMeta: metav1.ObjectMeta{
                                        Labels: map[string]string{
                                                d.LabelKey: d.LabelValue,
                                        },
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

	// Create Deployment
	logrus.Infof("Creating deployment: %s", d.Name)
	result, err := deployClient.Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}
	logrus.Infof("Created deployment: %s namespace: %s",
		result.GetObjectMeta().GetName(),
		d.Namespace)

	return nil
}

// DeleteDeployment will delete an deployment
//
// Args:
//      - Client struct from client module
//      - deployment name
//      - namespace
// Return:
//      - error or nil
func DeleteDeployment(c *Client, deployment string, namespace string) error {
        _, err := c.Clientset.AppsV1().Deployments(namespace).
                Get(context.TODO(), deployment, metav1.GetOptions{})
        if err != nil {
                return err
        }

	logrus.Info("\n")
        logrus.Infof("Deleting deployment: %s namespace: %s...",
                deployment,
                namespace)

        // Double check service is removed
        for i := 0; i < c.NumberMaxOfAttemptsPerTask; i++ {
                _, err := ExistsDeployment(c, deployment, namespace)
                if err != nil {
                        logrus.Infof("Deleted deployment: %s namespace: %s",
                                deployment,
                                namespace)
                        break
                }
		c.Clientset.AppsV1().Deployments(namespace).Delete(
			context.TODO(),
			deployment,
			metav1.DeleteOptions{})

                time.Sleep(time.Duration(c.TimeoutTaksInSec) * time.Second)
        }

        return nil
}

// Exists will check if the service exists or not
//
// Args:
// 	Pointer to a Client struct
//	Service Name
//	Namespace
//
// Returns:
//     string (namespace name) OR error type
//
func ExistsDeployment(c *Client, deployment string, namespace string) (string, error) {
	exists, err := c.Clientset.AppsV1().Deployments(namespace).
		Get(context.TODO(), deployment, metav1.GetOptions{})
        if err != nil {
                return "", err
        }

        return exists.Name, nil
}
