FROM golang:latest as builder

COPY . /chaoscoordinator
WORKDIR /chaoscoordinator
RUN make modules
RUN make build

#second stage
FROM gcr.io/distroless/base
WORKDIR /
COPY --from=builder /chaoscoordinator/bin/chaos /chaos
COPY --from=builder /chaoscoordinator/scripts/ /scripts