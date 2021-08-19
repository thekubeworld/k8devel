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

	// +optional
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

	s := pvc.Instance{
		Name:        "mypvc",
		NamePrefix:  "pvc-",
		Namespace:   "default",
		ClaimSize:   "2G",
		AccessModes: "rwx",
		Annotations: map[string]string{"name": "mypvc"},
		// +optional
		//Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
		StorageClassName: "nfs-medogz",
		VolumeMode:       "filesystem",
	}

	pvc, err := pvc.Create(&c, &s)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("pvc [%v] created!\n", pvc.Name)
}
