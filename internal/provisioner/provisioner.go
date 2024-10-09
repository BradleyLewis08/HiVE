package provisioner

import (
	"fmt"

	"github.com/BradleyLewis08/HiVE/deployments"
	k8sclient "github.com/BradleyLewis08/HiVE/internal/kubernetes"
	"github.com/BradleyLewis08/HiVE/internal/utils"
	"github.com/BradleyLewis08/HiVE/services"
)

var HTTPS_PORT = 80
var CODER_PORT = 8080

type Provisioner struct {
	k8sClient *k8sclient.Client
}

func NewProvisioner(k8sClient *k8sclient.Client) *Provisioner {
	return &Provisioner{k8sClient: k8sClient}
}

func (p* Provisioner) ProvisionStudentEnvironment(
	capacity int,
	courseName string,
	image string,
	netID string,
) error {
	// TODO: Make this custom
	image = "codercom/code-server:latest"
	environmentDeployment := deployments.NewEnvironmentDeployment(courseName, image, netID)
	fmt.Printf("Creating deployment for %s %s...\n", courseName, netID)
	err := p.k8sClient.DeployDeployment(environmentDeployment)

	if err != nil {
		fmt.Printf("Error creating deployment: %s\n", err)
		return err
	}

	// -- Create ClusterIP service
	fmt.Printf("Creating ClusterIP for %s %s...\n", courseName, netID)
	service := services.NewEnvironmentIPService(courseName, netID)
	err = p.k8sClient.DeployService(service)

	if err != nil {
		fmt.Printf("Error creating service: %s\n", err)
		return err
	}
	return nil
}

func (p* Provisioner) ProvisionCourseRouter(
	courseName string,
	netIDs []string,
) error {
	routes := utils.ConstructReverseProxyRoutes(netIDs, courseName)
	configMap := deployments.NewNginxConfigMap(courseName, routes)

	err := p.k8sClient.CreateConfigMap(configMap)

	if err != nil {
		fmt.Println("Failed to create config map")
		return err
	}

	nginxDeployment := deployments.NewNginxDeployment(courseName, configMap.Name)
	err = p.k8sClient.DeployDeployment(nginxDeployment)

	if err != nil {
		fmt.Println("Failed to deploy nginx deployment")
		return err
	}

	nginxService := services.NewNginxService(courseName)

	err = p.k8sClient.DeployService(nginxService)

	if err != nil {
		fmt.Println("Failed to deploy nginx service")
		return err
	}

	fmt.Printf("Successfully deployed course router for %s\n", courseName)
	return nil
}


