package main

import (
	"os"

	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployChart(ctx *pulumi.Context, namespaceName pulumi.StringPtrOutput, releaseName string) (*helmv3.Release, pulumi.StringOutput, error) {
	ChartName = os.Getenv("CHART_NAME")
	DeployKubeNamespace = os.Getenv("DEPLOY_KUBER_NAMESPACE")
	CIBranchName = os.Getenv("BRANCH_NAME")
	deployImageRepository = os.Getenv("DEPLOY_IMAGE_REPOSITORY") + os.Getenv("EXTRA_NAME")
	deployImageTag = os.Getenv("DEPLOY_IMAGE_TAG")
	domainName 	   = os.Getenv("UPPER_DOMAIN_NAME")
	hostIngress := pulumi.String("https://" + DeployKubeNamespace + "-" + CIBranchName + domainName)
	description := pulumi.Sprintf("Enjoy dynamic environments via pulumi + Helm. Url is %s-%s%s. Made by Grigory Lantsov.", DeployKubeNamespace, CIBranchName, domainName)

	release, err := helmv3.NewRelease(ctx, releaseName, &helmv3.ReleaseArgs{
		Chart:     pulumi.String(ChartName),
		Name:	   pulumi.String(releaseName),
		Namespace: namespaceName, // Using the generated namespace name
		ForceUpdate: pulumi.Bool(true),
		RecreatePods: pulumi.Bool(true),
		Replace: pulumi.Bool(true),
		ValueYamlFiles: pulumi.AssetOrArchiveArray{
			pulumi.NewFileAsset(".helm/values.yaml"),
		},
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://grigorylantsov.github.io/helm-library"),
		},
		Values: pulumi.Map{
			"global": pulumi.Map{
				"env": pulumi.String("pulumi"),
			},
			"fullnameOverride": pulumi.String(releaseName + "-" + ChartName),
			"ingress": pulumi.Map{
				"nginx": pulumi.Map{
					"enabled": pulumi.Bool(false),
				},
			},
			"image": pulumi.Map{
				"repository": pulumi.String(deployImageRepository),
				"tag":        pulumi.String(deployImageTag),
			},
			"secrets": pulumi.Map{
				"rolename": pulumi.String(DeployKubeNamespace + "-dynamic-role"),
			},
			"application": pulumi.Map{
				"env": pulumi.Map{
					"BACKEND_URL": pulumi.Map{
						"pulumi": pulumi.String(hostIngress),
					},
				},
			},
		},
		Description: pulumi.StringOutput(description),
	})

	return release, description, err
}
