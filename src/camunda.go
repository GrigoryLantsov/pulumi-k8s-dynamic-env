package main

import (
	"os"

	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CamundaInitChart(ctx *pulumi.Context, namespaceName pulumi.StringPtrOutput) (*helmv3.Release, error) {
	ChartName = "camunda-bpm-platform"
	DeployKubeNamespace = os.Getenv("DEPLOY_KUBER_NAMESPACE")
	CIBranchName = os.Getenv("BRANCH_NAME")
	deployImageRepository = ""
	deployImageTag = os.Getenv("DEPLOY_IMAGE_TAG")
	hostIngress := pulumi.String(DeployKubeNamespace + "-" + CIBranchName + ".apps.k8s.dev.domoy.ru")
	applicationConf := `
camunda.bpm:
  auto-deployment-enabled: false
  generic-properties:
    properties:
      enableExceptionsAfterUnhandledBpmnError: true
  admin-user:
    id: ${CAMUNDA_ADMIN_USER}
    password: ${CAMUNDA_ADMIN_PASSWORD}
    firstName: Admin
    lastName: Admin
  filter:
    create: All tasks
  run:
    cors:
      enabled: false
      allowed-origins: "http://localhost:9000"
      exposed-headers: "Access-Control-Allow-Origin"
      allow-credentials: true
    example:
      enabled: false
spring.web.resources:
  static-locations: NULL`
	return helmv3.NewRelease(ctx, ChartName, &helmv3.ReleaseArgs{
		Chart:     pulumi.String(ChartName),
		Namespace: namespaceName, // Using the generated namespace name
		Name: pulumi.String(ChartName),
		ForceUpdate: pulumi.Bool(true),
		RecreatePods: pulumi.Bool(true),
		Replace: pulumi.Bool(true),
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://nexus.dev.domoy.ru/repository/helm-local-core"),
			},
			Values: pulumi.Map{
			"global": pulumi.Map{
				"env": pulumi.String("pulumi"),
			},
			"replicaCount": pulumi.Int(1),
			"serviceAccount": pulumi.Map{
				"create": pulumi.Bool(false),
			},
			"metrics": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"autoscaling": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"startupProbe": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"readinessProbe": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"livenessProbe": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"service": pulumi.Map{
				"type": pulumi.String("ClusterIP"),
				"port": pulumi.Int(8080),
				"portName": pulumi.String("http"),
				"protocol": pulumi.String("TCP"),
				"enabled": pulumi.Bool(true),
			},
			"general": pulumi.Map{
				"debug":         pulumi.Bool(false),
				"replicaCount":  pulumi.Int(1),
				"nameOverride":  pulumi.String(""),
				"fullnameOverride": pulumi.String(DeployKubeNamespace + extraName + "-" + ChartName),
			},
			"resources": pulumi.Map{
				"limits": pulumi.Map{
					"cpu": pulumi.Map{
						"_default": pulumi.String("1000m"),
						},
						"memory": pulumi.Map{
						"_default": pulumi.String("1024Mi"),
						},
				},
				"requests": pulumi.Map{
					"cpu": pulumi.Map{
						"_default": pulumi.String("200m"),
					},
					"memory": pulumi.Map{
						"_default": pulumi.String("512Mi"),
					},
				},
			},
			"volumes": pulumi.Map{
				"enabled": pulumi.Bool(true),
				"path":    pulumi.String("/camunda/configuration/default.yml"),
				"subpath": pulumi.String("default.yml"),
			},
			"secrets": pulumi.Map{
				"rolename": pulumi.String(DeployKubeNamespace + "-dynamic-role"),
				"sa":       pulumi.String(DeployKubeNamespace + "-deployer"),
				"enabled":  pulumi.Bool(true),
				"path": pulumi.Array{
					pulumi.Map{
						"name": pulumi.String("DB_PASSWORD"),
						"path": pulumi.Map{
							"_default": pulumi.String("rpp/data/dnmcCreds"),
						},
						"key": pulumi.String("DB_PASSWORD"),
					},
					pulumi.Map{
						"name": pulumi.String("DB_USERNAME"),
						"path": pulumi.Map{
							"_default": pulumi.String("rpp/data/dnmcCreds"),
						},
						"key": pulumi.String("DB_USER"),
					},
					pulumi.Map{
						"name": pulumi.String("CAMUNDA_ADMIN_USER"),
						"path": pulumi.Map{
							"_default": pulumi.String("rpp/data/creds"),
						},
						"key": pulumi.String("CAMUNDA_ADMIN_USER"),
					},
					pulumi.Map{
						"name": pulumi.String("CAMUNDA_ADMIN_PASSWORD"),
						"path": pulumi.Map{
							"_default": pulumi.String("rpp/data/creds"),
						},
						"key": pulumi.String("CAMUNDA_ADMIN_PASSWORD"),
					},
				},
			},
			"ingress": pulumi.Map{
				"nginx": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
				"annotations": pulumi.Map{
					"kubernetes.io/ingress.class": pulumi.String("nginx"),
				},
				"hosts": pulumi.Array{
					pulumi.Map{
						"host": pulumi.Map{
							"pulumi": pulumi.String(hostIngress),
						},
						"paths": pulumi.Array{
							pulumi.Map{
								"path":     pulumi.String("/camunda"),
								"pathType": pulumi.String("Prefix"),
							},
							pulumi.Map{
								"path":     pulumi.String("/engine-rest"),
								"pathType": pulumi.String("Prefix"),
							},
							pulumi.Map{
								"path":     pulumi.String("/swaggerui"),
								"pathType": pulumi.String("Prefix"),
							},
						},
					},
				},
			},
			"image": pulumi.Map{
				"repository": pulumi.String("nexus.dev.domoy.ru/docker-core/camunda"),
				"tag":        pulumi.String("stable"),
				},
				"database": pulumi.Map{
				"driver": pulumi.String("org.postgresql.Driver"),
				"url": pulumi.Map{
					"pulumi": pulumi.String("jdbc:postgresql://postgresql:5432/postgres?reWriteBatchedInserts=true"),
					},
				},
				"application": pulumi.Map{
				    "conf": pulumi.Map{
						"default.yml": pulumi.String(applicationConf),
					},
				},
			},
	})
}
