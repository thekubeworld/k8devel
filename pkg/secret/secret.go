package secret

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
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/thekubeworld/k8devel/pkg/client"
)

type Instance struct {
	Name      string
	Namespace string
	Type      string
	Key       string
	Value     string
	// "k8s.io/api/core/v1/types.go"
	// Type:
	// Opaque				arbitrary user-defined data
	// kubernetes.io/service-account-token	service account token
	// kubernetes.io/dockercfg		serialized ~/.dockercfg file
	// kubernetes.io/dockerconfigjson	serialized ~/.docker/config.json file
	// kubernetes.io/basic-auth		credentials for basic authentication
	// kubernetes.io/ssh-auth		credentials for SSH authentication
	// kubernetes.io/tls			data for a TLS client or server
	// bootstrap.kubernetes.io/token	bootstrap token dataExample: "Opaque" "kubernetes.io/dockerconfigjson"
}

// Create will create a secret
//
// Args:
//     - Pointer to a Client struct
//     - Point to the Instance struct
//
// Returns:
//     error or nil
//
func Create(c *client.Client, i *Instance) error {
	secretType, err := detectSecretType(i.Type)
	if err != nil {
		return err
	}

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.Name,
			Namespace: i.Namespace,
		},
		Data: map[string][]byte{
			i.Key: []byte(i.Value),
		},
		Type: secretType,
	}
	_, err = c.Clientset.CoreV1().Secrets(i.Namespace).Create(
		context.TODO(),
		&secret,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// detectSecretType will detect the type provided as string
// and convert to v1.Secret.Type
//
// Args:
//     - string
//
// Returns:
//     v1.Secret.Type or error
//
func detectSecretType(s string) (v1.SecretType, error) {
	switch strings.ToLower(s) {
	case "opaque":
		return v1.SecretTypeOpaque, nil
	case "kubernetes.io/basic-auth":
		return v1.SecretTypeBasicAuth, nil
	case "kubernetes.io/tls":
		return v1.SecretTypeTLS, nil
	case "kubernetes.io/ssh-auth":
		return v1.SecretTypeSSHAuth, nil
	case "kubernetes.io/service-account-token":
		return v1.SecretTypeServiceAccountToken, nil
	case "kubernetes.io/dockercfg":
		return v1.SecretTypeDockercfg, nil
	case "kubernetes.io/dockerconfigjson":
		return v1.SecretTypeDockerConfigJson, nil
	}
	return "", errors.New("unknown secretType yet\n")
}

// Exists will check if thsee secret  exists or not
//
// Args:
//     - Pointer to a Client struct
//     - secret name
//      - namespace name
//
// Returns:
//     string (namespace name) OR error type
//
func Exists(c *client.Client, secretname string, namespace string) (string, error) {
	exists, err := c.Clientset.CoreV1().Secrets(namespace).Get(
		context.TODO(),
		secretname,
		metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return exists.Name, nil
}
