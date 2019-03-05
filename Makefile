OUT := credential-template-engine
PKG := ./cmd
DOCKERFILE := ./build/package/Dockerfile
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: run

docker-image:
	docker build -f ${DOCKERFILE} -t ${OUT} .

deps:
	GO111MODULES=on vgo get -v ${PKG}

build:
	GO111MODULES=on vgo build -v -o ${OUT} ${PKG}

test:
	GO111MODULES=on vgo test -short ${PKG_LIST}

vet:
	GO111MODULES=on vgo vet ${PKG_LIST}

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
	protoc -I ./protos ./protos/template.proto ./protos/payload.proto ./protos/transaction_receipts.proto --js_out=import_style=commonjs,binary:./dev-cli/proto

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run protos runtp build docker-image vet lint out deps