package ingress

import (
	"fmt"
	"log"

	k8sclient "github.com/BradleyLewis08/HiVE/internal/kubernetes"
	"github.com/BradleyLewis08/HiVE/internal/utils"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
    SERVICE_PORT = 80
)

type IngressManager struct {
    k8sClient *k8sclient.Client
}

type IngressRule struct {
	Path string
	ServiceName string
	ServicePort int32
}

func NewIngressManager(k8sClient *k8sclient.Client) *IngressManager {
    return &IngressManager{k8sClient: k8sClient}
}

func (im *IngressManager) ProvisionIngressController() {
    rules := []IngressRule {
        {
            Path: fmt.Sprintf("/ping"),
            ServiceName: "ping-service",
            ServicePort: SERVICE_PORT,   
        },
    }

    controller := NewEnvironmentIngressController(rules)
    err := im.k8sClient.DeployIngressController(controller)
    if err != nil {
        log.Fatalf("Failed to deploy ingress controller: %v", err)
    }
}

func (im *IngressManager) AddRouteToIngress(assignmentName string, courseName string, netID string) error {
    ingress := im.k8sClient.GetIngressController()
    if ingress == nil {
        return fmt.Errorf("Ingress controller not found")
    }

    serviceName := utils.ConstructLoadBalancerServiceName(assignmentName, courseName, netID)
    path := fmt.Sprintf("/environment/%s/%s/%s", courseName, assignmentName, netID)
    pathType := networkingv1.PathTypePrefix

    newPath := networkingv1.HTTPIngressPath{
        Path: path,
        PathType: &pathType,
        Backend: networkingv1.IngressBackend{
            Service: &networkingv1.IngressServiceBackend{
                Name: serviceName, 
                Port: networkingv1.ServiceBackendPort{
                    Number: SERVICE_PORT,
                },
            },
        },
    }

    ingress.Spec.Rules[0].HTTP.Paths = append(ingress.Spec.Rules[0].HTTP.Paths, newPath)

    err := im.k8sClient.UpdateIngressController(ingress)

    if err != nil {
        log.Fatalf("Failed to update ingress controller: %v", err)
    }

    return nil
}

// NewEnvironmentIngress creates an Ingress resource for routing to student environments
func NewEnvironmentIngressController(rules []IngressRule) *networkingv1.Ingress {
    pathType := networkingv1.PathTypePrefix
    
    var ingressPaths []networkingv1.HTTPIngressPath
    for _, rule := range rules {
        ingressPaths = append(ingressPaths, networkingv1.HTTPIngressPath{
            Path: rule.Path,
            PathType: &pathType,
            Backend: networkingv1.IngressBackend{
                Service: &networkingv1.IngressServiceBackend{
                    Name: rule.ServiceName,
                    Port: networkingv1.ServiceBackendPort{
                        Number: rule.ServicePort,
                    },
                },
            },
        })
    }

    return &networkingv1.Ingress{
        ObjectMeta: metav1.ObjectMeta{
            Name: "hive-environments",
            Annotations: map[string]string{
                "nginx.ingress.kubernetes.io/rewrite-target": "/$2",
                "nginx.ingress.kubernetes.io/proxy-read-timeout": "3600",
                "nginx.ingress.kubernetes.io/proxy-send-timeout": "3600",
                "nginx.ingress.kubernetes.io/proxy-body-size": "10m",
                // Enable websocket support
                "nginx.ingress.kubernetes.io/proxy-http-version": "1.1",
                "nginx.ingress.kubernetes.io/proxy-buffering": "off",
            },
        },
        Spec: networkingv1.IngressSpec{
            Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue {
							Paths: ingressPaths,
						},
					},
				},
            },
        },
    }
}

