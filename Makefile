include .env

# development

run:
	go run ./cmd/server

# cloudflare tunnel

tunnel-init:
	cloudflared tunnel login
	cloudflared tunnel create afonsodev-api || true
	cloudflared tunnel route dns --overwrite-dns afonsodev-api dev-afonsodev-api.afonso.dev

tunnel-run:
	cloudflared tunnel --config=.devcontainer/cloudflared/config.yml run afonsodev-api

# docker
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

docker-build:
	@echo "Building Docker image $(IMAGE_NAME) with GO_VERSION=$(GO_VERSION)"
	@ARGS=; \
	for tag in $$(echo "$(IMAGE_TAGS)" | tr ',' ' '); do \
		echo "Adding tag: $$tag"; \
		ARGS="$$ARGS --tag $(IMAGE_NAME):$$tag"; \
	done; \
	docker buildx build . \
		--platform linux/amd64,linux/arm64 \
		--file docker/Dockerfile.production \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		$$ARGS

docker-push: docker-build
	@for tag in $$(echo "$(IMAGE_TAGS)" | tr ',' ' '); do \
		echo "Pushing tag: $$tag"; \
		docker push $(IMAGE_NAME):$$tag; \
	done
