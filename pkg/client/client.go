package client

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
	"os"
	"path/filepath"

	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

// Client struct contains all info about client
// to connect into the cluster
type Client struct {
	Clientset                  *kubernetes.Clientset
	Restclientset              *kubernetes.Clientset
	Namespace                  string
	Restconfig                 *rest.Config
	Kubeconfig                 clientcmd.ClientConfig
	TimeoutTaskInSec           int
	NumberMaxOfAttemptsPerTask int
	QPS                        float32 // Queries per second
	Burst                      int     // Suddently increase of call to API

	// TODO: remove NumberMaxOfAttemptsPerTask and add some Pool mechanism
	// for modules that still use it. that
}

// Connect will connect to specific Cluster
// read from kubeconfig
//
// Args:
//
// Returns:
//   - Client struct
func (client *Client) Connect() (*Client, error) {
	// TODO: Users can specify by dynamic the HOME for kubeconfig
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = os.Getenv("USERPROFILE") // windows
	}

	configPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", configPath)

	// QPS and Burst settings helps users that want to increase
	// the number of call per second against API Server, common
	// warning message:
	//
	// ... request.go:XX Waited for 1.19126963s
	// due to client-side throttling, not priority and fairness,
	// request: POST:https://foobar:6443/api/v1...
	//
	// Examples of increase:
	// QPS: 6000 Bust: 30000
	// QPS: 1000 Bust: 1000 etc
	if client.QPS > 0 {
		config.QPS = client.QPS
	}

	if client.Burst > 0 {
		config.Burst = client.Burst
	}

	if err != nil {
		return nil, err
	}

	client.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client.KubeClientFromConfig()
	return client, nil
}

// KubeClientFromConfig will provide the REST interface
// to cluster
//
// Args:
//
// Returns:
//   - Client struct or error
func (client *Client) KubeClientFromConfig() (*Client, error) {
	var err error

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	client.Kubeconfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		configOverrides)

	client.Restconfig, err = client.Kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client.Restclientset, err = kubernetes.NewForConfig(client.Restconfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
