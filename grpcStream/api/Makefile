PACKAGE_NAME := github.com/brinkpku/training_center/grpcStream/api
PROTOS_WITH_PATH := $(shell find $$GOPATH/src/$(PACKAGE_NAME) -maxdepth 1 \( -name "*.proto" \))
ENV_GOPATH = ${GOPATH}/src

.DEFAULT: all
all: build

build: clean 
	for i in $(PROTOS_WITH_PATH); do \
	protoc --proto_path=$$GOPATH/src:\
	$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway:\
	$$GOPATH/src/github.com/googleapis/googleapis:. \
	--go_out=:$(ENV_GOPATH) \
	--go-grpc_out=:$(ENV_GOPATH) \
	$$i; \
	done

clean:
	rm -rf $(CURDIR)/../pb

.PHONY: all build swagger prepare clean
# $(VERBOSE).SILENT: