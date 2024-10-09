package utils

import (
	"fmt"
)

func ConstructReverseProxyRoutes(netIDs []string, courseName string) map[string]string {
	routes := make(map[string]string)

	for _, netID := range netIDs {
		proxy_pass := fmt.Sprintf("%s-%s-lb.default.svc.cluster.local", courseName, netID)
		routes[netID] = proxy_pass
	}

	return routes
}

func Int32ptr(i int32) *int32 { return &i }