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
        "github.com/thekubeworld/k8devel/pkg/client"
        "github.com/thekubeworld/k8devel/pkg/logschema"
        "github.com/thekubeworld/k8devel/pkg/deployment"
        "github.com/thekubeworld/k8devel/pkg/util"
        "github.com/thekubeworld/k8devel/pkg/service"
        "github.com/thekubeworld/k8devel/pkg/namespace"
        "github.com/thekubeworld/k8devel/pkg/endpoint"

	"github.com/sirupsen/logrus"

)

func main() {
	logschema.SetLogrusLogging()
	// Initial set
        c := client.Client{}
	c.Namespace = "kptesting"
	c.NumberMaxOfAttemptsPerTask = 5
        c.TimeoutTaksInSec = 20

	// Connect to cluster from:
	//	- $HOME/kubeconfig (Linux)
	//	- os.Getenv("USERPROFILE") (Windows)
        c.Connect()

	KPTestNamespaceName := c.Namespace

	randStr, err := util.GenerateRandomString(6, "lower")
	if err != nil {
		logrus.Fatal(err)
	}
	KPTestServiceName := KPTestNamespaceName +
			"service" +
			randStr

	KPTestNginxDeploymentName := KPTestNamespaceName +
			"nginxdeployment" +
			randStr
	// END: kube-proxy variables

	// START: Namespace
	_, err = namespace.Exists(&c,
			KPTestNamespaceName)
	if err != nil {
		err = namespace.Create(&c,
			KPTestNamespaceName)
		if err != nil {
			logrus.Fatal("exiting... failed to create: ", err)
		}
	}
	// END: Namespace

	// START: Deployment
	d := deployment.Instance {
		Name: KPTestNginxDeploymentName,
		Namespace: KPTestNamespaceName,
		Replicas: 1,
		LabelKey: "app",
		LabelValue: "nginx",
	}

	d.Pod.Name = "nginx"
	d.Pod.Image = "nginx:1.14.2"
	d.Pod.ContainerPortName = "http"
	d.Pod.ContainerPortProtocol = "TCP"
	d.Pod.ContainerPort = 80

	err = deployment.Create(&c, &d)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	// END: Deployment

	// START: Service
	s := service.Instance {
		Name: KPTestServiceName,
		Namespace: KPTestNamespaceName,
		LabelKey: "k8sapp",
		LabelValue: "kproxy-testing",
		SelectorKey: "k8sapp",
		SelectorValue: "kproxy-testing",
		ClusterIP: "",
		Port: 80,
	}
	err = service.CreateClusterIP(&c, &s)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	IPService, err := service.GetIP(&c, KPTestServiceName, KPTestNamespaceName)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	logrus.Info(IPService)
	// END: Service

	// START: Endpoint
	e := endpoint.Instance {
		Name: "kproxy-service",
		Namespace: KPTestNamespaceName,
		IP: "172.16.0.80",
	}
	e.EndpointPort.Name = "http"
	e.EndpointPort.Port = 80
	e.EndpointPort.Protocol = "TCP"

	epoint, _ := endpoint.Exists(&c, &e)
	if epoint != "" {
		err = endpoint.Patch(&c, &e)
		if err != nil {
			logrus.Fatal("exiting... failed to update: ", err)
		}
	} else {
		err = endpoint.Create(&c, &e)
		if err != nil {
			logrus.Fatal("exiting... failed to create: ", err)
		}
	}
	endpoint.Show(&c, "kproxy-service", c.Namespace)
	endpoint.List(&c, &e)
	// END: Endpoint
}
