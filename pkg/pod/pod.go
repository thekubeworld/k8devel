package pod

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
	"bytes"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/thekubeworld/k8devel/pkg/client"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Instance type refers to the Pod object
type Instance struct {
	Name        string
	Namespace   string
	Image       string
	Command     []string
	CommandArgs []string
	LabelKey    string
	LabelValue  string
}

// ExecCmd executes a command inside a POD
//
// Args:
//      Client - struct from client module
//	podName	- The pod name
//	cmd - Array (string)
//
// Returns:
//	stdout, stderr as bytes.Buffer or error
func ExecCmd(c *client.Client,
	podName string,
	nameSpace string,
	cmd []string) (bytes.Buffer, bytes.Buffer, error) {

	restClient := c.Clientset.CoreV1().RESTClient()

	req := restClient.Post().Resource("pods").Name(podName).
		Namespace(nameSpace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	var stdout, stderr bytes.Buffer
	exec, err := remotecommand.NewSPDYExecutor(
		c.Restconfig,
		"POST",
		req.URL())
	if err != nil {
		return stdout, stderr, err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return stdout, stderr, err
	}
	return stdout, stderr, nil
}

// GetLastTimeConditionHappened
// Get the last time a condition happened in a pod
//
// Conditions: 
//
//     PodScheduled: the Pod has been scheduled to a node.
//     ContainersReady: all containers in the Pod are ready.
//     Initialized: all init containers have started successfully.
//     Ready: the Pod is able to serve requests and should be added
//            to the load balancing pools of all matching
//
// Args:
//
//	- Client struct from client module
//	- pod name
//	- namespace
//
// Returns:
//	- the IP as string or error
func GetLastTimeConditionHappened(c *client.Client,
	condition string,
	podName string,
	nameSpace string) (metav1.Time, error) {

	pod, err := c.Clientset.CoreV1().Pods(nameSpace).Get(
		context.TODO(),
		podName,
		metav1.GetOptions{})
	if err != nil {
		return metav1.Time{}, err
	}

	var requiredCondition v1.PodConditionType

	switch condition {
	case "ContainersReady":
		requiredCondition = v1.ContainersReady
		break
	case "Initialized":
		requiredCondition = v1.PodInitialized
		break
	case "Ready":
		requiredCondition = v1.PodReady
		break
	case "PodScheduled":
		requiredCondition = v1.PodScheduled
		break
	default:
		return metav1.Time{}, errors.New("condition not recognized, use: " +
				"ContainersReady, Initialized, " +
				"Ready or PodScheduled")
	}

	for _, cond := range pod.Status.Conditions {
               if cond.Type == requiredCondition && cond.Status == v1.ConditionTrue {
		        // Type is the type of the condition.
			// More info: https://kubernetes.io/docs/
			//concepts/workloads/pods/pod-lifecycle#pod-conditions
			//fmt.Println(cond.Type)

			// Status is the status of the condition.
			// Can be True, False, Unknown.
			// More info: https://kubernetes.io/docs/concepts/
			//workloads/pods/pod-lifecycle#pod-conditions
			//fmt.Println(cond.Status)

			// Last time we probed the condition.
                        // +optional
			//fmt.Println(cond.LastProbeTime)

			// Last time the condition transitioned
			//from one status to another. +optional
			return cond.LastTransitionTime, nil
		}
	}
	return metav1.Time{}, errors.New("unable to get the last time the " +
				"condition happened")
}
// GetIP will return the pod IP address
//
// Args:
//
//	- Client struct from client module
//	- pod name
//	- namespace
//
// Returns:
//	- the IP as string or error
func GetIP(c *client.Client,
	podName string,
	nameSpace string) (string, error) {

	pod, err := c.Clientset.CoreV1().Pods(nameSpace).Get(
		context.TODO(),
		podName,
		metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return pod.Status.PodIP, nil

}

// FindPodsWithNameContains will find pods with
// substring provided
//
// Args:
//      - Client struct from client module
//      - substring to be found
//      - namespace
//
// Return:
//      - error or nil
func FindPodsWithNameContains(c *client.Client,
	substring string,
	namespace string) ([]string, int) {

	var podsFound []string
	listPods, _ := c.Clientset.CoreV1().Pods(namespace).List(
		context.TODO(),
		metav1.ListOptions{})

	for _, p := range listPods.Items {
		if strings.Contains(p.Name, substring) {
			podsFound = append(podsFound, p.Name)
		}
	}

	return podsFound, len(podsFound)
}

// isPodRunning will check if the pod is running
//
// Args:
//	- Pointer to a client struct
//	- podname
//	- namespace
//
// Returns:
//	bool or error
func isPodRunning(c *client.Client, podname, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := c.Clientset.CoreV1().Pods(namespace).Get(
			context.TODO(),
			podname,
			metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case v1.PodRunning:
			return true, nil
		case v1.PodFailed, v1.PodSucceeded:
			return false, errors.New("pod not running")
		}
		return false, nil
	}
}

// waitForPodRunning will execute wait.PollImmediate
//
// Args:
//	- Pointer to a client struct
//	- podname
//	- namespace
//
// Returns:
//	nil or error
func waitForPodRunning(c *client.Client, namespace, podname string, timeout time.Duration) error {
	return wait.PollImmediate(
		time.Second,
		timeout,
		isPodRunning(c, podname, namespace))
}

// WaitForPodInRunningState will execute waitForPodRunning
//
// Args:
//	- Pointer to a client struct
//	- podname
//	- namespace
//
// Returns:
//	nil or error
func WaitForPodInRunningState(c *client.Client, podname string, namespace string) error {
	if err := waitForPodRunning(c,
		namespace,
		podname,
		time.Duration(c.TimeoutTaskInSec)*time.Second); err != nil {
		return err
	}
	return nil
}

// Exists will check if the pod exists or not
//
// Args:
//     - Pointer to a Client struct
//
// Returns:
//     string (namespace name) OR error type
//
func Exists(c *client.Client, podName string, namespace string) (string, error) {
	exists, err := c.Clientset.CoreV1().Services(namespace).
		Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return exists.Name, nil
}

// Create will create a POD
//
// Args:
//      - Client struct from client module
//      - Instance struct from pod module
//
// Return:
//      - error or nil
func Create(c *client.Client, p *Instance) error {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
			Labels: map[string]string{
				p.LabelKey: p.LabelValue,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    p.Name,
					Image:   p.Image,
					Command: p.Command,
					Args:    p.CommandArgs,
				},
			},
		},
	}

	_, err := c.Clientset.CoreV1().Pods(p.Namespace).Create(
		context.TODO(),
		pod,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	err = WaitForPodInRunningState(c, p.Name, p.Namespace)
	if err != nil {
		return err
	}
	return nil
}
