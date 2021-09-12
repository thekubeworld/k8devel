package clusterrole

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

	"github.com/thekubeworld/k8devel/pkg/client"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Instance type refers to the Pod object
type Instance struct {
	Name           string
	LabelKey       string
	LabelValue     string
	APIGroups      []string
	Resources      []string
	ResourcesNames []string
	Verbs          []string
}

// Delete will delete a clusterrole
//
// Args:
//	- Pointer to a Client struct
//	- clusterrolename name
//
// Returns:
//     error or nil
//
func Delete(c *client.Client, clusterrolename string) error {
	err := c.Clientset.RbacV1().ClusterRoles().Delete(
		context.TODO(),
		clusterrolename,
		metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// List will list all ClusterRoles
//
// Args:
//
//      - Client struct from client module
//
// Returns:
//      - pointer v1.NamespaceList or error
func List(c *client.Client) (*rbacv1.ClusterRoleList, error) {
	clusterrolelist, err := c.Clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return clusterrolelist, nil
}

// Create will create a clusterrole
//
// Args:
//     - Pointer to a Client struct
//     - Pointer to Instance
//
// Returns:
//     error or nil
//
func Create(c *client.Client, cr *Instance) error {

	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name,
			Labels: map[string]string{
				cr.LabelKey: cr.LabelValue,
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     cr.APIGroups,
				Resources:     cr.Resources,
				Verbs:         cr.Verbs,
				ResourceNames: cr.ResourcesNames,
			},
		},
	}

	_, err := c.Clientset.RbacV1().ClusterRoles().Create(
		context.TODO(),
		role,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Exists will check if the namespace exists or not
//
// Args:
//     - Pointer to a Client struct
//	- namespace name
//
// Returns:
//     string (namespace name) OR error type
//
func Exists(c *client.Client, namespace string) (string, error) {
	exists, err := c.Clientset.CoreV1().Namespaces().Get(
		context.TODO(),
		namespace,
		metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return exists.Name, nil
}
