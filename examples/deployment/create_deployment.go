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
	"github.com/thekubeworld/k8devel/pkg/deployment"
        "github.com/sirupsen/logrus"
)

func main() {
        k8devel.SetLogrusLogging()

	newDeployment := "newdeployment" // Put here the new deployment name
	namespace := "default" // Put here the namespace name

	c := client.Client{}
        c.NumberMaxOfAttemptsPerTask = 10
        c.TimeoutTaksInSec = 2

	// Connect to cluster from:
        //      - $HOME/kubeconfig (Linux)
        //      - os.Getenv("USERPROFILE") (Windows)
        c.Connect()

	d := deployment.Instance {
                Name: newDeployment,
                Namespace: namespace,
                Replicas: 1,
                LabelKey: "app",
                LabelValue: "nginxtesting",
        }

        d.Pod.Name = "nginx"
        d.Pod.Image = "nginx:1.14.2"
        d.Pod.ContainerPortName = "http"
        d.Pod.ContainerPortProtocol = "TCP"
        d.Pod.ContainerPort = 80

	err := deployment.Create(&c, &d)
        if err != nil {
                logrus.Fatal("exiting... failed to create: ", err)
        }

	logrus.Infof("Deployment %s namespace %s created!", newDeployment, namespace)
}
