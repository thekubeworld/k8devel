package k8devel

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
	"context"
	"fmt"
	"strings"
	"errors"

	v1beta1 "k8s.io/api/policy/v1beta1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sirupsen/logrus"
)

// ListAllPodSecurityPolicy will list all PodSecurityPolicoies
//
// Args:
//      - Client struct from client module
//
// Return:
//      - error or nil
func ListAllPodSecurityPolicy(c *Client) error {
	psp, err := c.Clientset.PolicyV1beta1().PodSecurityPolicies().List(
		context.TODO(),
		metav1.ListOptions{})
	if err != nil {
		return errors.New("Error listing PodSecurityPolicies")
	}

	if psp == nil || len(psp.Items) == 0 {
		return errors.New("no PodSecurityPolicies found; assuming PodSecurityPolicy is disabled")
	}
	for i := 0; i < len(psp.Items); i++ {
		logrus.Infof("\n")
		logrus.Infof("Name: %s", psp.Items[i].Name)
		logrus.Infof("\tPrivileged: %v", psp.Items[i].Spec.Privileged)

		if psp.Items[i].Spec.AllowPrivilegeEscalation != nil {
			logrus.Infof("\tAllow PrivilegeEscalation: %v", *psp.Items[i].Spec.AllowPrivilegeEscalation)
		} else {
			logrus.Infof("\tAllow PrivilegeEscalation: unset")
		}

		logrus.Infof("\tRequired DropCapabilities: %v", psp.Items[i].Spec.RequiredDropCapabilities)
		logrus.Infof("\tDefault AddCapabilities: %v", psp.Items[i].Spec.DefaultAddCapabilities)
		logrus.Infof("\tRequired DropCapabilities: %v", psp.Items[i].Spec.RequiredDropCapabilities)
		logrus.Infof("\tAllowed Capabilities: %v", psp.Items[i].Spec.AllowedCapabilities)
		logrus.Infof("\tVolumes: %v", psp.Items[i].Spec.Volumes)

		if len(psp.Items[i].Spec.AllowedFlexVolumes) > 0 {
			logrus.Infof("\tAllowed FlexVolumes: %v", psp.Items[i].Spec.AllowedFlexVolumes)
		}
		logrus.Infof("\tAllowed CSIDrivers: %v", psp.Items[i].Spec.AllowedCSIDrivers)
		logrus.Infof("\tAllowed UnsafeSysctls: %v", psp.Items[i].Spec.AllowedUnsafeSysctls)
		logrus.Infof("\tForbidden Sysctls: %v", psp.Items[i].Spec.ForbiddenSysctls)
		logrus.Infof("\tAllow Host Network: %v", psp.Items[i].Spec.HostNetwork)
		logrus.Infof("\tAllow Host Ports: %v", psp.Items[i].Spec.HostPorts)
		logrus.Infof("\tAllow Host PID: %v", psp.Items[i].Spec.HostPID)
		logrus.Infof("\tAllow Host IPC: %v", psp.Items[i].Spec.HostIPC)
		logrus.Infof("\tRead Only Root Filesystem: %v", psp.Items[i].Spec.ReadOnlyRootFilesystem)
		logrus.Infof("\tSELinux Context Strategy: %v", psp.Items[i].Spec.SELinux.Rule)
                if psp.Items[i].Spec.SELinux.SELinuxOptions != nil {
			logrus.Infof("\t\tUser: %s", psp.Items[i].Spec.SELinux.SELinuxOptions.User)
			logrus.Infof("\t\tRole: %s", psp.Items[i].Spec.SELinux.SELinuxOptions.Role)
			logrus.Infof("\t\tSELinux Type: %s", psp.Items[i].Spec.SELinux.SELinuxOptions.Type)
			logrus.Infof("\t\tLevel: %s", psp.Items[i].Spec.SELinux.SELinuxOptions.Level)
                }
		logrus.Infof("\tRun As User Strategy: %s", psp.Items[i].Spec.RunAsUser.Rule)
		logrus.Infof("\tRanges: %s", psp.Items[i].Spec.RunAsUser.Ranges)

		logrus.Infof("\tFSGroup Strategy: %s", psp.Items[i].Spec.FSGroup.Rule)
		logrus.Infof("\tRanges: %s", psp.Items[i].Spec.FSGroup.Ranges)

		logrus.Infof("\tSupplemental Groups Strategy: %s:", psp.Items[i].Spec.SupplementalGroups.Rule)
		logrus.Infof("\tRanges: %s", idRangeToString(psp.Items[i].Spec.SupplementalGroups.Ranges))
        }
	return nil
}

// idRangeToString will return string from idRange
// Params:
// 	[]v1beta1.IDRange
//
// Returns
//	string
func idRangeToString(ranges []v1beta1.IDRange) string {
        formattedString := ""
        if ranges != nil {
                strRanges := []string{}
                for _, r := range ranges {
                        strRanges = append(strRanges, fmt.Sprintf("%d-%d", r.Min, r.Max))
                }
                formattedString = strings.Join(strRanges, ",")
        }
        return formattedString
}
