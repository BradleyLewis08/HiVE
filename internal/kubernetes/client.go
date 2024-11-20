package k8sclient

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func GetKubernetesClient() (*Client, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		panic(err);
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)

	if err != nil {
		panic(err);
	}

	return &Client{clientset: clientset}, nil
}

func (c *Client) DeployService(service *apiv1.Service) error {
	_, err := c.clientset.CoreV1().Services(apiv1.NamespaceDefault).Create(context.TODO(), service, metav1.CreateOptions{})
	return err
}

func (c *Client) DeployDeployment(deployment *appsv1.Deployment) error {
	_, err := c.clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func (c* Client) CreateConfigMap(configMap *apiv1.ConfigMap) error {
	_, err := c.clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault).Create(context.TODO(), configMap, metav1.CreateOptions{})
	return err
}

func (c* Client) UpdateConfigMap(configMap *apiv1.ConfigMap) error {
	_, err := c.clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	return err
}

func (c* Client) GetServiceIP(serviceName string) (string, error) {
	service, err := c.clientset.CoreV1().Services((apiv1.NamespaceDefault)).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if service.Spec.Type == apiv1.ServiceTypeLoadBalancer {
		// Wait for LoadBalncer to be assigned IP
		for i := 0; i < 30; i++ {
			service, err = c.clientset.CoreV1().Services(apiv1.NamespaceDefault).Get(context.TODO(), serviceName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}

			if len(service.Status.LoadBalancer.Ingress) > 0 {
				if service.Status.LoadBalancer.Ingress[0].IP != "" {
					return service.Status.LoadBalancer.Ingress[0].IP, nil
				}

				if service.Status.LoadBalancer.Ingress[0].Hostname != "" {
					return service.Status.LoadBalancer.Ingress[0].Hostname, nil
				}
			}
			time.Sleep(time.Second);
		}
	}

	if len(service.Spec.ClusterIP) > 0 {
		return service.Spec.ClusterIP, nil
	}

	return "", fmt.Errorf("service IP not found")
}

func (c* Client) DeleteDeployment(deploymentName string) error {
	err := c.clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	return err
}

func (c* Client) DeleteService(serviceName string) error {
	err := c.clientset.CoreV1().Services(apiv1.NamespaceDefault).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	return err
}

func (c* Client) DeleteConfigMap(configMapName string) error {
	err := c.clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault).Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
	return err
}

func (c* Client) DeploymentExists(deploymentName string) bool {
	_, err := c.clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	return err == nil
}

func (c* Client) DeployIngressController(ingress *networkingv1.Ingress) error {
	_, err := c.clientset.NetworkingV1().Ingresses(apiv1.NamespaceDefault).Create(context.TODO(), ingress, metav1.CreateOptions{})
	return err
}

func (c* Client) GetIngressController() *networkingv1.Ingress {
	ingress, err := c.clientset.NetworkingV1().Ingresses("default").Get(
		context.TODO(),
		"hive-environments",
		metav1.GetOptions{},
	)

	if err != nil {
		fmt.Println("Failed to get ingress controller")
		return nil
	}
	return ingress
}

func(c* Client) UpdateIngressController(newIngress *networkingv1.Ingress) error {
	_, err := c.clientset.NetworkingV1().Ingresses(apiv1.NamespaceDefault).Update(context.TODO(), newIngress, metav1.UpdateOptions{})
	return err
}








