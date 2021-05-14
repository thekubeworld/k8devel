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
	"github.com/thekubeworld/k8devel/pkg/cronjob"
)

func main() {
	newCronJob := "mycronjob" // Put here the new cronjob name
	namespace := "default"    // Put here the namespace name

	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 2

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	j := cronjob.Instance{
		Name:                       newCronJob,
		Namespace:                  namespace,
		Schedule:                   "*/1 * * * *", // At every minute
		ConcurrencyPolicy:          "Allow",       // Allow, Forbid and Replace
		Parallelism:                int32(0),      // 0 means unlimited parallelism.
		Completions:                int32(1),
		RestartPolicy:              "Never", // OnFailure, Never, Always
		BackoffLimit:               6,       // default is 6
		SuccessfulJobsHistoryLimit: int32(1),
		FailedJobsHistoryLimit:     int32(1),
		Command:                    []string{"/bin/true"}, // Pod will complete with success. It can false or sleep here too
	}

	command := []string{"ls", "-la"}
	j.Pod.Name = "ubuntu"
	j.Pod.Image = "ubuntu:latest"
	j.Pod.Command = command

	err := cronjob.Create(&c, &j)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("cronjob %s namespace %s created\n", newCronJob, namespace)
}
