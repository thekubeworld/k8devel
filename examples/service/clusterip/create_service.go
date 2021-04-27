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
        "github.com/thekubeworld/k8devel"
        "github.com/sirupsen/logrus"
)

func main() {
        k8devel.SetLogrusLogging()

	newService := "newservice" // Put here the new Service name
	namespace := "default" // Put here the namespace name

	c := k8devel.Client{}
        c.NumberMaxOfAttemptsPerTask = 10
        c.TimeoutTaksInSec = 2

	// Connect to cluster from:
        //      - $HOME/kubeconfig (Linux)
        //      - os.Getenv("USERPROFILE") (Windows)
        c.Connect()

	s := k8devel.Service {
		Name: newService,
		Namespace: namespace,
		LabelKey: "app",
		LabelValue: "k8s",
		SelectorKey: "app",
		SelectorValue: "k8s",
		ClusterIP: "",
		Port: 80,
	}
	err := k8devel.CreateService(&c, &s)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}

	IPService, err := k8devel.GetIPFromService(
		&c,
		newService,
		namespace)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	logrus.Infof("Service %s namespace %s created!", newService, namespace)
	logrus.Infof("ClusterIP Service IP: %s", IPService)
}
