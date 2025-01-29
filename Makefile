IMAGE_VERSION = $(shell git rev-parse --short HEAD)

REGISTRY = docker.io/stdoutt

SCHEDULER_IMAGE = ${REGISTRY}/scheduler-plugins:${IMAGE_VERSION}
CONTROLLER_IMAGE = ${REGISTRY}/scheduler-plugins-controller:${IMAGE_VERSION}

.PHONY: build deploy

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/scheduler-plugins cmd/scheduler/main.go    
	CGO_ENABLED=0 GOOS=linux go build -o bin/controller cmd/controller/controller.go    

build-image:
	docker build -f build/scheduler/Dockerfile -t ${SCHEDULER_IMAGE} .
	docker build -f build/controller/Dockerfile -t ${CONTROLLER_IMAGE} .

push-image:
	docker push ${SCHEDULER_IMAGE}
	docker push ${CONTROLLER_IMAGE}

all:
	@$(MAKE) build
	@$(MAKE) build-image
	@$(MAKE) push-image