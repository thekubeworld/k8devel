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
	"github.com/thekubeworld/k8devel/pkg/emoji"
	"github.com/thekubeworld/k8devel/pkg/secret"
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

	s := secret.Instance{
		Name:      "mysecret",
		Namespace: "default",
		Type:      "Opaque",
		Key:       "myuser",
		Value:     "mypass",
	}

	err := secret.Create(&c, &s)
	if err != nil {
		fmt.Printf("%s %s\n", emoji.Show(e.CrossMark), err)
		os.Exit(1)
	}
	fmt.Printf("%s secret created %s, namespace %s, type: %s\n",
		emoji.Show(e.Rocket),
		s.Name,
		s.Namespace,
		s.Type)
}
