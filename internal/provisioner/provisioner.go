package provisioner

import (
	"context"
	"fmt"
	"time"

	"github.com/BradleyLewis08/HiVE/internal/imager"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/client-go/kubernetes"
)

type Provisioner struct {
	k8sClient *kubernetes.Clientset
	imager *imager.Imager
}

func NewProvisioner(k8sClient *kubernetes.Clientset, imager *imager.Imager) *Provisioner {
	return &Provisioner{k8sClient: k8sClient, imager: imager}
}

func (p* Provisioner) createLoadBalancerService(
	courseName string, 
) *apiv1.Service {
	service := &apiv1.Service {
		ObjectMeta: metav1.ObjectMeta {
			Name: fmt.Sprintf("%s-lb", courseName),
		},
		Spec: apiv1.ServiceSpec {
			Type: apiv1.ServiceTypeLoadBalancer,
			Ports: []apiv1.ServicePort{
				{
					Name: "ssh",
					Port: 22,
					TargetPort: intstr.FromInt(22),
				},
			},
			Selector: map[string]string {
				"app": "hive-course",
				"course": courseName,
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

				if service.Spec.Type == "LoadBalancer" {
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
) *appsv1.Deployment {
	deploymentName := fmt.Sprintf("hive-course-%s", courseName)
    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: deploymentName,
            Labels: map[string]string{
                "app": "hive-course",
                "course": courseName,
            },
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: int32Ptr(1),
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app": "hive-course",
                    "course": courseName,
                },
            },
            Template: apiv1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app": "hive-course",
                        "course": courseName,
                    },
                },
                Spec: apiv1.PodSpec{
					Containers: []apiv1.Container {
						{
							Name: "student-env",
							Image: imageName,
							Ports: []apiv1.ContainerPort {
								{
									ContainerPort: 22,
								},
							},
							Command: []string {
								"/bin/bash",
							},
							Args: []string{
								"-c",
								"while true; do echo hello; sleep 10; done",
							},
						},
					},
                },
            },
        },
    }
	return deployment
}


func (p* Provisioner) CreateEnvironment(
	capacity int,
	courseName string,
	dockerFile string,
) (string, error) {

	// Create and push class image
	// imageName, err := p.imager.CreateAndPushImage(
	// 	courseName,
	// 	dockerFile,
	// )

	// if err != nil {
	// 	fmt.Printf("Error creating and pushing image: %s\n", err)
	// 	return "", err
	// }	

	// Create deployment with new image
	deployment := p.createDeployment(courseName, "bradleylewis08/course-environments:test")
	
	deploymentsClient := p.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault)
	fmt.Printf("Creating deployment...\n")
	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf("Error creating deployment: %s\n", err)
		return "", err
	}

	// Create service
	service := p.createLoadBalancerService(courseName)
	fmt.Printf("Creating LoadBalancer service...\n")
	servicesClient := p.k8sClient.CoreV1().Services(apiv1.NamespaceDefault)
	_, err = servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})

	if err != nil {
		fmt.Printf("Error creating service: %s\n", err)
		return "", err
	}

	// Get LoadBalancer IP
	fmt.Printf("Getting LoadBalancer IP...\n")
	loadBalancerIP, err := p.getLoadBalancerIP(service.Name)

	if err != nil {
		fmt.Printf("Error getting LoadBalancer IP: %s\n", err)
		return "", err
	}

	return loadBalancerIP, nil
}

func int32Ptr(i int32) *int32 { return &i }

