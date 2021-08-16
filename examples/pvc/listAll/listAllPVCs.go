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
	"time"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/pvc"
)

func main() {
	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 2

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	pvc, err := pvc.List(&c, "default")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	t1 := time.Now()

	// Get duration.
	for _, p := range pvc.Items {
		fmt.Printf("Name: %s\n", p.ObjectMeta.Name)
		fmt.Printf("Namespace: %s\n", p.ObjectMeta.Namespace)
		fmt.Printf("Creation time: %s\n", p.ObjectMeta.CreationTimestamp.Time)
		fmt.Printf("Age: %vd\n", int(t1.Sub(p.ObjectMeta.CreationTimestamp.Time).Hours()/24))
		fmt.Printf("Status: %v\n", p.Status.Phase)
		fmt.Printf("Capacity: %v\n", p.Spec.Resources.Requests.Storage().String())
		fmt.Printf("Volume: %v\n", p.Spec.VolumeName)
		fmt.Printf("AccessModes: %s\n", p.Status.AccessModes[0])
		fmt.Printf("StorageClassName: %v\n", *p.Spec.StorageClassName)
		fmt.Printf("\n")
	}
}
