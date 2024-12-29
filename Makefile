# Copyright 2024 Ava AI. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

TAG?=latest
VERSION:="v0.1.0"

NAME:=ava
DOCKER_REPOSITORY:=matthish29
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)

ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --abbrev=0 --dirty --always --tags | sed 's/-/./g')
endif
BUILDAPP = "./bin/ava"

GIT_COMMIT:=$(shell git rev-parse HEAD)

## build: Build binaries by default
.PHONY: build
build: 
	@echo "$(shell go version)"
	@echo "===========> Building binary $(BUILDAPP) *[Git Info]: $(VERSION)-$(GIT_COMMIT)"
	@export CGO_ENABLED=0 && go build -o $(BUILDAPP) -ldflags "-s -w -X main.version=dev -X main.commit=$$(git rev-parse --short HEAD) -X main.date=$$(date +%FT%TZ)" $(BUILDFILE)

swagger:
	swag init

tidy:
	go mod tidy

vet:
	go vet ./...

fmt:
	gofmt -l -s -w ./
	goimports -l -w ./

version-set:
	@next="$(TAG)" && \
	current="$(VERSION)" && \
	/usr/bin/sed -i '' "s/$$current/$$next/g" pkg/version/version.go && \
	echo "Version $$next set in code"

build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

build-init-container:
	docker build -t $(DOCKER_IMAGE_NAME)-init:$(VERSION) . -f Dockerfile.init

prisma-push:
	cd prisma && go run github.com/steebchen/prisma-client-go db push

LICENSE_TEMPLATE ?= ./hacks/LICENSE_TEMPLATE

# Linux command settings-CODE DIRS Copyright
CODE_DIRS := ./pkg ./cmd

## copyright.add: Add the boilerplate headers for all files
.PHONY: copyright.add
copyright.add:
	@echo "===========> Adding $(LICENSE_TEMPLATE) the boilerplate headers for all files"
	@addlicense -y $(shell date +"%Y") -v -c "Ava AI." -f $(LICENSE_TEMPLATE) $(CODE_DIRS)
	@echo "===========> End the copyright is added..."