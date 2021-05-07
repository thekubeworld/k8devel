package node

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
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thekubeworld/k8devel/pkg/client"
)

// Instance type refers to the ConfigMap object
type Instance struct {
	Name      string
	Namespace string
}

// GetIPFromNodes will list all ConfigMaps
//
// Args:
//
//      - Client struct from client module
//
// Returns:
//      - pointer v1.ConfigMapList or error
func GetIPFromNodes(c *client.Client) ([]string, error) {
	nodes, err := c.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodeip := []v1.NodeAddress{}
	var nodeList []string
	for i := 0; i < len(nodes.Items); i++ {
		nodeip = nodes.Items[i].Status.Addresses
		nodeList = append(nodeList, fmt.Sprint(nodeip[0].Address))
	}
	return nodeList, nil
}
