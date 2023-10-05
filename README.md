# k8s-dynamic-env-pulumi

## Overview

This image allows you to create a dynamic environment using [Pulumi](https://www.pulumi.com/). Currently, it supports the following:

1. Namespace creation. Depends on variables ```DEPLOY_KUBER_NAMESPACE```-```CI_BRANCH_NAME```.
2. Creation of a service account for Vault access. Created with the pattern ```DEPLOY_KUBER_NAMESPACE```-deployer. Role configuration in Vault is required beforehand.
3. Creation of a role model for domain access via IPA with role bindings.
4. Helm chart creation using one of the [unified Helm charts](https://github.com/GrigoryLantsov/helm-library).
5. Ability to run two Helm charts simultaneously using the variables ```EXTRA_DEPLOY (bool)```, ```EXTRA_NAME_2```
6. Creation and launch of Kafka in a stateless format using [KRaft](https://developer.confluent.io/learn/kraft/) and the [Bitnami Kafka chart](https://github.com/bitnami/charts/tree/main/bitnami/kafka). Switching is controlled by the variable ```KAFKA_ENABLED.```
7. Creation and launch of PostgreSQL 16 using the [Bitnami PostgreSQL chart](https://github.com/bitnami/charts/blob/main/bitnami/postgresql). Switching is controlled by the variable ```POSTGRESQL_ENABLED```.
8. Creation and launch of Camunda using the [unified Helm chart](https://github.com/GrigoryLantsov/helm-library).
9. Separate ingress creation. It involves parsing the ```.helm/values.yaml``` file, substituting the necessary parameters, and deploying separately. *Note: Due to the fact that .status.LoadBalancer is checked after ingress creation, which is mainly relevant for cloud technologies, an error may occur.*


## Environment Variables

| Variable                | Description                                                       | Example                                                  |
|-------------------------|-------------------------------------------------------------------|----------------------------------------------------------|
| CHART_NAME              | Chart Name                                                        | node-app                                                 |
| DEPLOY_KUBER_NAMESPACE  | core prefix on namespace                                          | core                                                     |
| BRANCH_NAME             | Variable, mb predefined in CI Controller for current Branch       | dev                                                      |
| DEPLOY_IMAGE_REPOSITORY | Registry for docker container                                     | ```${DOCKER_REGISTRY}/${CI_PROJECT_NAME}${EXTRA_NAME}``` |
| DEPLOY_IMAGE_TAG        | Image TAG                                                         | ```${CI_COMMIT_REF_SLUG}.${CI_COMMIT_SHORT_SHA}```       |
| EXTRA_NAME              | Actual for MonoRep. Just an EXTRA_NAME for service                | service                                                  |
| EXTRA_DEPLOY            | Actual for MonoRep. If we want deploy 2 services at the same time | true                                                     |
| EXTRA_NAME_2            | Second name for service in Monorep                                | nonservice                                               |
| POSTGRESQL_ENABLED      | Deploy Stateless pg sql                                           | true                                                     |
| KAFKA_ENABLED           | Deploy stateless kafka                                            | true                                                     |
| CAMUNDA_ENABLED         | Deploy stateless camunda                                          | true                                                     |

## Environment Variables for ci/cd

| Key                                      | Value    |
|------------------------------------------|----------|
| PULUMI_CONFIG_PASSPHRASE                 | ""       |
| PULUMI_AUTOMATION_API_SKIP_VERSION_CHECK | true     |
| PULUMI_SKIP_UPDATE_CHECK                 | true     |

## How to Use

**Reference for gitlab. How to use that**

```sh
# --------------------------------------------------------------------------------------
.pulumi.deploy:
  stage: deploy
  variables:
    KAFKA_ENABLED: "true"
    POSTGRESQL_ENABLED: "true"
    PULUMI_CONFIG_PASSPHRASE: ""
    PULUMI_AUTOMATION_API_SKIP_VERSION_CHECK: "true"
    PULUMI_SKIP_UPDATE_CHECK: "true"
    BRANCH_NAME: ${CI_COMMIT_REF_SLUG}
  image:
    name: kaligrir/k8s-dynamic-env:stable
    pull_policy: always
  before_script:
    - echo ${KUBECONFIG_DEV} | base64 -d > ${PWD}/kubeconfig
    - export KUBECONFIG=${PWD}/kubeconfig
  script:
    - mv /app/* .
    - |
      if [[ -n "${EXTRA_NAME}" ]]; then
        mv .helm/values${EXTRA_NAME}.yaml .helm/values.yaml
      else
        echo "Values file doesn't changed"
      fi
    - ls -lah
    - echo ${BRANCH_NAME}
    - pulumi logout
    - pulumi login --local
    - pulumi stack ls -a
    - |
      if pulumi stack ls -a | grep "organization/k8s-dynamic-env/gl"; then
        pulumi rm organization/k8s-dynamic-env/gl -y -f
      else
        echo "Pulumi stack 'organization/k8s-dynamic-env/gl' does not exist. Skipping..."
      fi
    - pulumi stack init organization/k8s-dynamic-env/gl
    - pulumi preview --stack organization/k8s-dynamic-env/gl --logtostderr --logflow -v=9 2> out.txt
    - pulumi up --yes --non-interactive --skip-preview --stack organization/k8s-dynamic-env/gl || true
    # we use it without state, so we must delete secret for effective redeploy in the same branch
    - kubectl delete secret $(kubectl get secret -n ${DEPLOY_KUBER_NAMESPACE}-${BRANCH_NAME} | grep sh.helm | awk '{ print $1 }') -n ${DEPLOY_KUBER_NAMESPACE}-${BRANCH_NAME}
  rules:
    - if: '$CI_COMMIT_REF_NAME == "prod"'
      when: never
    - when: manual
  environment:
    name: pulumi
# --------------------------------------------------------------------------------------
```    

```sh
deploy_dynamic_backend_service:
  extends: .pulumi.deploy
  needs:
    - build_image
  before_script:
    - export EXTRA_NAME="-service"
    - export CAMUNDA_ENABLED="true"
    - export POSTGRESQL_ENABLED="true"
    - export EXTRA_DEPLOY="true"
    - export EXTRA_NAME_2="-engine"
    - echo ${KUBECONFIG_DEV} | base64 -d > ${PWD}/kubeconfig
    - export KUBECONFIG=${PWD}/kubeconfig
```

### Extra Links

For clean pods by time i use and a little bit refactor for my personal needs that repo [TwiN/k8s-ttl-controller](https://github.com/TwiN/k8s-ttl-controller). Big Thanks to him.