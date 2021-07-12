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
	"sync"
	"time"
	"math"
	"strconv"
	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/util"
)

var totalSec float64
var totalMinutes float64
var totalHour float64

func generatePod(c *client.Client, wg *sync.WaitGroup, podName string, nsName string) {
	defer wg.Done()
	p := pod.Instance{
                Name:       podName,
                Namespace:  nsName,
                Image:      "nginx",
                LabelKey:   "app",
                LabelValue: "podTest",
        }

	timeNow := time.Now()
        err := pod.Create(c, &p)
        if err != nil {
                fmt.Printf("%s\n", err)
                os.Exit(1)
        }

	lastTime, err := pod.GetLastTimeConditionHappened(c,
		"Ready",
		podName,
		nsName)
	if err != nil {
		fmt.Println(err)
	}

	hour := lastTime.Sub(timeNow).Hours()
	hour, mf := math.Modf(hour)
	totalHour = totalHour + hour

	minutes := mf * 60
	minutes, sf := math.Modf(minutes)
	totalMinutes = totalMinutes + minutes

	seconds := sf * 60
	totalSec = totalSec + seconds
        fmt.Printf("\n- %s is created and responsive in namespace %s\n", p.Name, p.Namespace)

	fmt.Println("  took:", math.Abs(hour), "hours",
			math.Abs(minutes), "minutes",
			math.Abs(seconds), "seconds")
}

func createNamespace(c *client.Client, nsName string) error {
	err := namespace.Create(c, nsName)
        if err != nil {
                return err
        }
	return nil
}

func main() {

	fmt.Printf("REPORT GENERATED AT: %v\n\n",  time.Now().Format("2006-01-02 3:4:5 PM"))

	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 1200 // 20 min

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()


	nsName, _ := util.GenerateRandomString(6, "lower")

	err := createNamespace(&c, nsName)
        if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	maxpod := 5

	fmt.Printf("Generating %d pods in namespace %v...", maxpod, nsName)
	for i := 0; i <= 100; i++ {
		wg.Add(1)
		generatePod(&c,
			&wg,
			"pod" + strconv.Itoa(i),
			nsName)
	}

	wg.Wait()

	fmt.Printf("\nTotal time creating pods:\n")
	fmt.Println("  Hour:", totalHour,
			"Minutes:", totalHour,
			"Seconds:", totalSec, "\n")

	fmt.Printf("Deleting namespace: %s\n", nsName)
	//namespace.Delete(&c, nsName)
}
