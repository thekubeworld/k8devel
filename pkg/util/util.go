package util

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
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// DetectReclaimPolicy is a helper for users to specify the reclaim policy
// as string: "retain", "delete", "recycle"
//
// Args:
//	string - retain, delete and recycle
//
// Returns:
//	v1.LimitType (v1.LimitTypePod, v1.LimitTypeContainer, v1.LimitTypePersistentVolumeClaim,  or error)
//
func DetectReclaimPolicy(reclaimPolicy string) (v1.PersistentVolumeReclaimPolicy, error) {
	switch strings.ToLower(reclaimPolicy) {
	case "retain":
		return v1.PersistentVolumeReclaimRetain, nil
	case "delete":
		return v1.PersistentVolumeReclaimDelete, nil
	case "recycle":
		return v1.PersistentVolumeReclaimRecycle, nil
	}
	return "", errors.New("unknown reclaim policy type")
}

// DetectConcurrencyPolicy is a helper for users
// use only Allow, Forbid and Replace instead of require
// them to manage k8s.io/api/batch/v1/types.go
//
// Args:
// 	allow, forbid, replace
//
// Returns:
//	batch.AllowConcurrent, batch.ForbidConcurrent, batch.ReplaceConcurrent
func DetectConcurrencyPolicy(currencyPolicy string) (batchv1.ConcurrencyPolicy, error) {
	switch strings.ToLower(currencyPolicy) {
	case "allow":
		return batchv1.AllowConcurrent, nil
	case "forbid":
		return batchv1.ForbidConcurrent, nil
	case "replace":
		return batchv1.ReplaceConcurrent, nil
	}
	return "", errors.New("unknown concurrency policy")
}

// LimitType is a helper for users to specify the limit type
// as string: "Pod", "Container" or "PersistentVolumeClaim" in a namespace
//
// Args:
//	string - pod, container or persistentvolumeclaim
//
// Returns:
//	v1.LimitType (v1.LimitTypePod, v1.LimitTypeContainer, v1.LimitTypePersistentVolumeClaim,  or error)
//
func DetectLimitType(limitType string) (v1.LimitType, error) {

	switch strings.ToLower(limitType) {
	case "pod":
		return v1.LimitTypePod, nil
	case "container":
		return v1.LimitTypeContainer, nil
	case "persistentvolumeclaim":
		return v1.LimitTypePersistentVolumeClaim, nil
	}
	return "", errors.New("unknown limit type")
}

// DetectVolumeAccessModes is a helper for users
// use access mode as: rox, rwo, rwx, readonlymany,
// readwritemany, readwriteonce
//
// Remember:
//
// ROX - ReadOnlyMany - can be mounted in read-only mode to many hosts
// RWO - ReadWriteOnce - can be mounted in read/write mode to exactly 1 host
// RWX - ReadWriteMany - can be mounted in read/write mode to many hosts
//
// Args:
//	string - rox, rwo, rwx, readonlymany,
//		 readwritemany, readwriteonce
//
// Returns:
//	v1.PersistentVolumeAccessMode (v1.ReadOnlyMany, v1.ReadWriteMany, v1.ReadWriteOnce) or error
func DetectVolumeAccessModes(access string) ([]v1.PersistentVolumeAccessMode, error) {

	switch strings.ToLower(access) {
	case "readonlymany", "rox":
		return []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany}, nil
	case "readwritemany", "rwx":
		return []v1.PersistentVolumeAccessMode{v1.ReadWriteMany}, nil
	case "readwriteonce", "rwo":
		return []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}, nil
	}
	return nil, errors.New("unknown access mode")
}

// DetectVolumeMode is a helper for users
// use mode as:
//
// persistentvolumeblock or block
// persistentvolumefilesystem or filesystem
//
// Args:
//	string - persistentvolumeblock or block
//               persistentvolumefilesystem or filesystem
//
// Returns:
//	v1.PersistentVolumeMode or error
func DetectVolumeMode(mode string) (v1.PersistentVolumeMode, error) {

	switch strings.ToLower(mode) {
	case "persistentvolumeblock", "block":
		return v1.PersistentVolumeBlock, nil
	case "persistentvolumefilesystem", "filesystem":
		return v1.PersistentVolumeFilesystem, nil
	}
	return "", errors.New("unknown volume mode")
}

// DetectContainerPortProtocol is a helper for users
// to use TCP or UDP words instead of require them
// to manage k8s.io/api/core/v1. This will convert
// TCP or UDP to v1.ProtocolTCP or v1.ProtocolUDP
//
// Args:
//     protocol - tcp or udp as string
//
// Returns:
//     v1.ProtocolUDP, v1.ProtocolTCP or error
func DetectContainerPortProtocol(protocol string) (v1.Protocol, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return v1.ProtocolTCP, nil
	case "udp":
		return v1.ProtocolUDP, nil
	}
	return "", errors.New("unknown protocol")
}

// DetectImagePullPolicy is a helper for users to use more
// friendly words like: always, ifnotpresent or never instead
// to require them to manage k8s.io/api/core/v1.
//
// Args:
// 	ifnotpresent, always or never
//
// Returns:
//	v1.PullAlways
//	v1.PullNever
//	v1.PullIfNotPresent or error
func DetectImagePullPolicy(pullpolicy string) (v1.PullPolicy, error) {
	switch strings.ToLower(pullpolicy) {
	case "always":
		return v1.PullAlways, nil
	case "never":
		return v1.PullNever, nil
	case "ifnotpresent":
		return v1.PullIfNotPresent, nil
	}
	return "", errors.New("unknown image pull policy")
}

// DetectContainerRestartPolicy is a helper for users to more
// friendly words like: onfailure, never or always instead of
// require them to manage k8s.io/api/core/v1.
//
// Args:
// 	onfailure, always, never
//
// Returns:
//	v1.RestartPolicyOnFailure
//	v1.RestartPolicyNever
//	v1.RestartPolicyAlways or error
func DetectContainerRestartPolicy(policy string) (v1.RestartPolicy, error) {
	switch strings.ToLower(policy) {
	case "onfailure":
		return v1.RestartPolicyOnFailure, nil
	case "never":
		return v1.RestartPolicyNever, nil
	case "always":
		return v1.RestartPolicyAlways, nil
	}
	return "", errors.New("unknown restart policy")
}

// CompareFiles will compare two files, byte by byte
// to see if they are equal
//
// Args:
//    fileone - first file to compare
//    filetwo - second file to compare
//
//   Returns:
//      bool and error
func CompareFiles(fileone string, filetwo string) (bool, error) {
	f1, err := ioutil.ReadFile(fileone)

	if err != nil {
		return false, err
	}

	f2, err := ioutil.ReadFile(filetwo)
	if err != nil {
		return true, err
	}

	return bytes.Equal(f1, f2), nil
}

// DiffCommand will diff two files
//
// Args:
//    fileone
//    filtwo
//
//   Returns:
//      bytes from the file or error
func DiffCommand(fileone string, filetwo string) ([]byte, error) {
	fmt.Printf("Diffing %s and %s\n", fileone, filetwo)
	path, err := exec.LookPath("diff")
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	fmt.Printf("%s -r -u -N %s %s\n", path, fileone, filetwo)
	cmd := exec.Command(path,
		"-r",
		"-u",
		"-N",
		fileone,
		filetwo)
	cmd.Stdout = &out

	// WARN: cannot use err = cmd.Run()
	// diff command if find a difference will return -1
	// and err will be empty
	cmd.Run()

	return out.Bytes(), nil
}

// CreateTempFile will create temporary file
//
// Args:
//    dirname - dir name
//    filename - file name
//
//   Returns:
//      filename as string or error
func CreateTempFile(dirname string, filename string) (*os.File, error) {
	file, err := ioutil.TempFile(dirname, filename+".*")
	if err != nil {
		return nil, err
	}
	// DO NOT ADD defer here, it returns the file pointers to the caller
	return file, nil
}

// DownloadFile will download a file specified as temporary file
//
// Args:
//    url - url to be download
//
//   Returns:
//      path as string or error
func DownloadFile(url string) (string, error) {
	out, err := CreateTempFile(os.TempDir(), "downloadedfile")
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

// GenerateRandomString will generate a random string
//
// Args:
//	number of chars to be generated
//	modeString lower or all
//
//   Returns:
//	random string
func GenerateRandomString(numberOfChars int, modeString string) (string, error) {
	rand.Seed(time.Now().UnixNano())

	var letters []rune
	if modeString == "lower" {
		letters = []rune("abcdefghijklmnopqrstuvwxyz")
	} else if modeString == "all" {
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	} else {
		return "", errors.New("modeString must be lower or all")
	}

	result := make([]rune, numberOfChars)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result), nil
}

// GetResourceList returns a ResourceList with the
// specified cpu and memory or disk resource values
func GetResourceList(cpu string, memory string, ephemeralStorage string) v1.ResourceList {
	res := v1.ResourceList{}
	if cpu != "" {
		res[v1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		res[v1.ResourceMemory] = resource.MustParse(memory)
	}
	if ephemeralStorage != "" {
		res[v1.ResourceEphemeralStorage] = resource.MustParse(ephemeralStorage)
	}
	return res
}
