package provisioner

import (
	"context"

	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Provisioner struct {
	k8sClient *kubernetes.Clientset
}

func NewProvisioner(k8sClient *kubernetes.Clientset) *Provisioner {
	return &Provisioner{k8sClient: k8sClient}
}

func (p* Provisioner) createDeployment() *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hive-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "hive",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "hive",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "hive",
							Image: "nginx:1.12",
						},
					},
				},
			},
		},
	}
	return deployment;
}

func (p* Provisioner) CreateEnvironment() error {
	deployment := p.createDeployment()
	deploymentsClient := p.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault)
	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if(err != nil) {
		return err
	}
	fmt.Printf("Deployment created successfully\n")
	return err
}

func int32Ptr(i int32) *int32 { return &i }

