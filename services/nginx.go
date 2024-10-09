package services

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewNginxService(courseName string) *apiv1.Service {
	serviceName := fmt.Sprintf("nginx-reverse-proxy-%s", courseName)

	labels := map[string]string{
		"app":            "nginx-reverse-proxy",
		"course":         courseName,
		"hive-component": "reverse-proxy",
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: labels,
		},
		Spec: apiv1.ServiceSpec{
			Selector: labels,
			Ports: []apiv1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}
	return service
}