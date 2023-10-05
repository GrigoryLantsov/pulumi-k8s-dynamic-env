package main

import (
	"os"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateServiceAccountStack is a Pulumi stack for creating a Kubernetes ServiceAccount.
func CreateServiceAccountStack(ctx *pulumi.Context, namespaceName pulumi.StringPtrOutput) error {
	DeployKubeNamespace = os.Getenv("DEPLOY_KUBER_NAMESPACE")
	serviceAccount, err := corev1.NewServiceAccount(ctx, "pulumi_dynamic_service_account", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(DeployKubeNamespace + "-deployer"),
			Namespace: namespaceName,
		},
	})
	if err != nil {
		return err
	}

	// Export the name of the ServiceAccount for reference in other stacks.
	ctx.Export("serviceAccountName", serviceAccount.Metadata.Name())

	return nil
}
