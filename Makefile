NAME = gAD-System
REPO = github.com/makometr/gAD-System
BUILD_DIR ?= build

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

proto-gen:
	@printf "${OK_COLOR}==> Generating proto code${NO_COLOR}\n"
	@protoc -I=. --go_out=. ./api/proto3/calculation_service/v1beta1/message.proto