GO := $(shell command -v go 2> /dev/null)
DOCKER := $(shell command -v docker 2> /dev/null)
RELEASE_VERSION ?= $(shell bash -c 'read -s -p "Release version: " pwd')

.PHONY: build
build:
ifndef GO
	@echo "go is required!"
endif
	@echo "building ddnsclient..."
	@go build -o bin/ddnsclient cmd/main.go
	@echo "module built!"

.PHONY: deploy-image-latest
deploy-image-latest:
ifndef DOCKER
	@echo "docker is required!"
endif
	@echo "Pushing image ddnsclient:latest to docker hub..."
	@docker push ddnsclient:latest
	@echo "Image pushed!"

.PHONY: deploy-image-release
deploy-image-release:
ifndef DOCKER
	@echo "docker is required!"
endif
	@echo "Pushing image ddnsclient:$(RELEASE_VERSION) to docker hub..."
	@docker push ddnsclient:$(RELEASE_VERSION)
	@echo "Image pushed!"
    
