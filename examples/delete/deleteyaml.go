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
	"io/ioutil"
	"os"

	"github.com/thekubeworld/k8devel/pkg/client"
	delcmd "github.com/thekubeworld/k8devel/pkg/delete"
)

func main() {
	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaksInSec = 2

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	yamlInput, err := ioutil.ReadFile("file.yaml")
	if err != nil {
		os.Exit(1)
	}

	output := delcmd.YAML(&c, yamlInput)
	for _, i := range output {
		fmt.Println(i)
	}
}
