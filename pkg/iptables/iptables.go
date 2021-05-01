package iptables

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
	"os"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/util"
)

// IPTables type refers to the iptables object
type Instance struct {
	ReadNatTable []string
	ReadFilterTable []string
	ReadNatTableKubeServices []string
	Save []string
}

// LoadPreDefinedIPTables loads several pre-defined
// iptables commands.
//
// Args:
//	None
//
// Returns:
//	Filled IPTablesIPTables struct with commands
//
func LoadPreDefinedCommands() Instance {
	IPTablesCmd := Instance{}

	IPTablesCmd.ReadNatTable = []string {
		"iptables",
		"-w",
		"-t",
		"nat",
		"-L",
		"-n",
		"-v",
	}

	IPTablesCmd.ReadFilterTable = []string {
		"iptables",
		"-w",
		"-t",
		"filter",
		"-L",
		"-n",
		"-v",
	}

	IPTablesCmd.ReadNatTableKubeServices = []string {
		"iptables",
		"-w",
		"-L",
		"-n",
		"-v",
		"KUBE-SERVICES",
		"-t",
		"nat",
	}

	IPTablesCmd.Save = []string {
		"iptables-save",
	}

	return IPTablesCmd
}

// IPTablesSave will save the current state of NAT table
// from the container provider in the parameter
//
// Args:
//	container name	
//	namespace
//
// Returns:
//	file object
//	filesystem which triggered this method
//
func Save(c *client.Client, i *Instance, container string, namespace string)(*os.File, error) {
	fileRef, err := util.CreateTempFile(os.TempDir(), "iptables")
        if err != nil {
		return nil, err
        }

	stdout, _, err := pod.ExecCmd(c,
              container,
              namespace,
              i.Save)
        if err != nil {
		return nil, err
        }

        fileRef.Write(stdout.Bytes())
        fileRef.Sync()

	return fileRef, nil
}
