FROM golang:1.18 AS builder
RUN apt-get update -y && apt-get upgrade -y
WORKDIR /go/src/gAD-System
COPY ./cmd/calc-controller/main.go ./cmd/calc-controller/main.go
COPY ./services/calc-controller/ ./services/calc-controller/
COPY ./internal/ ./internal/
COPY ./Makefile ./Makefile
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN make build-calc-controller

FROM alpine:latest AS runner
RUN apk -U upgrade
COPY --from=builder /go/src/gAD-System/build/calc-controller.exe ./
EXPOSE 50051
ENTRYPOINT [ "./calc-controller.exe" ]