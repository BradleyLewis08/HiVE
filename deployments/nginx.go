package deployments

import (
	"fmt"
	"strings"

	"github.com/BradleyLewis08/HiVE/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Make this dynamic
var LOAD_BALANCER_SUFFIX_URL = "-lb.default.svc.cluster.local"

func constructLocationBlocks(routes map[string]string) string {
	var locationBlocks strings.Builder

	for path, service := range routes {
		locationBlocks.WriteString(fmt.Sprintf(`
		location /%s/ {
			proxy_pass http://%s:80/;
		}
		`, path, service))
	}

	return locationBlocks.String()
}

func NewNginxConfigMap(courseName string, routes map[string]string) (*apiv1.ConfigMap) {
	configMapName := fmt.Sprintf("nginx-config-%s", courseName)
	nginxConfig := `
	events {}
	http {
		map $http_upgrade $connection_upgrade {
			default upgrade;
			'' close;
		}
		server {
			listen 80;
			
			# Proxy settings
			proxy_set_header Host $host;
			proxy_set_header X-Real-IP $remote_addr;
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
			proxy_set_header X-Forwarded-Proto $scheme;
			proxy_set_header Upgrade $http_upgrade;
			proxy_set_header Connection "upgrade";
			
			# WebSocket support
			proxy_http_version 1.1;
			
			# Increase max body size if needed
			client_max_body_size 10m;

			%s
		}
	}`

	locationBlocks := constructLocationBlocks(routes)

	configData := fmt.Sprintf(nginxConfig, locationBlocks)

	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configMapName,
		},
		Data: map[string]string{
			"nginx.conf": configData,
		},
	}

	return configMap
}

func NewNginxDeployment(courseName string, configMapName string) *appsv1.Deployment { 
	deploymentName := fmt.Sprintf("nginx-reverse-proxy-%s", courseName)

	labels := map[string]string {
		"app": "nginx-reverse-proxy",
		"course": courseName,
		"hive-component": "reverse-proxy",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "nginx-config",
									MountPath: "/etc/nginx/nginx.conf",
									SubPath:   "nginx.conf",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "nginx-config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: configMapName,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}
