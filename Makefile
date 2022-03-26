NAME = gAD-System
REPO = github.com/makometr/gAD-System
BUILD_DIR ?= build

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

all: clean build

proto-calculator:
	@echo "${OK_COLOR}==> Generating proto code for manager${NO_COLOR}\n"
	@protoc -I=. --go_out=. --go-grpc_out=. ./api/proto/grpc/calculator/service.proto

proto-expression:
	@echo "${OK_COLOR}==> Generating proto code for rmq${NO_COLOR}\n"
	@protoc -I=. --go_out=. ./api/proto/expression/event.proto

docker-build-gad-manager:
	@echo "${OK_COLOR}==> Building docker image for gad-manager${NO_COLOR}\n"
	@docker build --no-cache -t gad-manager:v0.1 ./ -f ./docker/gad-manager.Dockerfile

docker-build-calc-controller:
	@echo "${OK_COLOR}==> Building docker image for calc-controller${NO_COLOR}\n"
	@docker build --no-cache -t calc-controller:v0.1 ./ -f ./docker/calc-controller.Dockerfile

docker-build-rabbitmq:
	@echo "${OK_COLOR}==> Building docker image for rabbitmq${NO_COLOR}\n"
	@docker build --no-cache -t gad-rabbitmq:v0.1 ./ -f ./docker/rmq/rabbitmq.Dockerfile

docker-build: docker-build-gad-manager docker-build-calc-controller docker-build-rabbitmq

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

build-gad-manager:
	@echo "${OK_COLOR}==> Building gad-manager${NO_COLOR}\n"
	@CGO_ENABLED=0 go build -o ${BUILD_DIR}/gad-manager.exe cmd/gad-manager/main.go

build-calc-controller:
	@echo "${OK_COLOR}==> Building calc-controller${NO_COLOR}\n"
	@CGO_ENABLED=0 go build -o ${BUILD_DIR}/calc-controller.exe cmd/calc-controller/main.go

docker-build-gad-manager:
	@docker build -t gad-test ./ -f docker\gad-manager.Dockerfile


build: build-gad-manager build-calc-controller