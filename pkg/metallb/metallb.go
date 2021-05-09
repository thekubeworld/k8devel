package metallb

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
	"io/ioutil"
	"os"

	"github.com/thekubeworld/k8devel/pkg/apply"
	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/configmap"
	"github.com/thekubeworld/k8devel/pkg/secret"
	"github.com/thekubeworld/k8devel/pkg/util"
)

// InstanceConfig type refers to the Metallb Config
type InstanceConfig struct {
	Name                 string
	Namespace            string
	ConfigName           string
	AddressPoolName      string
	AddressPoolProtocol  string
	AddressPoolAddresses string
}

// Deploy will download yaml files from metallb and deploy in
// the environment
//
// Args:
//
//      - Client struct from client module
//	- version to deploy
//
// Returns:
//      - pointer v1.ConfigMapList or error
func Deploy(c *client.Client, version string) error {
	if len(version) == 0 {
		return errors.New("version must be specified")
	}

	baseURL := "https://raw.githubusercontent.com/metallb/metallb/" + version

	// Namespace
	url := baseURL + "/manifests/namespace.yaml"
	fileNamespace, err := util.DownloadFile(url)
	if err != nil {
		return err
	}
	yamlFile, err := ioutil.ReadFile(fileNamespace)
	if err != nil {
		return err
	}
	apply.YAML(c, yamlFile)

	// Metallb
	url = baseURL + "/manifests/metallb.yaml"
	fileMetallb, err := util.DownloadFile(url)
	if err != nil {
		return err
	}
	yamlFile, err = ioutil.ReadFile(fileMetallb)
	if err != nil {
		return err
	}
	apply.YAML(c, yamlFile)

	// Done, removing files
	os.Remove(fileNamespace)
	os.Remove(fileMetallb)

	return nil
}

// CreateSecret will create a config as configmap based on the
// struct InstanceConfig
//
// Args:
//
//      - Client struct from client module
//	- InstanceSecret struct
//
// Returns:
//      - nil or error
func CreateSecret(c *client.Client, s *secret.Instance) error {
	// Adding secret for metallb
	s = &secret.Instance{
		Name:      s.Name,
		Namespace: s.Namespace,
		Type:      s.Type,
		Key:       s.Key,
		Value:     s.Value,
	}

	err := secret.Create(c, s)
	if err != nil {
		return err
	}
	return nil
}

// CreateConfig will create a config as configmap based on the
// struct InstanceConfig
//
// Args:
//
//      - Client struct from client module
//	- InstanceConfig
//
// Returns:
//      - nil or error
func CreateConfig(c *client.Client, conf *InstanceConfig) error {
	cfgmap := configmap.Instance{
		Name:      conf.Name,
		Namespace: conf.Namespace,
		ConfigKey: conf.ConfigName,
		ConfigValue: "address-pools:\n- name: " + conf.AddressPoolName +
			"\n  protocol: " + conf.AddressPoolProtocol +
			"\n  addresses:\n  - " + conf.AddressPoolAddresses + "\n",
	}

	err := configmap.Create(c, &cfgmap)
	if err != nil {
		return err
	}
	return nil
}
