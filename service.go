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
	"context"
	"time"
	v1 "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sirupsen/logrus"
)

// Service type refers to the Service object
type Service struct {
        Name string
        Namespace string
        LabelKey string
        LabelValue string
        Port int32
        SelectorKey string
        SelectorValue string
        ClusterIP string
	// Possible values for clusterIP:
	//   - None: headless service when proxying is not required
	//   - empty string or "": Auto Generated
	//   - ValidIP: Address
}

//func List(c *Client, serviceName string, namespace string) error {
//	servicelist, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
//	if err != nil {
//		return err
//	}
//}

// DeleteService deletes a service
//
// Args:
//	Client - client struct from the client module
//	service - Service Name
//	namespace - Namespace
//
//   Returns:
//      error or nil
func DeleteService(c *Client, service string, namespace string) error {
	_, err := c.Clientset.CoreV1().Services(namespace).
		Get(context.TODO(), service, metav1.GetOptions{})
	if err != nil {
		return err
	}

	logrus.Info("\n")
	logrus.Infof("Deleting service: %s namespace: %s...",
		service,
		namespace)

	// Double check service is removed
	for i := 0; i < c.NumberMaxOfAttemptsPerTask; i++ {
		_, err := ExistsService(c, service, namespace)
		if err != nil {
			logrus.Infof("Deleted service: %s namespace: %s",
				service,
				namespace)
			break
		}
		c.Clientset.CoreV1().Services(namespace).Delete(
			context.TODO(),
			service,
			metav1.DeleteOptions{})

		time.Sleep(time.Duration(c.TimeoutTaksInSec) * time.Second)
	}

	return nil
}

// GetIPFromService will return the pod IP address
//
// Args:
//
//      - Client struct from client module
//      - pod name
//      - namespace
//
// Returns:
//      - the IP as string or error
func GetIPFromService(c *Client,
                svcName string,
                nameSpace string) (string, error) {

        svc, err := c.Clientset.CoreV1().Services(nameSpace).Get(
                context.TODO(),
                svcName,
                metav1.GetOptions{})
        if err != nil {
                return "", err
        }
        return svc.Spec.ClusterIP, nil

}

// ExistsService will check if the service exists or not
//
// Args:
//     - Pointer to a Client struct
//
// Returns:
//     string (namespace name) OR error type
//     
func ExistsService(c *Client, service string, namespace string) (string, error) {
	exists, err := c.Clientset.CoreV1().Services(namespace).
		Get(context.TODO(), service, metav1.GetOptions{})
	if err != nil {
                return "", err
        }

        return exists.Name, nil
}

// CreateClusterIPService creates a service using the values
// from the Service struct via the Client.Clientset 
//
// Args:
//    Service - Service struct
//    Client  - Client strucut
//
//   Returns:
//      error or nil
func CreateClusterIPService(c *Client, s *Service) error {
        service := &v1.Service {
                ObjectMeta: metav1.ObjectMeta {
                        Name: s.Name,
                        Namespace: s.Namespace,
                        Labels: map[string]string {
                                s.LabelKey: s.LabelValue,
                        },
                },
                Spec: v1.ServiceSpec {
                        Ports: []v1.ServicePort{
                        {
                                Port: s.Port,
                        },
                },
                Selector: map[string]string {
                        s.SelectorKey: s.SelectorValue,
                },
                        ClusterIP: s.ClusterIP,
                },
        }

	logrus.Infof("\n")
	logrus.Infof("Creating service: %s namespace: %s",
		s.Name,
		c.Namespace)

        _, err := c.Clientset.CoreV1().Services(s.Namespace).Create(
                context.TODO(),
		service,
                metav1.CreateOptions{})
        if err != nil {
                return err
        }

	logrus.Infof("Created service: %s namespace: %s",
		s.Name,
		s.Namespace)

        return nil
}
