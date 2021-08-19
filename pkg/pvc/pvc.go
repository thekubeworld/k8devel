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
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/util"
)

// Instance type refers to the PVC object
type Instance struct {
	Name             string
	NamePrefix       string // "pvc-" if unspecified
	Namespace        string
	ClaimSize        string
	AccessModes      string
	Annotations      map[string]string
	Selector         *metav1.LabelSelector
	StorageClassName string
	VolumeMode       string
}

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

// Create will create a PVC
//
// Args:
//      - Client struct from client module
//      - Instance struct from pod module
//
// Return:
//      - error or nil
func Create(c *client.Client, p *Instance) (*v1.PersistentVolumeClaim, error) {

	volumeMode, _ := util.DetectVolumeMode(p.VolumeMode)

	if len(p.ClaimSize) == 0 {
		return nil, fmt.Errorf("size is required for a PVC")
	}

	accessModes, err := util.DetectVolumeAccessModes(p.AccessModes)
	if len(p.AccessModes) == 0 {
		return nil, fmt.Errorf("accessMode is required for a PVC")
	}

	if len(p.NamePrefix) == 0 {
		p.NamePrefix = "pvc-"
	}

	pvcSpec := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:         p.Name,
			GenerateName: p.NamePrefix,
			Namespace:    p.Namespace,
			Annotations:  p.Annotations,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			Selector:    p.Selector,
			AccessModes: accessModes,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(p.ClaimSize),
				},
			},
			StorageClassName: &p.StorageClassName,
			VolumeMode:       &volumeMode,
		},
	}

	pvc, err := c.Clientset.CoreV1().PersistentVolumeClaims(p.Namespace).Create(context.TODO(), pvcSpec, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("PVC Create API error: %v", err)
	}
	return pvc, nil
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
