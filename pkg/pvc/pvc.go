package pvc

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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thekubeworld/k8devel/pkg/client"
)

// Delete will delete a pvc
//
// Args:
//      - Pointer to a Client struct
//      - namespace name
//	- pvc name
//
// Returns:
//     error or nil
//
func Delete(c *client.Client, namespace string, pvcname string) error {
	err := c.Clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcname, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// List will list all Namespaces
//
// Args:
//
//      - Client struct from client module
//
// Returns:
//      - pointer v1.NamespaceList or error
func List(c *client.Client, namespace string) (*v1.PersistentVolumeClaimList, error) {
	pvc, err := c.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pvc, nil
}
