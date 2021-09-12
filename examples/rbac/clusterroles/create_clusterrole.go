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
	"github.com/thekubeworld/k8devel/pkg/rbac/clusterrole"
)

func main() {
	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 2

	c.Connect()

	apigroups := []string{""}
	resources := []string{"pods"}
	resourcesnames := []string{""}
	verbs := []string{"get", "watch", "list"}

	cr := clusterrole.Instance{
		Name:           "myclusterrole",
		LabelKey:       "app",
		LabelValue:     "foobar",
		APIGroups:      apigroups,
		Resources:      resources,
		ResourcesNames: resourcesnames,
		Verbs:          verbs,
	}

	fmt.Printf("creating cluster role: %s\n", cr.Name)
	err := clusterrole.Create(&c, &cr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
