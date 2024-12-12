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

const NGINX_NAME = "master-router"
const NGINX_BASE_CONFIG = `
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
	}
`

func constructLocationBlocks(routes map[string]string) string {
	var locationBlocks strings.Builder

	for path, service := range routes {
		locationBlocks.WriteString(fmt.Sprintf(`
		location /%s/ {
			proxy_pass http://%s:8080/;
			proxy_set_header X-Original-URI $request_uri;
			proxy_set_header Accept-Encoding "";
		}
		`, path, service))
	}

	return locationBlocks.String()
}

func DefaultNginxConfigMap() *apiv1.ConfigMap {
	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: NGINX_NAME,
		},
		Data: map[string]string{
			"nginx.conf": fmt.Sprintf(NGINX_BASE_CONFIG, ""),
		},
	}

	return configMap
} 

func NewNginxConfigMap(routes map[string]string) (*apiv1.ConfigMap) {
	locationBlocks := constructLocationBlocks(routes)
	configData := fmt.Sprintf(NGINX_BASE_CONFIG, locationBlocks)

	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: NGINX_NAME,
		},
		Data: map[string]string{
			"nginx.conf": configData,
		},
	}

	return configMap
}

func NewNginxDeployment(configMapName string) *appsv1.Deployment { 
	labels := map[string]string {
		"app": "nginx-reverse-proxy",
		"hive-component": "reverse-proxy",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: NGINX_NAME,
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
