package podsecuritypolicy

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

	"github.com/thekubeworld/k8devel/pkg/client"
)

// ListAllPodSecurityPolicy will list all PodSecurityPolicoies
//
// Args:
//      - Client struct from client module
//
// Return:
//      - error or nil
func ListAllPodSecurityPolicy(c *client.Client) error {
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
		fmt.Println("\n")
		fmt.Println("Name:", psp.Items[i].Name)
		fmt.Println("\tPrivileged:", psp.Items[i].Spec.Privileged)

		if psp.Items[i].Spec.AllowPrivilegeEscalation != nil {
			fmt.Println("\tAllow PrivilegeEscalation:", *psp.Items[i].Spec.AllowPrivilegeEscalation)
		} else {
			fmt.Println("\tAllow PrivilegeEscalation: unset")
		}

		fmt.Println("\tRequired DropCapabilities:", psp.Items[i].Spec.RequiredDropCapabilities)
		fmt.Println("\tDefault AddCapabilities:", psp.Items[i].Spec.DefaultAddCapabilities)
		fmt.Println("\tRequired DropCapabilities:", psp.Items[i].Spec.RequiredDropCapabilities)
		fmt.Println("\tAllowed Capabilities:", psp.Items[i].Spec.AllowedCapabilities)
		fmt.Println("\tVolumes:", psp.Items[i].Spec.Volumes)

		if len(psp.Items[i].Spec.AllowedFlexVolumes) > 0 {
			fmt.Println("\tAllowed FlexVolumes:", psp.Items[i].Spec.AllowedFlexVolumes)
		}
		fmt.Println("\tAllowed CSIDrivers:", psp.Items[i].Spec.AllowedCSIDrivers)
		fmt.Println("\tAllowed UnsafeSysctls:", psp.Items[i].Spec.AllowedUnsafeSysctls)
		fmt.Println("\tForbidden Sysctls:", psp.Items[i].Spec.ForbiddenSysctls)
		fmt.Println("\tAllow Host Network:", psp.Items[i].Spec.HostNetwork)
		fmt.Println("\tAllow Host Ports:", psp.Items[i].Spec.HostPorts)
		fmt.Println("\tAllow Host PID:", psp.Items[i].Spec.HostPID)
		fmt.Println("\tAllow Host IPC:", psp.Items[i].Spec.HostIPC)
		fmt.Println("\tRead Only Root Filesystem:", psp.Items[i].Spec.ReadOnlyRootFilesystem)
		fmt.Println("\tSELinux Context Strategy:", psp.Items[i].Spec.SELinux.Rule)
                if psp.Items[i].Spec.SELinux.SELinuxOptions != nil {
			fmt.Println("\t\tUser:", psp.Items[i].Spec.SELinux.SELinuxOptions.User)
			fmt.Println("\t\tRole:", psp.Items[i].Spec.SELinux.SELinuxOptions.Role)
			fmt.Println("\t\tSELinux Type:", psp.Items[i].Spec.SELinux.SELinuxOptions.Type)
			fmt.Println("\t\tLevel:", psp.Items[i].Spec.SELinux.SELinuxOptions.Level)
                }
		fmt.Println("\tRun As User Strategy:", psp.Items[i].Spec.RunAsUser.Rule)
		fmt.Println("\tRanges: ", psp.Items[i].Spec.RunAsUser.Ranges)

		fmt.Println("\tFSGroup Strategy:", psp.Items[i].Spec.FSGroup.Rule)
		fmt.Println("\tRanges:", psp.Items[i].Spec.FSGroup.Ranges)

		fmt.Println("\tSupplemental Groups Strategy:", psp.Items[i].Spec.SupplementalGroups.Rule)
		fmt.Println("\tRanges:", idRangeToString(psp.Items[i].Spec.SupplementalGroups.Ranges))
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
