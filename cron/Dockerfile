FROM golang:1.12.7 as builder

COPY . /cron
COPY genesys-config-k8s-genesys /root/.kube/genesys-config-k8s.genesys
WORKDIR /cron
RUN make modules
RUN make build

FROM alpine:3.9
COPY --from=builder /cron /chaoscoordinatorservice
WORKDIR /chaoscoordinatorservice