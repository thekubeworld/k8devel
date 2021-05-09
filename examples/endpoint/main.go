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

package main

import (
	"fmt"
	"os"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/deployment"
	"github.com/thekubeworld/k8devel/pkg/endpoint"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/service"
	"github.com/thekubeworld/k8devel/pkg/util"
)

func main() {
	// Initial set
	c := client.Client{}
	c.Namespace = "kptesting"
	c.NumberMaxOfAttemptsPerTask = 5
	c.TimeoutTaskInSec = 20

	// Connect to cluster from:
	//	- $HOME/kubeconfig (Linux)
	//	- os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	randStr, err := util.GenerateRandomString(6, "lower")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// kube-proxy variables
	NamespaceName := c.Namespace + randStr
	ServiceName := "service" + randStr
	DeploymentName := "nginxdeployment" + randStr
	EndpointName := "kproxy-service" + randStr

	// START: Namespace
	err = namespace.Create(&c, NamespaceName)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}
	// END: Namespace

	// START: Deployment
	d := deployment.Instance{
		Name:       DeploymentName,
		Namespace:  NamespaceName,
		Replicas:   1,
		LabelKey:   "app",
		LabelValue: "nginx",
	}

	d.Pod.Name = "nginx"
	d.Pod.Image = "nginx:1.14.2"
	d.Pod.ContainerPortName = "http"
	d.Pod.ContainerPortProtocol = "TCP"
	d.Pod.ContainerPort = 80

	err = deployment.Create(&c, &d)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}
	// END: Deployment

	// START: Service
	s := service.Instance{
		Name:          ServiceName,
		Namespace:     NamespaceName,
		LabelKey:      "k8sapp",
		LabelValue:    "kproxy-testing",
		SelectorKey:   "k8sapp",
		SelectorValue: "kproxy-testing",
		ClusterIP:     "",
		Port:          80,
	}
	err = service.CreateClusterIP(&c, &s)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}

	IPService, err := service.GetIP(&c, ServiceName, NamespaceName)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("IP Service %s\n", IPService)
	// END: Service

	// START: Endpoint
	e := endpoint.Instance{
		Name:      EndpointName,
		Namespace: NamespaceName,
		IP:        "172.16.0.80",
	}
	e.EndpointPort.Name = "http"
	e.EndpointPort.Port = 80
	e.EndpointPort.Protocol = "TCP"

	epoint, _ := endpoint.Exists(&c, &e)
	if len(epoint) > 0 {
		err = endpoint.Patch(&c, &e)
		if err != nil {
			fmt.Printf("exiting... failed to update: %s\n", err)
			os.Exit(1)
		}
	} else {
		err = endpoint.Create(&c, &e)
		if err != nil {
			fmt.Printf("exiting... failed to create: %s\n", err)
			os.Exit(1)
		}
	}
	endpoint.Show(&c, EndpointName, c.Namespace)
	endpoint.List(&c, &e)
	// END: Endpoint
}
