package provisioner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/BradleyLewis08/HiVE/internal/imager"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/client-go/kubernetes"
)

var HTTPS_PORT = 80
var CODER_PORT = 8080

type Provisioner struct {
	k8sClient *kubernetes.Clientset
	imager *imager.Imager
}

func NewProvisioner(k8sClient *kubernetes.Clientset, imager *imager.Imager) *Provisioner {
	return &Provisioner{k8sClient: k8sClient, imager: imager}
}

func (p* Provisioner) createDeploymentIPService(
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
					Name: "ssh",
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

func (p* Provisioner) getLoadBalancerIP(
	serviceName string,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
			case <-ctx.Done():
				return "", fmt.Errorf("timeout waiting for LoadBalancer IP")
			case <-ticker.C:
				service, err := p.k8sClient.CoreV1().Services(apiv1.NamespaceDefault).Get(context.TODO(), serviceName, metav1.GetOptions{})
				if err != nil {
					fmt.Printf("Error getting service: %s\n", err)
					return "", err
				}

				if service.Spec.Type == "ClusterIP" && service.Name != "kubernetes" {
					ingress := service.Status.LoadBalancer.Ingress
					if len(ingress) > 0 {
						if ingress[0].IP != "" {
							return ingress[0].IP, nil
						} else if ingress[0].Hostname != "" {
							return ingress[0].Hostname, nil
						}
					}
				}

				// Wait for 10 seconds before checking again
				fmt.Printf("Waiting for LoadBalancer IP...\n")
				time.Sleep(10 * time.Second)
		}
	}
}

func (p* Provisioner) createDeployment(
	courseName string,
	imageName string,
	netId string,
) *appsv1.Deployment {
	deploymentName := fmt.Sprintf("hive-course-%s-%s", courseName, netId)
    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: deploymentName,
            Labels: map[string]string{
                "app": "hive-course",
                "course": courseName,
				"student": netId,
            },
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: int32Ptr(1),
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app": "hive-course",
                    "course": courseName,
					"student": netId,
                },
            },
            Template: apiv1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app": "hive-course",
                        "course": courseName,
						"student": netId,
                    },
                },
                Spec: apiv1.PodSpec{
					Containers: []apiv1.Container {
						{
							Name: "code-server",
							Image: imageName,
							Ports: []apiv1.ContainerPort {
								{
									ContainerPort: int32(CODER_PORT),
								},
							},
							Env: []apiv1.EnvVar {
								{
									Name: "PASSWORD",
									Value: "password",
								},
							},
							VolumeMounts: []apiv1.VolumeMount {
								{
									Name: "workspace",
									MountPath: fmt.Sprintf("home/coder/proj/%s", netId),
								},
							},
						},
					},
					Volumes: []apiv1.Volume {
						{
							Name: "workspace",
							VolumeSource: apiv1.VolumeSource {
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
					},
                },
				
            },
        },
    }
	return deployment
}

func (p *Provisioner) createNginxConfigMap(routes map[string]string) (*apiv1.ConfigMap, error) {
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

	var locationBlocks strings.Builder
	for path, service := range routes {
		locationBlocks.WriteString(fmt.Sprintf(`
		location /%s/ {
			proxy_pass http://%s:80/;
		}
		`, path, service))
	}

	configData := fmt.Sprintf(nginxConfig, locationBlocks.String())

	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-config",
		},
		Data: map[string]string{
			"nginx.conf": configData,
		},
	}

	return configMap, nil
}

func (p *Provisioner) deployNginx(configMapName string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-reverse-proxy",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx-reverse-proxy",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx-reverse-proxy",
					},
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

	_, err := p.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// Create LoadBalancer service for NGINX
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-reverse-proxy-service",
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeLoadBalancer,
			Ports: []apiv1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"app": "nginx-reverse-proxy",
			},
		},
	}

	_, err = p.k8sClient.CoreV1().Services(apiv1.NamespaceDefault).Create(context.TODO(), service, metav1.CreateOptions{})
	return err
}


func (p* Provisioner) CreateReverseProxy(
	routes map[string]string,
) error {
	// Create ConfigMap
	configMap, err := p.createNginxConfigMap(routes);
	fmt.Printf("Creating ConfigMap...\n")

	if err != nil {
		fmt.Printf("Error creating ConfigMap: %s\n", err)
		return err
	}

	_, err = p.k8sClient.CoreV1().ConfigMaps(apiv1.NamespaceDefault).Create(context.TODO(), configMap, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf("Error creating ConfigMap: %s\n", err)
		return err
	}

	fmt.Printf("ConfigMap created successfully\n")

	fmt.Printf("Deploying NGINX...\n")
	err = p.deployNginx(configMap.Name)

	if err != nil {
		fmt.Printf("Error deploying NGINX: %s\n", err)
		return err
	}

	// Return the IP of the NGINX service


	fmt.Printf("NGINX deployed successfully\n")
	return nil
}


func (p* Provisioner) CreateEnvironment(
	capacity int,
	courseName string,
	image string,
	netID string,
) error {
	// // Create and push class image
	// imageName, err := p.imager.CreateAndPushImage(
	// 	courseName,
	// 	dockerFile,
	// )

	// if err != nil {
	// 	fmt.Printf("Error creating and pushing image: %s\n", err)
	// 	return "", err
	// }	

	// Create deployment with new image
	deployment := p.createDeployment(courseName, "codercom/code-server:latest", netID)
	
	deploymentsClient := p.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault)
	fmt.Printf("Creating deployment...\n")

	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf("Error creating deployment: %s\n", err)
		return err
	}

	// -- Create ClusterIP service
	service := p.createDeploymentIPService(courseName, netID)
	fmt.Printf("Creating ClusterIP for %s %s...\n", courseName, netID)
	servicesClient := p.k8sClient.CoreV1().Services(apiv1.NamespaceDefault)
	_, err = servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf("Error creating service: %s\n", err)
		return err
	}

	if err != nil {
		fmt.Printf("Error getting LoadBalancer IP: %s\n", err)
		return err
	}

	return nil
}

func int32Ptr(i int32) *int32 { return &i }

