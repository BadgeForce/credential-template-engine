OUT := credential-template-engine
PKG := github.com/BadgeForce/credential-template-engine/cmd
DOCKERFILE := ./build/package/Dockerfile
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: run

docker-image:
	docker build -f ${DOCKERFILE} -t ${OUT} .

deps:
	GO111MODULE=off go get ${PKG}

build:
	GO111MODULE=off go build -i -v -o ${OUT} ${PKG}

test:
	GO111MODULE=off @go test -short ${PKG_LIST}

vet:
	GO111MODULE=off @go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

runtp:
	./${OUT}-v${VERSION} --validator-endpoint $(validator)

out:
	@echo ${OUT}-v${VERSION}

protos:
	protoc -I ./protos ./protos/template.proto ./protos/payload.proto ./protos/transaction_receipts.proto --go_out=./core/template_pb

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run protos runtp build docker-image vet lint out deps