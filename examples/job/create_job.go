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
	"github.com/thekubeworld/k8devel/pkg/job"
)

func main() {
	newJob := "mynewjob"   // Put here the new job name
	namespace := "default" // Put here the namespace name

	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 2

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	j := job.Instance{
		Name:          newJob,
		Namespace:     namespace,
		RestartPolicy: "Never", // Never or Always or OnFailure
		BackoffLimit:  6,       // default number is 6 attemps before it calls as failure
	}

	command := []string{"ls", "-la"}
	j.Pod.Name = "ubuntu"
	j.Pod.Image = "ubuntu:latest"
	j.Pod.Command = command

	err := job.Create(&c, &j)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Job %s namespace %s created\n", newJob, namespace)
}
