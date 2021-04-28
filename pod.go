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
	"bytes"
	"context"
	"strings"
	"errors"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	v1 "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"
)

// Pod type refers to the Pod object
type Pod struct {
	Name string
	Namespace string
	Image string
	Command []string
	CommandArgs []string
	LabelKey string
	LabelValue string
}

// ExecCmdPod executes a command inside a POD
//
// Args:
//      Client - struct from client module
//	podName	- The pod name
//	cmd - Array (string)
//
// Returns:
//	stdout, stderr as bytes.Buffer or error
func ExecCmdPod(c *Client,
		podName string,
		nameSpace string,
		cmd []string) (bytes.Buffer, bytes.Buffer, error) {

	logrus.Infof("\n")
	logrus.Infof("Executing command: %s", cmd)
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

// GetIPFromPod will return the pod IP address
//
// Args:
//
//	- Client struct from client module
//	- pod name
//	- namespace
//
// Returns:
//	- the IP as string or error
func GetIPFromPod(c *Client,
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
func FindPodsWithNameContains(c *Client,
                substring string,
                namespace string) ([]string, int){

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

// ExistsPod will check if the pod exists or not
//
// Args:
//     - Pointer to a Client struct
//
// Returns:
//     string (namespace name) OR error type
//
func ExistsPod(c *Client, podName string, namespace string) (string, error) {
        exists, err := c.Clientset.CoreV1().Services(namespace).
                Get(context.TODO(), podName, metav1.GetOptions{})
        if err != nil {
                return "", err
        }

        return exists.Name, nil
}

// CreatePod will create a POD
//
// Args:
//      - Client struct from client module
//      - Instance struct from pod module
//
// Return:
//      - error or nil
func CreatePod(c *Client, p *Pod) error {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta {
			Name: p.Name,
			Namespace: p.Namespace,
			Labels: map[string]string {
				p.LabelKey: p.LabelValue,
			},
		},
		Spec: v1.PodSpec {
			Containers: []v1.Container{
				{
					Name: p.Name,
					Image: p.Image,
					Command: p.Command,
					Args: p.CommandArgs,
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

	logrus.Infof("Creating pod: %s namespace: %s", p.Name, p.Namespace)

	// Double check if pod is created
        for i := 0; i < c.NumberMaxOfAttemptsPerTask; i++ {
                exists, _ := ExistsPod(c, p.Name, p.Namespace)
		if exists != "" {
			break
		}
                time.Sleep(time.Duration(c.TimeoutTaksInSec) * time.Second)
		if i == c.NumberMaxOfAttemptsPerTask {
			return errors.New("cannot create pod")
		}
        }
	return nil
}
