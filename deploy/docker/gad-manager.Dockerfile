FROM golang:1.18 AS builder
RUN apt-get update -y && apt-get upgrade -y
COPY ./cmd/gad-manager/main.go /go/src/gAD-System/cmd/gad-manager/main.go
COPY ./services/gad-manager/ /go/src/gAD-System/services/gad-manager/
COPY ./Makefile /go/src/gAD-System/Makefile
COPY ./go.mod /go/src/gAD-System/go.mod
COPY ./go.sum /go/src/gAD-System/go.sum
WORKDIR /go/src/gAD-System
RUN make build-gad-manager

FROM alpine:latest AS runner
RUN apk -U upgrade
COPY --from=builder /go/src/gAD-System/build/gad-manager.exe .
ENTRYPOINT [ "./gad-manager.exe" ]