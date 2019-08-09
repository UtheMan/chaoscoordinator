FROM golang:latest as builder

COPY . /chaoscoordinator
WORKDIR /chaoscoordinator
RUN make modules
RUN make build

#second stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /chaoscoordinator/bin/chaos /chaos
#CMD ["./chaos",  "vm", "kill"]