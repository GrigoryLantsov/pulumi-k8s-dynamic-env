package main

import (
	"fmt"
	"io/ioutil"
	"os"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v2"
)

func CreateIngress(ctx *pulumi.Context, namespaceName pulumi.StringPtrOutput, name string) (*networkingv1.Ingress, error) {
	DeployKubeNamespace := os.Getenv("DEPLOY_KUBER_NAMESPACE")
	extraName := os.Getenv("EXTRA_NAME")
	ChartName := os.Getenv("CHART_NAME")
	CIBranchName := os.Getenv("BRANCH_NAME")
	domainName 	 := os.Getenv("UPPER_DOMAIN_NAME")
	hostIngress := pulumi.String(DeployKubeNamespace + "-" + CIBranchName + domainName)

	valuesFile, err := ioutil.ReadFile(".helm/values.yaml")
	if err != nil {
		return nil, err
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(valuesFile, &values)
	if err != nil {
		return nil, err
	}

	// Extract the annotations from the parsed values
	ingressValues, ok := values["ingress"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("ingress values not found or not in the expected format")
	}

	annotationsValue, ok := ingressValues["annotations"].(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("annotations not found or not in the expected format")
	}

	annotations := map[string]string{}

	for key, value := range annotationsValue {
		// Convert key and value to strings
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("annotation key is not a string: %v", key)
		}

		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("annotation value is not a string: %v", value)
		}

		annotations[keyStr] = valueStr
	}

	ingressAnnotations := pulumi.StringMap{}

	for key, value := range annotations {
		// Convert key to a string explicitly
		keyStr := fmt.Sprintf("%v", key)
		ingressAnnotations[keyStr] = pulumi.String(value)
	}

	// Extract service port
	servicePort := int(values["service"].(map[interface{}]interface{})["ports"].(map[interface{}]interface{})["http"].(int))

	// Extract the "hosts" section from the "ingress" section
	ingressHosts, ok := ingressValues["hosts"]
	if !ok {
		return nil, fmt.Errorf("hosts config not found or not in the expected format")
	}

	// Check if the "hosts" section is a slice of interfaces
	hostsConfig, ok := ingressHosts.([]interface{})
	if !ok {
		return nil, fmt.Errorf("hosts is not a slice of interfaces")
	}

	var ingressRules networkingv1.IngressRuleArray

	for _, hostValue := range hostsConfig {
		hostConfig, ok := hostValue.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("host config not found or not in the expected format")
		}

		pathsConfig, ok := hostConfig["paths"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("paths config not found or not in the expected format")
		}

		var ingressPaths networkingv1.HTTPIngressPathArray

		for _, pathValue := range pathsConfig {
			pathConfig, ok := pathValue.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("path config not found or not in the expected format")
			}

			path := pathConfig["path"].(string)
			pathType := pathConfig["pathType"].(string)

			ingressPath := &networkingv1.HTTPIngressPathArgs{
				Backend: &networkingv1.IngressBackendArgs{
					Service: &networkingv1.IngressServiceBackendArgs{
						Name: pulumi.String(DeployKubeNamespace + extraName + "-" + ChartName),
						Port: &networkingv1.ServiceBackendPortArgs{
							Number: pulumi.Int(servicePort),
						},
					},
				},
				Path:     pulumi.String(path),
				PathType: pulumi.String(pathType),
			}
			ingressPaths = append(ingressPaths, ingressPath)
		}

		if len(ingressPaths) > 0 {
			ingressRule := networkingv1.IngressRuleArgs{
				Host: pulumi.String(hostIngress),
				Http: &networkingv1.HTTPIngressRuleValueArgs{
					Paths: ingressPaths,
				},
			}
			ingressRules = append(ingressRules, ingressRule)
		}
	}

	if len(ingressRules) == 0 {
		return nil, fmt.Errorf("no valid ingress rules found")
	}

	// Create Ingress resource
	ingress, err := networkingv1.NewIngress(ctx, name, &networkingv1.IngressArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace:   namespaceName,
			Name:        pulumi.StringPtr(name),
			Annotations: ingressAnnotations,
		},
		Spec: &networkingv1.IngressSpecArgs{
			Rules: ingressRules,
		},
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "15s", Update: "15s"}))
	if err != nil {
		return ingress, err
	}

	return ingress, nil
}
