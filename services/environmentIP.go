package services

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var HTTPS_PORT = 80
var CODER_PORT = 8080

func NewEnvironmentIPService(
	courseName string, 
	netId string,
) *apiv1.Service {
	service := &apiv1.Service {
		ObjectMeta: metav1.ObjectMeta {
			Name: fmt.Sprintf("%s-%s-lb", courseName, netId),
		},
		Spec: apiv1.ServiceSpec {
			Type: apiv1.ServiceTypeClusterIP,
			Ports: []apiv1.ServicePort{
				{
					Name: "environmentip",
					Port: int32(HTTPS_PORT),
					TargetPort: intstr.FromInt(CODER_PORT),
				},
			},
			Selector: map[string]string {
				"app": "hive-course",
				"course": courseName,
				"student": netId,
			},
		},
	}
	return service
}