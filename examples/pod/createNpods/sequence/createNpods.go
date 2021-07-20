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
	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/util"
	"math"
	"os"
	"strconv"
	"time"
)

var totalSec float64
var totalMinutes float64
var totalHour float64

var sumSec []float64
var sumMin []float64
var sumHour []float64

var numberNamespaces = 10
var numberPods = 100

var imageSource = "docker.io/nginx"

// generatePod create a pod and compare the time spend
// between the creation and the pod is in running state
//
// Args:
//      Client - struct from client module
//      podName - pod name
//      nsName - namespace name
//
func generatePod(c *client.Client, podName string, nsName string) {
	p := pod.Instance{
		Name:            podName,
		Namespace:       nsName,
		Image:           imageSource,
		LabelKey:        "app",
		ImagePullPolicy: "ifnotpresent",
		LabelValue:      "podTest",
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

	fmt.Printf("\n- %s is created and responsive in namespace %s ✅\n", p.Name, p.Namespace)
	fmt.Printf("- image used: %s\n", imageSource)

	fmt.Println("  took:", hour, "hours",
		minutes, "minutes",
		seconds, "seconds")

}

// createNamespace create a namespace
//
// Args:
//      Client - struct from client module
//      nsName - namespace name
//
// Return:
// 	error or nil
func createNamespace(c *client.Client, nsName string) error {
	err := namespace.Create(c, nsName)
	if err != nil {
		return err
	}
	return nil
}

// sumResults will sum results from exections
//
// Args:
//      []float64 - slice with all results to be sum
//
// Return:
// 	float64
func sumResults(sumResult []float64) float64 {
	result := 0.0
	for _, s := range sumResult {
		result += s
	}
	return result
}

// main
func main() {
	fmt.Printf("REPORT GENERATED AT: %v\n", time.Now().Format("2006-01-02 3:4:5 PM"))

	c := client.Client{}
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaskInSec = 1200 // 20 min

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	nsName := ""
	for i := 0; i < numberNamespaces; i++ {
		nsName, _ = util.GenerateRandomString(6, "lower")
		err := createNamespace(&c, nsName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\nNamespace %s created ✅\n", nsName)
		fmt.Printf("Creating pods and waiting for running state ⏳")
		for i := 0; i < numberPods; i++ {
			generatePod(&c,
				"pod"+strconv.Itoa(i),
				nsName)
		}

		sumSec = append(sumSec, totalSec)
		sumMin = append(sumMin, totalMinutes)
		sumHour = append(sumHour, totalHour)

		totalHour = 0
		totalMinutes = 0
		totalSec = 0

		err = namespace.Delete(&c, nsName)
		if err != nil {
			fmt.Println("cannot delete namespace: %s\n", nsName)
			os.Exit(1)
		}
	}

	fmt.Printf("\n🏁 Summary 🏁\n")
	fmt.Printf("-----------------------------\n")
	fmt.Printf("Namespaces created: %v\n", numberNamespaces)
	fmt.Printf("Pods per Namespaces created: %v\n", numberPods)
	fmt.Printf("Hours: %v\n", sumResults(sumHour))
	fmt.Printf("Minutes: %v\n", sumResults(sumMin))
	fmt.Printf("Seconds: %v\n", sumResults(sumSec))

}
