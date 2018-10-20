OUT := credential-template-engine
PKG := github.com/BadgeForce/credential-template-engine/cmd
DOCKERFILE := ./build/package/Dockerfile
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: run

docker-image:
	docker build -f ${DOCKERFILE} -t ${OUT} .

goget:
	go get ${PKG}

build:
	go build -i -v -o ${OUT} -ldflags="-X main.repoVersion=${VERSION}" ${PKG}

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

runtp:
	./${OUT}-v${VERSION} --validator-endpoint $(validator)

out:
	@echo ${OUT}-v${VERSION}

protos:
	rm ./core/proto/pending_props_pb/*
	rm ./dev-cli/proto/*
	protoc -I ./protos ./protos/earning.proto ./protos/payload.proto --go_out=./core/proto/pending_props_pb
	protoc -I ./protos ./protos/earning.proto ./protos/payload.proto --js_out=import_style=commonjs,binary:./dev-cli/proto

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run protos runtp build docker-image vet lint out goget