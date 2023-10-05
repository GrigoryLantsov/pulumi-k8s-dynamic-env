package main

import (
	"os"

	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func PostgresqlInitChart(ctx *pulumi.Context, namespaceName pulumi.StringPtrOutput) (*helmv3.Release, error) {
	ChartName = "postgresql"
	DeployKubeNamespace = os.Getenv("DEPLOY_KUBER_NAMESPACE")
	CIBranchName = os.Getenv("BRANCH_NAME")
	deployImageRepository = ""
	deployImageTag = os.Getenv("DEPLOY_IMAGE_TAG")
	return helmv3.NewRelease(ctx, ChartName, &helmv3.ReleaseArgs{
		Chart:     pulumi.String(ChartName),
		Namespace: namespaceName, // Using the generated namespace name
		Name: pulumi.String(ChartName),
		ForceUpdate: pulumi.Bool(true),
		RecreatePods: pulumi.Bool(true),
		Replace: pulumi.Bool(true),
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://grigorylantsov.github.io/helm-library"),
		},
		Values: pulumi.Map{
			"image": pulumi.Map{
				"registry": pulumi.String(deployImageRepository),
			},
			"primary": pulumi.Map{
				"persistence": pulumi.Map{
					"enabled": pulumi.Bool(false),
				},
			},
			"fullnameOverride": pulumi.String("postgresql"),
			"global": pulumi.Map{
				"postgresql": pulumi.Map{
					"auth": pulumi.Map{
						"postgresPassword": pulumi.String("tmp_password_132!"),
					},
				},
			},
		},
		Description: pulumi.String("Enjoy dynamic environments via pulumi + Helm. Url is " + DeployKubeNamespace + CIBranchName + domainName + "Made by Grigory Lantsov."),
	})
}
