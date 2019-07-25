FROM golang:latest as builder
RUN go get github.com/spf13/cobra

WORKDIR /chaoscoordinator
COPY cmd/loadbalancer/loadbalancer.go /go/src/github.com/UtheMan/chaosCoordinator/cmd/loadbalancer/loadbalancer.go
COPY cmd/loadbalancer/kill/kill_loadbalancer.go /go/src/github.com/UtheMan/chaosCoordinator/cmd/loadbalancer/kill/kill_loadbalancer.go
COPY cmd/vm/vm.go /go/src/github.com/UtheMan/chaosCoordinator/cmd/vm/vm.go
COPY cmd/vm/kill/kill_vm.go /go/src/github.com/UtheMan/chaosCoordinator/cmd/vm/kill/kill_vm.go

COPY cmd/ .
RUN go build chaos.go

# FROM alpine:3.9
CMD ["./chaos;"]
# COPY --from=builder /chaos .
 