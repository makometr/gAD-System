FROM golang:1.18 AS builder
RUN apt-get update -y && apt-get upgrade -y
WORKDIR /go/src/gAD-System
COPY ./cmd/gad-manager/main.go ./cmd/gad-manager/main.go
COPY ./services/gad-manager/ ./services/gad-manager/
COPY ./internal/ ./internal/
COPY ./Makefile ./Makefile
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN make build-gad-manager

FROM alpine:latest AS runner
RUN apk -U upgrade
COPY --from=builder /go/src/gAD-System/build/gad-manager.exe ./
EXPOSE 8080
ENTRYPOINT [ "./gad-manager.exe" ]