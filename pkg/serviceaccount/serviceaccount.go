package serviceaccount

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

// Instance holds values for serviceaccount
type Instance struct {
	Name string
	Namespace string
	AutomountServiceAccountToken bool
}

// Create will create a new serviceaccount
//
// Args:
//     - Pointer to a Client struct
//     - Instance structure
//
// Returns:
//     error or nil
//     
func Create(c *client.Client, i *Instance) error {

	// by defaut we set AutomountServiceAccountToken as true
	autoservice := true

	if i.AutomountServiceAccountToken == false {
		autoservice = i.AutomountServiceAccountToken
	}
	SA := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: i.Namespace,
			Name: i.Name,
                },
                AutomountServiceAccountToken: &autoservice,
        }
        _, err := c.Clientset.CoreV1().ServiceAccounts(i.Namespace).Create(
			context.TODO(),
			SA,
			metav1.CreateOptions{})

	return err
}
