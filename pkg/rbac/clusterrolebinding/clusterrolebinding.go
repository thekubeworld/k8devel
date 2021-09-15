package clusterrolebinding

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
	Name            string
	LabelKey        string
	LabelValue      string
	SubjectKind     string
	SubjectName     string
	SubjectAPIGroup string
	RoleRefName     string
	RoleRefKind     string
	RoleRefAPIGroup string
	Annotations     map[string]string
}

// Delete will delete a ClusteRroleBinding
//
// Args:
//	- Pointer to a Client struct
//	- clusterrolebinding name
//
// Returns:
//     error or nil
//
func Delete(c *client.Client, clusterrolebindingname string) error {
	err := c.Clientset.RbacV1().ClusterRoleBindings().Delete(
		context.TODO(),
		clusterrolebindingname,
		metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// List will list all ClusterRoleBinding
//
// Args:
//
//      - Client struct from client module
//
// Returns:
//      - pointer v1.NamespaceList or error
func List(c *client.Client) (*rbacv1.ClusterRoleBindingList, error) {
	clusterrolebindinglist, err := c.Clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return clusterrolebindinglist, nil
}

// Create will create a ClusterRoleBinding
//
// Args:
//     - Pointer to a Client struct
//     - Pointer to Instance
//
// Returns:
//     error or nil
//
func Create(c *client.Client, crb *Instance) error {

	role := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        crb.Name,
			Annotations: crb.Annotations,
			Labels: map[string]string{
				crb.LabelKey: crb.LabelValue,
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     crb.SubjectKind,
				APIGroup: crb.SubjectAPIGroup,
				Name:     crb.SubjectName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Name:     crb.RoleRefName,
			APIGroup: crb.RoleRefAPIGroup,
			Kind:     crb.RoleRefKind,
		},
	}

	_, err := c.Clientset.RbacV1().ClusterRoleBindings().Create(
		context.TODO(),
		role,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
