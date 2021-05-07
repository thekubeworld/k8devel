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
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// DetectContainerPortProtocol is a helper for users
// use only TCP or UDP string instead of require them
// to manage k8s.io/api/core/v1. This will convert
// TCP or UDP to v1.ProtocolTCP or v1.ProtocolUDP
//
// Args:
// 	protocol - TCP or UDP as string
//
// Returns:
//	v1.ProtocolUDP, v1.ProtocolTCP or error
func DetectContainerPortProtocol(protocol string) (v1.Protocol, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return v1.ProtocolTCP, nil
	case "udp":
		return v1.ProtocolUDP, nil
	}
	return "", errors.New("unknown protocol")
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
	logrus.Infof("Diffing %s and %s", fileone, filetwo)
	path, err := exec.LookPath("diff")
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	logrus.Infof("%s -r -u -N %s %s", path, fileone, filetwo)
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
