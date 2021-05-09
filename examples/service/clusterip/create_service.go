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
	"github.com/thekubeworld/k8devel/pkg/service"
)

func main() {
	newService := "newservice" // Put here the new Service name
	namespace := "default"     // Put here the namespace name

	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaksInSec = 2

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()
	c.Namespace = namespace

	s := service.Instance{
		Name:          newService,
		Namespace:     namespace,
		LabelKey:      "app",
		LabelValue:    "k8s",
		SelectorKey:   "app",
		SelectorValue: "k8s",
		ClusterIP:     "",
		Port:          80,
	}
	err := service.CreateClusterIP(&c, &s)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}

	IPService, err := service.GetIP(
		&c,
		newService,
		namespace)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Service %s namespace %s created\n", newService, namespace)
	fmt.Printf("ClusterIP Service IP: %s\n", IPService)
}
