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
	//"os/signal"
	"strconv"
	"sync"
	//"syscall"
	"time"
)

var totalSec float64
var totalMinutes float64
var totalHour float64

var sumSec []float64
var sumMin []float64
var sumHour []float64
var nsSlice []string

var numberNamespaces = 10
var numberPods = 100
var totalPodsRunning = 0

var imageSource = "docker.io/nginx"
var c = client.Client{}

// generatePod create a pod and compare the time spend
// between the creation and the pod is in running state
//
// Args:
//      Client - struct from client module
//      podName - pod name
//      nsName - namespace name
//
func generatePod(c *client.Client, podName string, nsName string, wg *sync.WaitGroup) {
	defer wg.Done()

	p := pod.Instance{
		Name:            podName,
		Namespace:       nsName,
		Image:           imageSource,
		LabelKey:        "app",
		ImagePullPolicy: "ifnotpresent",
		LabelValue:      "podTest",
	}

	timeNow := time.Now()
	fmt.Printf("creating pod %s in namespace %s\n", podName, nsName)
	err := pod.CreateWaitRunningState(c, &p)
	//if err != nil {
	//	fmt.Printf("%s\n", err)
	//	os.Exit(1)
	//}

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
	sumSec = append(sumSec, totalSec)
	sumMin = append(sumMin, totalMinutes)
	sumHour = append(sumHour, totalHour)
	totalPodsRunning = totalPodsRunning + 1
	fmt.Printf("TOTAL NUMBER OF PODS RUNNING: %v\n", totalPodsRunning)
	fmt.Printf("TIME NOW: %v\n", time.Now().Format("2006-01-02 3:4:5 PM"))

	totalHour = 0
	totalMinutes = 0
	totalSec = 0
}

// createNamespace create a namespace
//
// Args:
//      Client - struct from client module
//      nsName - namespace name
//
// Return:
// 	error or nil
func createNamespaces(c *client.Client, nsName string, wgNs *sync.WaitGroup) {
	defer wgNs.Done()
	err := namespace.Create(c, nsName)
	if err != nil {
		fmt.Println("Failed to create namespace...")
		os.Exit(1)
	}
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

func generateNamespaces(wgNs *sync.WaitGroup) {
	nsName := ""
	for i := 1; i <= numberNamespaces; i++ {
		nsName, _ = util.GenerateRandomString(6, "lower")
		wgNs.Add(1)
		go createNamespaces(&c, nsName, wgNs)
		nsSlice = append(nsSlice, nsName)
	}
}

// main
func main() {
	fmt.Printf("REPORT GENERATED AT: %v\n", time.Now().Format("2006-01-02 3:4:5 PM"))

	c.NumberMaxOfAttemptsPerTask = 10
	//c.TimeoutTaskInSec = 1200  // 20 min
	//c.TimeoutTaskInSec = 10800 // 3 hours
	c.TimeoutTaskInSec = 3600 // 1 hours

	// Connect to cluster from:
	//      - $HOME/kubeconfig (Linux)
	//      - os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	var wgPods sync.WaitGroup
	var wgNamespaces sync.WaitGroup

	generateNamespaces(&wgNamespaces)
	wgNamespaces.Wait()
	fmt.Printf("created %s namespaces\n", nsSlice)

	fmt.Printf("Creating pods...\n")
	for _, nsName := range nsSlice {
		for j := 1; j <= numberPods; j++ {
			wgPods.Add(1)
			go generatePod(&c,
				"pod"+strconv.Itoa(j),
				nsName,
				&wgPods)

		}
	}
	wgPods.Wait()

	fmt.Printf("\n🏁 Summary 🏁\n")
	fmt.Printf("-----------------------------\n")
	fmt.Printf("Namespaces created: %v\n", numberNamespaces)
	fmt.Printf("Pods per Namespaces created: %v\n", numberPods)
	fmt.Printf("Hours: %v\n", sumResults(sumHour))
	fmt.Printf("Minutes: %v\n", sumResults(sumMin))
	fmt.Printf("Seconds: %v\n", sumResults(sumSec))
	fmt.Println("=========================")

	fmt.Println("Absolute time:")
	sec := int(sumResults(sumSec))
	min := 0
	if sec > 60 {
		min := sec / 60
		fmt.Printf("Minutes: %v\n", min)
	} else {
		fmt.Printf("Seconds: %v\n", sec)
	}

	if min > 60 {
		hour := min / 60
		fmt.Printf("Hour: %v\n", hour)
	}

	fmt.Printf("\nCleaning created objects during the tests...\n")
	cleanup()

	fmt.Println("done!")

}

func cleanup() {
	for _, n := range nsSlice {
		fmt.Printf("Deleting %s\n", n)
		err := namespace.Delete(&c, n)
		if err != nil {
			fmt.Printf("cannot delete namespace %s\n", n)
			os.Exit(1)
		}
	}
}

/*
func init() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cleanup()
		os.Exit(1)
	}()
}*/
