package namespace

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
        "github.com/sirupsen/logrus"

	"github.com/thekubeworld/k8devel/pkg/client"

)

// Delete will delete a namespace
//
// Args:
//	- Pointer to a Client struct
//	- namespace name
//
// Returns:
//     error or nil
//     
func Delete(c *client.Client, namespace string) error {
        err := c.Clientset.CoreV1().Namespaces().Delete(
                context.TODO(),
                namespace,
                metav1.DeleteOptions{})
        if err != nil {
                return err
        }
        return nil
}

// Create will create a namespace
//
// Args:
//     - Pointer to a Client struct
//	- namespace name
//
// Returns:
//     error or nil
//     
func Create(c *client.Client, namespace string) error {
        logrus.Infof("Creating namespace: %s", namespace)

        ns := &v1.Namespace {
                ObjectMeta: metav1.ObjectMeta {
                Name: namespace,
                        Labels: map[string]string {
                                "name": namespace,
                        },
                },
        }

        _, err := c.Clientset.CoreV1().Namespaces().Create(
                context.TODO(),
                ns,
                metav1.CreateOptions{})
        if err != nil {
                return err
        }
        logrus.Infof("Created namespace: %s", c.Namespace)

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
