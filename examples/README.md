To run any example, make sure the system has k8devel module for Go and execute `go run myexample.go`

```
$ export GOPATH=$(go env GOPATH)
$ cd $GOPATH/src
$ go get github.com/thekubeworld/k8devel

$ mkdir $HOME/tests && cd $HOME/tests
$ git clone https://github.com/thekubeworld/k8devel
$ cd k8devel/examples
$ go run namespace/create_namespace.go
```

External Projects:  
- [kubeproxy-testing](https://github.com/thekubeworld/kubeproxy-testing)
