package main

import (
	"os"
	"strings"
//	"fmt"
//	"context"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//	"github.com/spf13/viper"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/clientcmd"
//	metav2 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateNamespaceName() string {
	DeployKubeNamespace := os.Getenv("DEPLOY_KUBER_NAMESPACE")
	CIBranchName := os.Getenv("BRANCH_NAME")
	// If environment variables are not set, default names will be used
	if DeployKubeNamespace == "" {
		DeployKubeNamespace = "defaultNamespace"
	}
	if CIBranchName == "" {
		CIBranchName = "defaultBranch"
	}
	return strings.Join([]string{DeployKubeNamespace, CIBranchName}, "-")
}

func createNamespace(ctx *pulumi.Context) (*corev1.Namespace, error) {
	nsName := generateNamespaceName()

	namespace, err := corev1.NewNamespace(ctx, nsName , &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(nsName),
			Annotations: pulumi.StringMap{
				"ct-dynamic-env/ttl": pulumi.String("1h30m"),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return namespace, nil
}