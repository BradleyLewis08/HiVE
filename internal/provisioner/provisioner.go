package provisioner

import (
	"context"

	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type Provisioner struct {
	k8sClient *kubernetes.Clientset
}

func NewProvisioner(k8sClient *kubernetes.Clientset) *Provisioner {
	return &Provisioner{k8sClient: k8sClient}
}

func (p* Provisioner) createSSHService(courseName string) *apiv1.Service {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-ssh", courseName),
			Labels: map[string]string{
				"app": "hive-course",
				"course": courseName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "hive-course",
				"course": courseName,
			},
			Ports: []apiv1.ServicePort{
				{
					Name: "ssh",
					Port: 2222,
					TargetPort: intstr.FromInt(2222),
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}
	return service
}

func (p* Provisioner) createContainers() []apiv1.Container {
	return []apiv1.Container{
		{
			Name: "student-env",
			Image: "ubuntu:latest",
			Command: []string{"/bin/sh"},
			Args:    []string{"-c", "echo 'Container starting'; while true; do echo 'Container still running'; sleep 30; done"},
			Ports: []apiv1.ContainerPort{
				{ContainerPort: 8080},
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					// TODO: Make these configurable
					apiv1.ResourceCPU: resource.MustParse("100m"),
					apiv1.ResourceMemory: resource.MustParse("128Mi"),
				},
			},
			VolumeMounts: []apiv1.VolumeMount{
				{Name: "workspace", MountPath: "/home/student/workspace"},
			},
		}, 
		{
			Name: "sshd",
			Image: "linuxserver/openssh-server:latest",
			Ports: []apiv1.ContainerPort{
				{ContainerPort: 2222},
			},
			Env: []apiv1.EnvVar{
				{Name: "PASSWORD_ACCESS", Value: "false"},
				{Name: "USER_NAME", Value: "student"},
				{Name: "PUBLIC_KEY", Value: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxpi8lfZXTcYLpd5EuN2bTMggarJ0XTrDkhrJ4zmIJkp3ByF8mqHoYZ5FKbyOeMn9z/N/s8FC1m9mD/3QkJu62sIOww+E0LFuO2tu2nPBBDFfEOrJUJ/CXRxaMTJy3+XshGtjEr0gQtILl8jktPi0O0OJ+6eltq5qPw1SIV2FIeDMp8Odl81vSH6E13nq+nceZBDjVlmmL2IBwvnn8YB0gDqRdU38HO9CKJHqtqWuYfewYIFPWsZLsQD458yVBTrfNt++nvgRuC4KIVeRCd/f2AuW72NVmlQZDn9k+aXg1vKWMR3iN8UdWua3GVT+a+8+KjDeS9MN9HKZzmBpHimW1y3m4ixEql9lqXCJ2AeUN6CPJJjlXfS0+17dgQzZ3hEf6ZhKHMDs5qr5Oir6FfQ+j+dfdcNVeRcPnGpNvxQP4DPshUygChLi7zgLFBJ6bMSJipVxPViBh86vYow2G4kiqzps7WQOg/oBs/NIozlPf1NCwYGnW0secM8NvYJCwad/mZnELKDm8b7tmqKoqlBuYlmpDZo5n2JM/F+/Beu7hjN1gL8Mlstk16Ci9ccitP4Evb1r4y6ZkXBzSy3ZFnFuX9WdWtiNelNZ8ndBT7mJ2Tk+/EfNvpdM7tboukfQoy1NTQQaVqWRJmAW1t/4UxVIDredH7ftKOhEz/VtYbrKkzw== bradleyevanlewis@gmail.com"},
			},
			VolumeMounts: []apiv1.VolumeMount{
				{Name: "workspace", MountPath: "/home/student/workspace"},
			},
		},
	}
}

func (p* Provisioner) createDeployment(
	capacity int,
	courseName string,
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
            Replicas: int32Ptr(int32(capacity)),
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
					Containers: p.createContainers(),
					// Containers: p.createContainers
                    Volumes: []apiv1.Volume{
                        {
                            Name: "workspace",
							VolumeSource: apiv1.VolumeSource{
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

func (p* Provisioner) CreateEnvironment(
	capacity int,
	courseName string,
) error {
	// Create deployment
	deployment := p.createDeployment(
		capacity,
		courseName,
	)

	fmt.Printf("Deployment: %v\n", deployment)

	deploymentsClient := p.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault)
	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if(err != nil) {
		fmt.Printf("Error creating deployment: %s\n", err)
		return err
	}

	service := p.createSSHService(courseName)
	fmt.Printf("Service: %v\n", service)
	serviceClient := p.k8sClient.CoreV1().Services(apiv1.NamespaceDefault) 
	_, err = serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})

	if(err != nil) {
		fmt.Printf("Error creating service: %s\n", err)
		return err
	}

	fmt.Printf("Deployment created successfully\n")
	return err
}

func int32Ptr(i int32) *int32 { return &i }

