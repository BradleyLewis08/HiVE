package services

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewNginxService() *apiv1.Service {
	labels := map[string]string{
		"app":            "nginx-reverse-proxy",
		"hive-component": "reverse-proxy",
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "master-router",
			Labels: labels,
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeLoadBalancer,
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