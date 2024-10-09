package k8sclient

import (
	"context"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
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






