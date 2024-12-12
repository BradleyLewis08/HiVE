package deployments

import (
	"fmt"

	"github.com/BradleyLewis08/HiVE/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CODER_PORT = 8080

func NewEnvironmentDeployment(assignmentName string, courseName string, imageName string, netId string) *appsv1.Deployment {
	deploymentName := utils.ConstructEnvironmentDeploymentName(assignmentName, courseName, netId)
	labels := map[string]string{
		"app": "hive-course",
		"course": courseName,
		"assignment": assignmentName,
		"student": netId,
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
					Containers: []apiv1.Container {
						{
							Name: "code-server",
							Image: imageName,
							Ports: []apiv1.ContainerPort {
								{
									Name: "http",
									ContainerPort: 8080,
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