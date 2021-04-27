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
        "github.com/thekubeworld/k8devel"
        "github.com/sirupsen/logrus"
)

func main() {
        k8devel.SetLogrusLogging()

	newNamespace := "newnamespace" // Put here your new namespace name
	c := k8devel.Client{}
        c.NumberMaxOfAttemptsPerTask = 10
        c.TimeoutTaksInSec = 2

	// Connect to cluster from:
        //      - $HOME/kubeconfig (Linux)
        //      - os.Getenv("USERPROFILE") (Windows)
        c.Connect()
	_, err := k8devel.ExistsNamespace(&c, newNamespace)
        if err != nil {
                err = k8devel.CreateNamespace(&c, newNamespace)
                if err != nil {
                        logrus.Fatal("exiting... failed to create: ", err)
                }
        }
	logrus.Infof("Namespace %s created!", newNamespace)
}
