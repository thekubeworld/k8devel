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

	"github.com/thekubeworld/k8devel/pkg/base64"
	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/loadbalancer/metallb"
	"github.com/thekubeworld/k8devel/pkg/secret"
)

func main() {
	c := client.Client{}

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	err := metallb.Deploy(&c, "v0.9.6")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("metallb deployed...\n")

	base64Str, _ := base64.GenerateRandomString(128)
	metallbSecret := secret.Instance{
		Name:      "memberlist",
		Namespace: "metallb-system",
		Type:      "Opaque",
		Key:       "secretkey",
		Value:     base64Str,
	}
	err = metallb.CreateSecret(&c, &metallbSecret)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("metallb secret created...\n")

	// Metallb Config
	conf := metallb.InstanceConfig{
		Name:                 "config",
		Namespace:            "metallb-system",
		ConfigName:           "config",
		AddressPoolName:      "default",
		AddressPoolProtocol:  "layer2",
		AddressPoolAddresses: "172.17.255.1-172.17.255.250",
	}
	err = metallb.CreateConfig(&c, &conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("metallb created configmap %s namespace %s\n", conf.Name, conf.Namespace)
	fmt.Printf("mettalb deployed successfully\n")
}
