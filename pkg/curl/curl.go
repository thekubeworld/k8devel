package curl

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

import (
	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/pod"
)
// ExecuteHTTPReqInsideContainer will
// execute http request to a specific address
//
// Args:
//     
//      client struct
//      container name
//	namespace
//	URL
//
// Returns:
//	output as string or error
//
func ExecuteHTTPReqInsideContainer(c *client.Client,
                container string,
                namespace string,
		URL string) (string, error) {

	Cmd := []string{"curl"}
	Cmd = append(Cmd, URL)
        stdout, _, err := pod.ExecCmd(c,
              container,
              namespace,
              Cmd)
        if err != nil {
                return "", err
        }

        return string(stdout.Bytes()), nil
}