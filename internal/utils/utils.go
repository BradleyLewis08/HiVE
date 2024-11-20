package utils

import (
	"fmt"
	"strings"
)

func LowerCaseAndStrip(courseName string) string {
	lowerCase := strings.ToLower(courseName)
	transformed := strings.ReplaceAll(lowerCase, " ", "-")
	return transformed
}

func ConstructRelativeRoute(assignmentName string, courseName string, netID string) string {
	return fmt.Sprintf("%s/%s/%s", assignmentName, courseName, netID)
}

func ConstructEnvironmentDeploymentName(assignmentName string, courseName string, netID string) string {
	return fmt.Sprintf("hive-environment-%s-%s-%s", assignmentName, courseName, netID)
}

func ConstructLoadBalancerServiceName(assignmentName string, courseName string, netID string) string {
	return fmt.Sprintf("%s-%s-%s-lb", assignmentName, courseName, netID)
}

func ConstructLoadBalancerRoute(assignmentName string, courseName string, netID string) string {
	return fmt.Sprintf("%s-%s-%s-lb.default.svc.cluster.local", assignmentName, courseName, netID)
}

func Int32ptr(i int32) *int32 { return &i }