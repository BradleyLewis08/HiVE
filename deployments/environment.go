package deployments

import (
	"fmt"

	"github.com/BradleyLewis08/HiVE/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CODER_PORT = 8080

func NewEnvironmentDeployment(courseName string, imageName string, netId string) *appsv1.Deployment {
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
            Replicas: utils.Int32ptr(1),
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