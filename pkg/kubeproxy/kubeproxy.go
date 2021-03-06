package kubeproxy

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
	"errors"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/distro/debian/apt"
	"github.com/thekubeworld/k8devel/pkg/distro/debian/dpkg"
	"github.com/thekubeworld/k8devel/pkg/firewall"
	"github.com/thekubeworld/k8devel/pkg/pod"
)

// SaveCurrentFirewallState will save the current state
// from the kubeproxy pod
//
// Args:
//	- Pointer to a Client struct
//	- podname A substring of kube-proxy pod name
//	- namespace
//
// Returns:
//	filename storing the firewall rules or error
//
func SaveCurrentFirewallState(c *client.Client,
	configmapname string,
	containerName string,
	namespace string) (string, error) {

	mode, err := DetectKubeProxyMode(c,
		configmapname,
		containerName,
		namespace)
	if err != nil {
		return "", err
	}

	containerName, err = FindKubeProxyPod(c, containerName, namespace)
	if err != nil {
		return "", err
	}

	if mode == "ipvs" {
		// check if ipvsadm exists, if not install it
		_, err := dpkg.CheckPackageInstalled(c, containerName, namespace, "ipvsadm")
		if err != nil {
			// apt update
			_, err = apt.UpdateInsidePod(
				c,
				containerName,
				namespace)
			if err != nil {
				return "", err
			}

			// apt install ipvsadm
			_, err = apt.InstallPackageInsidePod(
				c,
				containerName,
				namespace,
				"ipvsadm")
			if err != nil {
				return "", err
			}
		}
	}
	filesaved, err := firewall.Save(c,
		mode,
		containerName,
		namespace)
	if err != nil {
		return "", err
	}

	return filesaved.Name(), nil
}

// FindKubeProxyPod will return one of the daemonsets
// pods names for kubeproxy so we can connect to pod
// and execute commands or other actions
//
// Args:
//	- Pointer to a Client struct
//	- containerName A substring of kube-proxy pod name
//	- namespace
//
// Returns:
//     the first kube-proxy pod found from the daemonsets
//	or error
//
func FindKubeProxyPod(c *client.Client,
	containerName string,
	namespace string) (string, error) {
	// Validation
	kyPods, kyNumberPods := pod.FindPodsWithNameContains(c,
		containerName, namespace)
	if kyNumberPods < 0 {
		return "", errors.New(
			"exiting... unable to find kube-proxy pod")
	}
	return kyPods[0], nil
}

// DetectKubeProxyMode will detect kube-proxy mode
//
// Args:
//	- Pointer to a Client struct
//	- configmapname
//	- containerName
//	- namespace
//
// Returns:
//     string (ipvs or iptables) OR error type
//
func DetectKubeProxyMode(c *client.Client,
	configmapname string,
	containerName string,
	namespace string) (string, error) {

	// make sure we find at least one kube-proxy pod
	_, err := FindKubeProxyPod(c, containerName, namespace)
	if err != nil {
		return "", err
	}

	// Get configmapname from kube-proxy
	kproxyConfig, err := c.Clientset.CoreV1().ConfigMaps(namespace).Get(
		context.TODO(),
		configmapname,
		metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Detect if it's iptables
	if strings.Contains(
		fmt.Sprint(kproxyConfig.Data),
		"mode: iptables") {
		return "iptables", nil
	}

	// Detect if it's ipvs
	if strings.Contains(
		fmt.Sprint(kproxyConfig.Data), "mode: ipvs") {

		return "ipvs", nil
	}

	return "", errors.New("unable to detect the kube-proxy mode")
}
