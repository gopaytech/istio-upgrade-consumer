.EXPORT_ALL_VARIABLES:

BIN_DIR := ./bin
OUT_DIR := ./output
$(shell mkdir -p $(BIN_DIR) $(OUT_DIR))

APP_NAME=istio-upgrade-consumer
PACKAGE=github.com/gopaytech/istio-upgrade-consumer
TRACK=stable
ENV?=stg
IMAGE_REGISTRY=ghcr.io/gopaytech
IMAGE_NAME=$(IMAGE_REGISTRY)/istio-upgrade-consumer
IMAGE_TAG?=$(shell git rev-parse --short HEAD)

CURRENT_DIR=$(shell pwd)
VERSION=$(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --exact-match --tags HEAD 2>/dev/null)
GIT_TREE_STATE=$(shell if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi)

DEPLOYMENT_TIMEOUT=600

STATIC_BUILD?=true

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \
  -X ${PACKAGE}.gitTreeState=${GIT_TREE_STATE}

ifeq (${STATIC_BUILD}, true)
override LDFLAGS += -extldflags "-static"
endif

ifneq (${GIT_TAG},)
IMAGE_TAG=${GIT_TAG}
IMAGE_TRACK=stable
LDFLAGS += -X ${PACKAGE}.gitTag=${GIT_TAG}
else
IMAGE_TAG?=$(GIT_COMMIT)
IMAGE_TRACK=latest
endif

.PHONY: image.build
image.build:
	echo "building container image"
	DOCKER_BUILDKIT=1 docker build \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		--build-arg GITCONFIG=$(GITCONFIG) --build-arg BUILDKIT_INLINE_CACHE=1 .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest

.PHONY: image.release
image.release:
	echo "pushing container image"
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: build.binaries
build.binaries:
	CGO_ENABLED=0 GO111MODULE=on go build -a -ldflags '${LDFLAGS}' -o ${BIN_DIR}/istio-upgrade-consumer ./main.go
