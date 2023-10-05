package main

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	pulumi.RunErr(func(ctx *pulumi.Context) error {
		ns, err := createNamespace(ctx) // Create namespace if it doesn't exist
		if err != nil {
			return err
		}

		namespaceName := ns.Metadata.Name() // Execute the function to get the actual namespace name
		ctx.Export("Namespace", namespaceName)
		CreateServiceAccountStack(ctx, namespaceName) // Create Service Account

		CreateNewRoleBindings(ctx, namespaceName)

		postgreSQLEnabled := os.Getenv("POSTGRESQL_ENABLED")

		if postgreSQLEnabled == "true" {

			_, err := PostgresqlInitChart(ctx, namespaceName) // Create POSTGRESQL
			if err != nil {
				// Handle error
			}
		} else {
			fmt.Println("PostgreSQL initialization is not enabled.")
		}

		kafkaEnabled := os.Getenv("KAFKA_ENABLED")

		if kafkaEnabled == "true" {

			_, err := KafkaInitChart(ctx, namespaceName) // Create Kafka
			if err != nil {
				// Handle error
			}
		} else {
			fmt.Println("kafka initialization is not enabled.")
		}

		camundaEnabled := os.Getenv("CAMUNDA_ENABLED")

		if camundaEnabled == "true" {

			_, err := CamundaInitChart(ctx, namespaceName) // Create Kafka
			if err != nil {
				// Handle error
			}
		} else {
			fmt.Println("Camunda initialization is not enabled.")
		}
		extraName := os.Getenv("EXTRA_NAME")
		_, description, err := deployChart(ctx, namespaceName, DeployKubeNamespace + extraName) // Deploy the Helm chart using the namespace name
		if err != nil {
			return err
		}

		extraChart := os.Getenv("EXTRA_DEPLOY")

		if extraChart == "true" {
			extraName := os.Getenv("EXTRA_NAME_2")
			_, _, err := deployChart(ctx, namespaceName, DeployKubeNamespace + extraName)
			// Create Additional Helm Chart
			if err != nil {
				return err
			}
		} else {
			fmt.Println("ExtraChart initialization is not enabled.")
		}

		// Attempt to create the Ingress, but don't raise exceptions if it fails.
		_, ingressErr := CreateIngress(ctx, namespaceName, DeployKubeNamespace + "-" + CIBranchName + extraName)

		if ingressErr != nil {
			// Log the error but continue program execution.
			ctx.Log.Warn("Failed to create Ingress", nil)
		} else {
			// Check if the Ingress was created successfully and log it.
			ctx.Log.Info("Ingress created", nil)
		}

		ctx.Export("Made By:", description)
		return nil // Return nil if all operations are successful
	})
}
