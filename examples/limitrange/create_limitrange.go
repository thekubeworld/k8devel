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
	"github.com/thekubeworld/k8devel/pkg/limitrange"
)

func main() {
	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 2

	c.Connect()

	// Setting the limits
	ltype := "container" //pod, container, persistentvolumeclaim
	min := [3]string{"5m", "100Mi", "5Gi"}
	max := [3]string{"100m", "10000Mi", "10Gi"}
	defaultRequest := [3]string{"10m", "200Mi", ""}
	defaultLimit := [3]string{"50m", "500Mi", ""}
	maxLimitRequestRatio := [3]string{"10", "", ""}

	l := limitrange.Instance{
		Name:                 "mylimitrange",
		Namespace:            "default",
		LabelKey:             "app",
		LabelValue:           "foobar",
		LimitType:            ltype,
		Min:                  min,
		Max:                  max,
		Default:              defaultLimit,
		DefaultRequest:       defaultRequest,
		MaxLimitRequestRatio: maxLimitRequestRatio,
	}

	fmt.Printf("creating limit range: %s\n", l.Name)
	err := limitrange.Create(&c, &l)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
