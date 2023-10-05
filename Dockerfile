FROM dtzar/helm-kubectl:3.13.0 as builder

FROM pulumi/pulumi-go:3.86.0

WORKDIR /app

COPY main k8s-dynamic-env
COPY Pulumi.yaml .
COPY Pulumi.dev.yaml .
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm
COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl

RUN apt remove curl git -y && \
    rm -f /usr/local/go && \
    apt-get update && \
    apt-get upgrade -y && \
    apt autoremove -y

CMD bash
