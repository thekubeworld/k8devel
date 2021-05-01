package logschema

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
	"github.com/sirupsen/logrus"
)

// SetLogrusLogging will set the log format for logrus
//
// Args:
//	- None
//
// Returns:
//	- None
//
func SetLogrusLogging() {
	customFormatter := new(logrus.TextFormatter)
        customFormatter.TimestampFormat = "2006-01-02 15:04:05"
        logrus.SetFormatter(customFormatter)
        customFormatter.FullTimestamp = true
	logrus.Infof("Finished logrus log format settings...")
        logrus.Infof("\n")
}
