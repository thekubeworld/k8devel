package firewall

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
	"errors"
	"os"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/util"
)

// Save will save the current state of firewall
//
// Args:
//	client struct
//	firewallMode - iptables or ipvs
//	container name
//	namespace
//
// Returns:
//	file object
//	filesystem which triggered this method
//
func Save(c *client.Client,
	firewallMode string,
	container string,
	namespace string) (*os.File, error) {

	var cmdSave []string
	if firewallMode == "iptables" {
		cmdSave = append(cmdSave, "iptables-save")
	} else if firewallMode == "ipvs" {
		cmdSave = append(cmdSave, "ipvsadm")
		cmdSave = append(cmdSave, "--save")
	} else {
		return nil, errors.New("unknown firewall mode")
	}

	fileRef, err := util.CreateTempFile(os.TempDir(), "firewall")
	if err != nil {
		return nil, err
	}

	stdout, _, err := pod.ExecCmd(c,
		container,
		namespace,
		cmdSave)
	if err != nil {
		return nil, err
	}

	fileRef.Write(stdout.Bytes())
	fileRef.Sync()

	return fileRef, nil
}
