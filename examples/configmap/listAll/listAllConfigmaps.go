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
	"github.com/thekubeworld/k8devel/pkg/configmap"
	"github.com/thekubeworld/k8devel/pkg/emoji"
)

func main() {
	e := emoji.LoadEmojis()
	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 2

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	configMapList, err := configmap.ListAll(&c)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	for _, cm := range configMapList.Items {
		fmt.Printf("Name %s\n", cm.ObjectMeta.Name)
		fmt.Printf("Namespace: %s\n", cm.ObjectMeta.Namespace)
		for _, d := range cm.Data {
			fmt.Printf("%s", d)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Number total if configmaps: %v %s %s\n",
		len(configMapList.Items),
		emoji.Show(e.Rocket),
		emoji.Show(e.Collision))
}
