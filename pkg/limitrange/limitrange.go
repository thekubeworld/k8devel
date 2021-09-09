package limitrange

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
	"github.com/thekubeworld/k8devel/pkg/util"
)

// Instance type refers to the Pod object
type Instance struct {
	Name                 string
	Namespace            string
	LabelKey             string
	LabelValue           string
	LimitType            string
	Min                  [3]string // +optional
	Max                  [3]string // +optional
	Default              [3]string // +optional
	DefaultRequest       [3]string // +optional
	MaxLimitRequestRatio [3]string // +optional
}

// Delete will delete a LimitRange
//
// Args:
//	- Pointer to a Client struct
//	- namespace name
//	- limitrange name
//
// Returns:
//     error or nil
//
func Delete(c *client.Client, namespace string, limitrange string) error {
	err := c.Clientset.CoreV1().LimitRanges(namespace).Delete(
		context.TODO(),
		limitrange,
		metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// List will list all LimitRange
//
// Args:
//
//      - Client struct from client module
//
// Returns:
//      - pointer v1.NamespaceList or error
func List(c *client.Client, namespace string) (*v1.LimitRangeList, error) {
	limitRanges, err := c.Clientset.CoreV1().LimitRanges(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return limitRanges, nil
}

// Create will create a LimitRange
//
// Args:
//     - Pointer to a Client struct
//
// Returns:
//     error or nil
//
func Create(c *client.Client, l *Instance) error {
	limitType, err := util.DetectLimitType(l.LimitType)
	if err != nil {
		return err
	}

	lrange := &v1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.Name,
			Namespace: l.Namespace,
			Labels: map[string]string{
				l.LabelKey: l.LabelValue,
			},
		},
		Spec: v1.LimitRangeSpec{
			Limits: []v1.LimitRangeItem{
				{
					Type:                 limitType,
					Min:                  util.GetResourceList(l.Min[0], l.Min[1], l.Min[2]),
					Max:                  util.GetResourceList(l.Max[0], l.Max[1], l.Max[2]),
					Default:              util.GetResourceList(l.Default[0], l.Default[1], l.Default[2]),
					DefaultRequest:       util.GetResourceList(l.DefaultRequest[0], l.DefaultRequest[1], l.DefaultRequest[2]),
					MaxLimitRequestRatio: util.GetResourceList(l.MaxLimitRequestRatio[0], l.MaxLimitRequestRatio[1], l.MaxLimitRequestRatio[2]),
				},
			},
		},
	}
	lrange, err = c.Clientset.CoreV1().LimitRanges(l.Namespace).Create(
		context.TODO(),
		lrange,
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
//	- limitrange name
//
// Returns:
//     string (namespace name) OR error type
//
func Exists(c *client.Client, namespace string, limitrange string) (string, error) {
	exists, err := c.Clientset.CoreV1().LimitRanges(namespace).Get(
		context.TODO(),
		limitrange,
		metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return exists.Name, nil
}
