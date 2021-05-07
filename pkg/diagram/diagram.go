package diagram

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
	"fmt"
)

// ClusterIP show a diagram about ClusterIP Service
//
// Args:
//	None
//
// Returns:
//	None
//
func ClusterIP() {
	fmt.Println("           POD                  ")
	fmt.Println("            |                   ")
	fmt.Println("         Traffic                ")
	fmt.Println("            |                   ")
	fmt.Println("            |                   ")
	fmt.Println("   +-------------------+        ")
	fmt.Println("   |     kube-proxy    |        ")
	fmt.Println("   +-------------------+        ")
	fmt.Println("       |           |            ")
	fmt.Println("+------------------------------+")
	fmt.Println("|     Service (Cluster IP)     |")
	fmt.Println("|+-----------------------------+")
	fmt.Println("|       |           |          |")
	fmt.Println("|   +-------+   +-------+      |")
	fmt.Println("|   |  Pod  |   |  Pod  |      |")
	fmt.Println("|   +-------+   +-------+      |")
	fmt.Println("|                              |")
	fmt.Println("| kubernetes cluster           |")
	fmt.Println("+------------------------------+")
}

// NodePort show a diagram about NodePort Service
//
// Args:
//	None
//
// Returns:
//	None
//
func NodePort() {
	fmt.Println("               +---------------------+              ")
	fmt.Println("               |      TRAFFIC        |              ")
	fmt.Println("               |    FROM USERS       |              ")
	fmt.Println("               +---------------------+              ")
	fmt.Println("                          |                         ")
	fmt.Println("                          v                         ")
	fmt.Println("             +-- kube-proxy/iptables-+              ")
	fmt.Println("             |                       |              ")
	fmt.Println("             v                       v              ")
	fmt.Println("          Node IPs               Node IPs           ")
	fmt.Println("         and Ports               and ports          ")
	fmt.Println("+------------|-------+    +-----------|------------+")
	fmt.Println("| 10.10.50.54|:30001 |    |10.10.50.51|:30001      |")
	fmt.Println("+------------|--------    +-----------|------------+")
	fmt.Println("|            +------------+-----------+            |")
	fmt.Println("|                         |                        |")
	fmt.Println("+-----------------------+-V-+----------------------|")
	fmt.Println("| ClusterIP 10.111.239.7|:80|                      |")
	fmt.Println("+-----------------------|-|-|----------------------|")
	fmt.Println("| Service 1             | | |                      |")
	fmt.Println("|   Selector:           | | |                      |")
	fmt.Println("|   app: nginx          | | |                      |")
	fmt.Println("+-----------------------|-|-|----------------------|")
	fmt.Println("|               +-------|-+-|-----------+          |")
	fmt.Println("|               |       |   |           |          |")
	fmt.Println("+---------------|-------|   |-----------|----------|")
	fmt.Println("| EndpointIP and| ports |   |EndpointIP | and ports|")
	fmt.Println("|               v       |   |           v          |")
	fmt.Println("|   10.244.2.2 :80 :8080|   |10.244.2.3 :80  :8080 |")
	fmt.Println("+---------------|---|---|   |-----------|----|-----|")
	fmt.Println("|               |   |   |   |           |    |     |")
	fmt.Println("+---------------|---|---|   |-----------|----|-----|")
	fmt.Println("| Container port|   |   |   | Container |Port|     |")
	fmt.Println("|               |   |   |   |           |    |     |")
	fmt.Println("|             :80 :8080 |   |          :80  :8080  |")
	fmt.Println("+---------------|---|---|---|-----------|----|     |")
	fmt.Println("|               |   |   |   |           |    |     |")
	fmt.Println("|               |   |   |   |           |    |     |")
	fmt.Println("|               |   |   |   |           |    |     |")
	fmt.Println("|               v   |   |   |           v    |     |")
	fmt.Println("|      Container 1  |   |   |  Container 1   |     |")
	fmt.Println("|                   |   |   |                |     |")
	fmt.Println("|                   v   |   |                v     |")
	fmt.Println("|         Container 2   |   |          Container 2 |")
	fmt.Println("|    Labels: app nginx  |   |    Labels: app nginx |")
	fmt.Println("|                       |   |                      |")
	fmt.Println("| Pod 1                 |   | Pod 2                |")
	fmt.Println("| Node 1                |   | Node 2               |")
	fmt.Println("+-----------------------+   +----------------------+")
}

// LoadBalancer show a diagram about LoadBalancer Service
//
// Args:
//	None
//
// Returns:
//	None
//
func LoadBalancer() {
	fmt.Println("               +---------------------+              ")
	fmt.Println("               |      TRAFFIC        |              ")
	fmt.Println("               |    FROM USERS       |              ")
	fmt.Println("               +---------------------+              ")
	fmt.Println("                          |                         ")
	fmt.Println("                          v                         ")
	fmt.Println("             +-----------------------+              ")
	fmt.Println("             |    Load Balancer      |              ")
	fmt.Println("             +-----------------------+              ")
	fmt.Println("  +----------|-----------------------|-----------+  ")
	fmt.Println("  |          v                       v           |  ")
	fmt.Println("  |       +-------------------------------+      |  ")
	fmt.Println("  |       |            Service            |      |  ")
	fmt.Println("  |       +-------------------------------+      |  ")
	fmt.Println("  |           |            |           |         |  ")
	fmt.Println("  |           v            v           v         |  ")
	fmt.Println("  |      +-------+     +------+    +------+      |  ")
	fmt.Println("  |      |  POD  |     |  POD |    |  POD |      |  ")
	fmt.Println("  |      +-------+     +------+    +------+      |  ")
	fmt.Println("  |                                              |  ")
	fmt.Println("  | Kubernetes Cluster                           |  ")
	fmt.Println("  +----------------------------------------------+  ")
}
