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

// Provisions pod and ClusterIP service for student environment
func (p* Provisioner) ProvisionStudentEnvironment(
	assignmentName string,
	courseName string,
	image string,
	netID string,
) error {
	environmentDeployment := deployments.NewEnvironmentDeployment(assignmentName, courseName, image, netID)
	fmt.Printf("Creating deployment for %s %s...\n", courseName, netID)
	err := p.k8sClient.DeployDeployment(environmentDeployment)

	if err != nil {
		fmt.Printf("Error creating deployment: %s\n", err)
		return err
	}

	// -- Create ClusterIP service
	fmt.Printf("Creating ClusterIP for %s:%s %s...\n", courseName, assignmentName, netID)
	service := services.NewEnvironmentService(assignmentName, courseName, netID)
	err = p.k8sClient.DeployService(service)

	if err != nil {
		fmt.Printf("Error creating service: %s\n", err)
		return err
	}

	return nil
}

func (p* Provisioner) DeleteEnvironment(assignmentName string, courseName string, netID string) error {
	deploymentName := utils.ConstructEnvironmentDeploymentName(assignmentName, courseName, netID)
	// Initialize error variable
	err := p.k8sClient.DeleteDeployment(deploymentName)
	if err != nil {
		fmt.Printf("Failed to delete deployment %s\n", deploymentName)
	}

	// Delete ClusterIP service
	serviceName := utils.ConstructLoadBalancerServiceName(assignmentName, courseName, netID)
	err = p.k8sClient.DeleteService(serviceName)

	if err != nil {
		fmt.Printf("Failed to delete service %s\n", serviceName)
	}

	fmt.Printf("Successfully deleted environment for %s %s\n", courseName, netID)
	return err
}

