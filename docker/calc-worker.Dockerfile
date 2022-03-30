FROM golang:1.18 AS builder
RUN apt-get update -y && apt-get upgrade -y
WORKDIR /go/src/gAD-System
COPY ./cmd/calc-worker/main.go ./cmd/calc-worker/main.go
COPY ./services/calc-worker/ ./services/calc-worker/
COPY ./internal/ ./internal/
COPY ./Makefile ./Makefile
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN make build-calc-worker

FROM alpine:latest AS runner
RUN apk -U upgrade
COPY --from=builder /go/src/gAD-System/build/calc-worker.exe ./
# EXPOSE 50051
ENTRYPOINT [ "./calc-worker.exe" ]