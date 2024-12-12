package services

import (
	"github.com/BradleyLewis08/HiVE/internal/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	CODER_PORT = 8080
	SERVICE_PORT = 80
)

func NewEnvironmentService(
	assignmentName string,
	courseName string, 
	netId string,
) *apiv1.Service {
	labels := map[string]string {
		"app": "hive-course",
		"course": courseName,
		"student": netId,
		"assignment": assignmentName,
	}
	service := &apiv1.Service {
		ObjectMeta: metav1.ObjectMeta {
			Name: utils.ConstructLoadBalancerServiceName(assignmentName, courseName, netId),
		},
		Spec: apiv1.ServiceSpec {
			Type: apiv1.ServiceTypeClusterIP,
			Ports: []apiv1.ServicePort{
				{
					Name: "environmentip",
					Port: SERVICE_PORT,
					TargetPort: intstr.FromInt(CODER_PORT),
				},
			},
			Selector: labels,
		},
	}
	return service
}