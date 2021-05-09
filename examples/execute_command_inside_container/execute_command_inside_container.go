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
	"github.com/thekubeworld/k8devel/pkg/pod"
)

func main() {
	containerName := "mytesting" // Put here the Pod name
	namespace := "default"       // Put here the namespace name
	cmd := []string{"ls", "-la"} // Put here the command to be executed inside container

	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 2

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	p := pod.Instance{
		Name:      containerName,
		Namespace: namespace,
		Image:     "nginx",
	}
	// POD Settings

	pod.Create(&c, &p)

	stdout, _, err := pod.ExecCmd(&c, containerName, namespace, cmd)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Output from command:\n")
	fmt.Printf("%s\n", stdout.String())
}
