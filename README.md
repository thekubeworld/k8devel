# k8devel

[![Go Report Card](https://goreportcard.com/badge/github.com/thekubeworld/k8devel)](https://goreportcard.com/report/github.com/thekubeworld/k8devel)
[![GoDoc](https://godoc.org/github.com/thekubeworld/k8devel?status.svg)](https://pkg.go.dev/github.com/thekubeworld/k8devel)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


**An Easy to use** Go **framework** for Kubernetes based on [kubernetes/client-go](https://github.com/kubernetes/client-go), see [examples](https://github.com/thekubeworld/k8devel/tree/main/examples) dir for a quick start.

How to test it ?

1. Download the module  
```
$ GO111MODULE=on go get github.com/thekubeworld/k8devel
```

2. Download the source
```
$ git clone https://github.com/thekubeworld/k8devel.git
$ cd k8devel/examples/pod
```

3. Build and Run
```
$ go build pod.go
$ ./pod
Pod mytesting namespace default created!
```
