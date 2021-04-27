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
	"errors"
	"fmt"
	"strings"
	"context"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DetectKubeProxyMode will detect kube-proxy mode
//
// Args:
//	- Pointer to a Client struct
//	- configmap
//	- namespace
//
// Returns:
//     string (namespace name) OR error type
//     
func DetectKubeProxyMode(c *Client, configmap string, namespace string) (string, error) {

	kproxyConfig, err := c.Clientset.CoreV1().ConfigMaps(namespace).Get(
		context.TODO(),
		configmap,
		metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if strings.Contains(fmt.Sprint(kproxyConfig.Data), "mode: iptables") {
                return "iptables", nil
        }

	if strings.Contains(fmt.Sprint(kproxyConfig.Data), "mode: ipvs") {
                return "ipvs", nil
        }

	return "", errors.New("unable to detect the kube-proxy mode")
}
