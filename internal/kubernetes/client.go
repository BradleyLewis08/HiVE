package k8sclient

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetKubernetesClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	clientset, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		panic(err);
	}

	return kubernetes.NewForConfig(clientset)
}

