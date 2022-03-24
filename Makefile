NAME = gAD-System
REPO = github.com/makometr/gAD-System
BUILD_DIR ?= build

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

all: clean deps build

proto-gen:
	@echo "${OK_COLOR}==> Generating proto code${NO_COLOR}\n"
	@protoc -I=. --go_out=. ./api/proto3/calculation_service/v1beta1/message.proto

deps:
	git config --global http.https://gopkg.in.followRedirects true
	@echo "${OK_COLOR}==> Downloading dependencies${NO_COLOR}\n"
	@go mod vendor

clean:
	@echo "${OK_COLOR}==> Cleaning... ${NO_COLOR}\n"
	@rm -rf ./build

lint:
	@echo "${OK_COLOR}==> Linting... ${NO_COLOR}\n"
	@golangci-lint -c .golangci.yml run ./..

dev-fix-lint:
	@echo "${OK_COLOR}==> Fixing... ${NO_COLOR}\n"
	@gofmt -s -w .
	@goimports -l -w .

build-gad-mamanger:
	@echo "${OK_COLOR}==> Building gad-manager${NO_COLOR}\n"
	@CGO_ENABLED=0 go build -o ${BUILD_DIR}/gad-manager.app cmd/gad-manager/main.go

build-calc-controller:
	@echo "${OK_COLOR}==> Building calc-controller${NO_COLOR}\n"
	@CGO_ENABLED=0 go build -o ${BUILD_DIR}/calc-controller.app cmd/calc-controller/main.go

build: build-gad-mamanger build-calc-controller