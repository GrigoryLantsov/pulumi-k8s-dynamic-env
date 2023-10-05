package main
//
//import (
//	"os"
//
//	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//)
//
//func KafkaUIInitChart(ctx *pulumi.Context, namespaceName pulumi.StringPtrOutput) (*helmv3.Release, error) {
//	ChartName = "kafka-ui"
//	DeployKubeNamespace = os.Getenv("DEPLOY_KUBER_NAMESPACE")
//	CIBranchName = os.Getenv("BRANCH_NAME")
//	deployImageRepository = "nexus.dev.domoy.ru/docker-remote"
//	deployImageTag = os.Getenv("DEPLOY_IMAGE_TAG")
//	return helmv3.NewRelease(ctx, ChartName, &helmv3.ReleaseArgs{
//		Chart:     pulumi.String(ChartName),
//		Namespace: namespaceName, // Using the generated namespace name
//		Name: pulumi.String(ChartName),
//		ForceUpdate: pulumi.Bool(true),
//		RecreatePods: pulumi.Bool(true),
//		Replace: pulumi.Bool(true),
//		RepositoryOpts: helmv3.RepositoryOptsArgs{
//			Repo: pulumi.String("https://nexus.dev.domoy.ru/repository/helm-local-core"),
//			},
//			Values: pulumi.Map{
//			"global": pulumi.Map{
//				"imageRegistry": pulumi.String(deployImageRepository),
//				},
//				"controller": pulumi.Map{
//				"persistence": pulumi.Map{
//					"enabled": pulumi.Bool(false),
//					},
//					"extraConfig": pulumi.String("offsets.topic.replication.factor=1"),
//					"replicaCount": pulumi.Int(1),
//					},
//					"broker": pulumi.Map{
//				"persistence": pulumi.Map{
//					"enabled": pulumi.Bool(false),
//					},
//					},
//					"zookeeper": pulumi.Map{
//				"persistence": pulumi.Map{
//					"enabled": pulumi.Bool(false),
//					},
//					},
//					"listeners": pulumi.Map{
//				"client": pulumi.Map{
//					"protocol": pulumi.String("PLAINTEXT"),
//					},
//					"controller": pulumi.Map{
//					"protocol": pulumi.String("PLAINTEXT"),
//					},
//					"interbroker": pulumi.Map{
//					"protocol": pulumi.String("PLAINTEXT"),
//					},
//					"external": pulumi.Map{
//					"protocol": pulumi.String("PLAINTEXT"),
//					},
//					},
//					"fullnameOverride": pulumi.String("kafka"),
//					},
//					Description: pulumi.String("Enjoy dynamic environments via pulumi + Helm. Url is " + DeployKubeNamespace + CIBranchName + ".apps.k8s.dev.domoy.ru. Made by Grigory Lantsov."),
//					})
//}
