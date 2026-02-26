include .env

run:
	go run .

docker-build:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--tag registry.gitlab.com/afonsodemori/afonsodev-api:latest \
		--push \
		.
