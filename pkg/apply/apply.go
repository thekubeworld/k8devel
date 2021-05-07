package apply

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
	"bytes"
	"context"
	"fmt"

	"github.com/thekubeworld/k8devel/pkg/client"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

const yamlDelimiter = "---"

// decode will Decode data to object
//
// Args:
//     - date []byte
//
// Returns:
//     runtime.Object or error
//
func decode(data []byte) (runtime.Object, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(data, nil, nil)
	return obj, err
}

// YAML will go by the read object and create it via API
//
// Args:
//      - client struct
//	- yamlInput []bytes
//
// Returns:
//	- None
//
// TODO: Can we simplify the switch?
//
func YAML(c *client.Client, yamlInput []byte) []string {
	var output []string
	yamlFiles := bytes.Split(yamlInput, []byte(yamlDelimiter))

	for _, f := range yamlFiles {
		if len(f) == 0 || string(f) == "\n" {
			continue
		}

		obj, err := decode(f)
		if err != nil {
			continue
		}

		switch obj.(type) {
		case *v1.ServiceAccount:
			// If namespace not declared, use default
			namespace := ""
			if len(obj.(*v1.ServiceAccount).Namespace) == 0 {
				namespace = "default"
			}
			_, err = c.Clientset.CoreV1().ServiceAccounts(namespace).Create(
				context.TODO(),
				obj.(*v1.ServiceAccount),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("service ",
						obj.(*v1.ServiceAccount).Name,
						" created"))
			}
		case *v1.Namespace:
			_, err = c.Clientset.CoreV1().Namespaces().Create(
				context.TODO(),
				obj.(*v1.Namespace),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("namespace ",
						obj.(*v1.Namespace).Name,
						" created"))
			}
		case *v1.ConfigMap:
			_, err = c.Clientset.CoreV1().ConfigMaps(obj.(*v1.ConfigMap).Namespace).Create(
				context.TODO(),
				obj.(*v1.ConfigMap),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("configmap ",
						obj.(*v1.ConfigMap).Name,
						" created"))
			}
		case *appsv1.Deployment:
			_, err = c.Clientset.AppsV1().Deployments(obj.(*appsv1.Deployment).Namespace).Create(
				context.TODO(),
				obj.(*appsv1.Deployment),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("deployment ",
						obj.(*appsv1.Deployment).Name,
						" created"))
			}
		case *appsv1.DaemonSet:
			_, err = c.Clientset.AppsV1().DaemonSets(obj.(*appsv1.DaemonSet).Namespace).Create(
				context.TODO(),
				obj.(*appsv1.DaemonSet),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("daemonset ",
						obj.(*appsv1.DaemonSet).Name,
						" created"))
			}
		case *appsv1.StatefulSet:
			_, err = c.Clientset.AppsV1().StatefulSets(obj.(*appsv1.StatefulSet).Namespace).Create(
				context.TODO(),
				obj.(*appsv1.StatefulSet),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("statefulset ",
						obj.(*appsv1.StatefulSet).Name,
						" created"))
			}
		case *v1beta1.PodSecurityPolicy:
			_, err = c.Clientset.PolicyV1beta1().PodSecurityPolicies().Create(
				context.TODO(),
				obj.(*v1beta1.PodSecurityPolicy),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("podsecuritypolicy ",
						obj.(*v1beta1.PodSecurityPolicy).Name,
						" created"))
			}
		case *rbacv1.ClusterRole:
			_, err = c.Clientset.RbacV1().ClusterRoles().Create(
				context.TODO(),
				obj.(*rbacv1.ClusterRole),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("clusterrole ",
						obj.(*rbacv1.ClusterRole).Name,
						" created"))
			}
		case *rbacv1.RoleBinding:
			_, err = c.Clientset.RbacV1().RoleBindings(obj.(*rbacv1.RoleBinding).Namespace).Create(
				context.TODO(),
				obj.(*rbacv1.RoleBinding),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("rolebinding ",
						obj.(*rbacv1.RoleBinding).Name,
						" created"))
			}
		case *rbacv1.ClusterRoleBinding:
			_, err = c.Clientset.RbacV1().ClusterRoleBindings().Create(
				context.TODO(),
				obj.(*rbacv1.ClusterRoleBinding),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("clusterrolebinding ",
						obj.(*rbacv1.ClusterRoleBinding).Name,
						" created"))
			}
		case *rbacv1.Role:
			_, err = c.Clientset.RbacV1().Roles(obj.(*rbacv1.Role).Namespace).Create(
				context.TODO(),
				obj.(*rbacv1.Role),
				metav1.CreateOptions{})
			if err != nil {
				output = append(output, fmt.Sprint(err))
			} else {
				output = append(
					output,
					fmt.Sprint("role ",
						obj.(*rbacv1.Role).Name,
						" created"))
			}
		default:
			output = append(
				output,
				"error, unknown object kind for applying, verify yaml provided...")
		}

	}
	return output
}
