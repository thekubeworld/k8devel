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

	// Setting the limits:
	//	pod, container or persistentvolumeclaim
	ltype := "container"

	// Min usage constraints on this kind by resource name.
	min := [3]string{
		"5m",    // cpu
		"100Mi", // memory
		"5Gi"}   // disk

	// Max usage constraints on this kind by resource name.
	max := [3]string{
		"100m",    // cpu
		"10000Mi", // memory
		"10Gi"}    // disk

	// DefaultRequest is the default resource requirement
	// request value by resource name if resource request
	// is omitted.
	defaultRequest := [3]string{
		"10m",   // cpu
		"200Mi", // memory
		""}      // disk

	// Default resource requirement limit value by resource
	// name if resource limit is omitted.
	defaultLimit := [3]string{
		"50m",   // cpu
		"500Mi", // memory
		""}      // disk

	// MaxLimitRequestRatio if specified, the named resource
	// must have a request and limit that are both non-zero
	// where limit divided by request is less than or equal
	// to the enumerated value; this represents the max burst
	// for the named resource.
	maxLimitRequestRatio := [3]string{
		"10", // cpu
		"",   // memory
		""}   // disk

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
