FROM golang:1.12.7 as builder

COPY . /cron
WORKDIR /cron
RUN make modules
RUN make build

FROM alpine:3.9
COPY --from=builder /cron /chaoscoordinatorservice
WORKDIR /chaoscoordinatorservice