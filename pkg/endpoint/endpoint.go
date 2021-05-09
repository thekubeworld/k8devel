package endpoint

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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/util"
)

// Instance type refers to the Endpoint object
type Instance struct {
	Name         string
	IP           string
	LabelKey     string
	LabelValue   string
	Namespace    string
	EndpointPort struct {
		Name     string
		Port     int32
		Protocol string
	}
}

// List will list ALL endpoints from a namespace
//
// Args:
//	- Client struct from client module
//	- Instance from endpoint module
//
// Return:
//	- error or nil
func List(c *client.Client, e *Instance) {
	epoints, _ := c.Clientset.CoreV1().Endpoints(e.Namespace).List(context.TODO(), metav1.ListOptions{})
	fmt.Printf("\n")
	fmt.Printf("Listing endpoints in namespace %s:\n", e.Namespace)
	for _, ep := range epoints.Items {
		fmt.Printf("\tName: " + ep.Name)
	}
}

// Patch will patch an endpoint object
//
// Args:
//	- Client struct from client module
//	- Instance from endpoint module
//
// Return:
//	- error or nil
func Patch(c *client.Client, e *Instance) error {
	fmt.Printf("\n")
	fmt.Printf("Patching endpoint: %s namespace: %s\n",
		e.Name,
		e.Namespace)
	_, err := c.Clientset.CoreV1().Endpoints(e.Namespace).Get(
		context.TODO(),
		e.Name,
		metav1.GetOptions{})
	if err != nil {
		return err
	}

	endpointPatch, err := json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": map[string]string{
				e.LabelKey: e.LabelValue,
			},
		},
		"subsets": []map[string]interface{}{
			{
				"addresses": []map[string]string{
					{
						"ip": e.IP,
					},
				},
				"ports": []map[string]interface{}{
					{
						"name":     e.EndpointPort.Name,
						"port":     e.EndpointPort.Port,
						"protocol": e.EndpointPort.Protocol,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// Executing the patch
	_, err = c.Clientset.CoreV1().Endpoints(e.Namespace).Patch(
		context.TODO(),
		e.Name,
		types.StrategicMergePatchType,
		[]byte(endpointPatch),
		metav1.PatchOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Patched endpoint: %s namespace %s\n", e.Name, e.Namespace)
	return nil
}

// Create will create an endpoint
//
// Args:
//	- Client struct from client module
//	- Instance from endpoint module
//
// Return:
//	- error or nil
func Create(c *client.Client, e *Instance) error {
	fmt.Printf("\n")
	fmt.Printf("Creating endpoint: %s namespace: %s\n",
		e.Name,
		e.Namespace)

	proto, err := util.DetectContainerPortProtocol(e.EndpointPort.Protocol)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	epoints := &v1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name: e.Name,
		},
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: e.IP,
					},
				},
				Ports: []v1.EndpointPort{
					{
						Name:     e.EndpointPort.Name,
						Port:     e.EndpointPort.Port,
						Protocol: proto,
					},
				},
			},
		},
	}
	_, err = c.Clientset.CoreV1().Endpoints(e.Namespace).Create(
		context.TODO(),
		epoints,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created endpoint: %s\n", e.Name)
	return nil
}

// Show will display a specific endpoint
// Args:
// 	- Client struct from client module
//	- endpoint name
func Show(c *client.Client, endpoint string, namespace string) error {
	epoints, err := c.Clientset.CoreV1().Endpoints(namespace).Get(
		context.TODO(),
		endpoint,
		metav1.GetOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("\n")
	fmt.Printf("Showing information about endpoint: %s namespace: %s\n",
		endpoint,
		namespace)

	if len(epoints.Subsets) == 0 {
		fmt.Printf("Subsets: []")
	}

	port := ""
	if len(epoints.Subsets[0].Ports) > 0 {
		port = strconv.FormatInt(int64(epoints.Subsets[0].Ports[0].Port), 10)
		for _, p := range epoints.Subsets[0].Ports {
			if p.Name != "" {
				fmt.Printf("\tEndpointPort Name: %s\n", p.Name)
				fmt.Printf("\tEndpointPort Port: %s\n", strconv.FormatInt(int64(p.Port), 10))
			}
		}

		for _, address := range epoints.Subsets[0].Addresses {
			fmt.Printf("\tEndpointAddress: %s\n", address.IP)
			fmt.Printf("\tEndpointAddress Port: %s\n", port)
		}
	}

	return nil
}

// Delete will delete an endpoint
//
// Args:
// 	- Client struct from client module
//	- endpoint name
//	- namespace
// Return:
//	- error or nil
func Delete(c *client.Client, endpoint string, namespace string) error {
	inst := Instance{Name: endpoint, Namespace: namespace}

	_, err := c.Clientset.CoreV1().Endpoints(inst.Namespace).
		Get(context.TODO(), inst.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("\n")
	fmt.Printf("Deleting endpoint: %s namespace: %s...\n",
		inst.Name,
		inst.Namespace)

	// Double check endpoint is removed
	for i := 0; i < c.NumberMaxOfAttemptsPerTask; i++ {
		_, err := Exists(c, &inst)
		if err != nil {
			fmt.Printf("Deleted endpoint: %s namespace: %s\n",
				inst.Name,
				inst.Namespace)
			break
		}
		c.Clientset.CoreV1().Endpoints(inst.Namespace).Delete(
			context.TODO(),
			inst.Name,
			metav1.DeleteOptions{})

		time.Sleep(time.Duration(c.TimeoutTaskInSec) * time.Second)
	}

	return nil
}

// Exists will check if the endpoint exists or not
//
// Args:
//	- Client struct from client module
//	- Instance struct from this module
//
// Returns:
//     bool OR error type
//
func Exists(c *client.Client, e *Instance) (string, error) {
	exists, err := c.Clientset.CoreV1().Endpoints(e.Namespace).Get(
		context.TODO(),
		e.Name,
		metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return exists.Name, nil
}
