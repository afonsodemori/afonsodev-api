.PHONY: run
run:
	@set -a; \
	[ -f .env ] && . ./.env; \
	set +a; \
	go run .

.PHONY: gemini
gemini:
	docker sandbox run gemini .

.PHONY: docker-build
docker-build:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t registry.gitlab.com/afonsodemori/afonso-dev-api:latest \
		--push \
		.
