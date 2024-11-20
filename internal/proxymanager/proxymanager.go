package proxymanager

import (
	"fmt"
	"sync"

	"github.com/BradleyLewis08/HiVE/deployments"
	k8sclient "github.com/BradleyLewis08/HiVE/internal/kubernetes"
	"github.com/BradleyLewis08/HiVE/internal/utils"
	"github.com/BradleyLewis08/HiVE/services"
)

type ProxyManager struct {
	k8sClient *k8sclient.Client
	mu sync.RWMutex
	routes map[string]string // location -> proxyPass
	proxyIPAddress string
}

func NewProxyManager(k8sClient *k8sclient.Client) *ProxyManager {
	return &ProxyManager{k8sClient: k8sClient, routes: make(map[string]string)}
}

func (pm *ProxyManager) DeleteExistingRouter() {
	pm.k8sClient.DeleteDeployment("master-router")
	pm.k8sClient.DeleteService("master-router")
	pm.k8sClient.DeleteConfigMap("master-router")
}

func (pm *ProxyManager) ProvisionMasterRouter() error {
	configMap := deployments.DefaultNginxConfigMap()
	err := pm.k8sClient.CreateConfigMap(configMap)

	if err != nil {
		fmt.Println("Failed to create config map")
		return err
	}

	nginxDeployment := deployments.NewNginxDeployment(configMap.Name);
	err = pm.k8sClient.DeployDeployment(nginxDeployment)

	if err != nil {
		fmt.Println("Failed to deploy master router")
		return err
	}

	nginxService := services.NewNginxService()
	err = pm.k8sClient.DeployService(nginxService)

	if err != nil {
		fmt.Println("Failed to deploy nginx router service")
		return err
	}

	serviceAddr, err := pm.k8sClient.GetServiceIP(nginxService.Name);

	if err != nil {
		fmt.Println("Failed to get service IP")
		return err
	}

	pm.proxyIPAddress = serviceAddr
	return nil
}

func (pm *ProxyManager) AddRoute(assignmentName string, courseName string, netID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	relativeRoute := utils.ConstructRelativeRoute(assignmentName, courseName, netID)
	pm.routes[relativeRoute] = utils.ConstructLoadBalancerRoute(assignmentName, courseName, netID)
	return pm.updateNginxConfig()
}

func (pm *ProxyManager) RemoveRoute(assignmentName string, courseName string, netID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	relativeRoute := utils.ConstructRelativeRoute(assignmentName, courseName, netID)
	delete(pm.routes, relativeRoute)
	return pm.updateNginxConfig()
}

func (pm* ProxyManager) updateNginxConfig() error {
	configMap := deployments.NewNginxConfigMap(pm.routes);
	return pm.k8sClient.UpdateConfigMap(configMap)
}

func (pm* ProxyManager) GetProxyIPAddress() string {
	return pm.proxyIPAddress
}

