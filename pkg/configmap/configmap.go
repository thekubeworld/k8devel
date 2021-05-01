package configmap

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
	"time"
	v1 "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sirupsen/logrus"

	"github.com/thekubeworld/k8devel/pkg/client"
)

// Instance type refers to the ConfigMap object
type Instance struct {
        Name string
        Namespace string
        ConfigKey string
        ConfigValue string
}


// Show will list all ConfigMaps
//
// Args:
//
//      - Client struct from client module
//
// Returns:
//      - pointer v1.ConfigMapList or error
func Show(c *client.Client, configmap string, namespace string) (*v1.ConfigMap, error) {
	cfmap, err := c.Clientset.CoreV1().ConfigMaps(namespace).Get(
		context.TODO(),
		configmap,
		metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cfmap, nil
}

// ListAll will list all ConfigMaps
//
// Args:
//
//      - Client struct from client module
//
// Returns:
//      - pointer v1.ConfigMapList or error
func ListAll(c *client.Client) (*v1.ConfigMapList, error) {
	configmap, err := c.Clientset.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return configmap, nil
}

// Delete deletes a configmap
//
// Args:
//	Client - client struct from the client module
//	configmap - ConigMap Name
//	namespace - Namespace
//
//   Returns:
//      error or nil
func Delete(c *client.Client, configmap string, namespace string) error {
	_, err := c.Clientset.CoreV1().ConfigMaps(namespace).
		Get(context.TODO(), configmap, metav1.GetOptions{})
	if err != nil {
		return err
	}

	logrus.Info("\n")
	logrus.Infof("Deleting configmap: %s namespace: %s...",
		configmap,
		namespace)

	// Double check configmap is removed
	for i := 0; i < c.NumberMaxOfAttemptsPerTask; i++ {
		_, err := Exists(c, configmap, namespace)
		if err != nil {
			logrus.Infof("Deleted configmap: %s namespace: %s",
				configmap,
				namespace)
			break
		}
		c.Clientset.CoreV1().ConfigMaps(namespace).Delete(
			context.TODO(),
			configmap,
			metav1.DeleteOptions{})

		time.Sleep(time.Duration(c.TimeoutTaksInSec) * time.Second)
	}

	return nil
}


// Exists will check if the configmap exists or not
//
// Args:
//     - Pointer to a Client struct
//
// Returns:
//     string (namespace name) OR error type
//     
func Exists(c *client.Client, configmap string, namespace string) (string, error) {
	exists, err := c.Clientset.CoreV1().ConfigMaps(namespace).
		Get(context.TODO(), configmap, metav1.GetOptions{})
	if err != nil {
                return "", err
        }

        return exists.Name, nil
}

// Create creates a configmap using the values
// from the ConfigMap struct via the Client.Clientset 
//
// Args:
//    ConfigMap - ConfigMap struct
//    Client  - Client strucut
//
//   Returns:
//      error or nil
func Create(c *client.Client, cm *Instance) error {
        configmap := &v1.ConfigMap {
                ObjectMeta: metav1.ObjectMeta {
                        Name: cm.Name,
                        Namespace: cm.Namespace,
                },
		Data: map[string]string {
			cm.ConfigKey: cm.ConfigValue,
                },
        }

	logrus.Infof("\n")
	logrus.Infof("Creating configmap: %s namespace: %s",
		cm.Name,
		c.Namespace)

        _, err := c.Clientset.CoreV1().ConfigMaps(cm.Namespace).Create(
                context.TODO(),
		configmap,
                metav1.CreateOptions{})
        if err != nil {
                return err
        }

	logrus.Infof("Created configmap: %s namespace: %s",
		cm.Name,
		cm.Namespace)

        return nil
}
