include .env

run:
	go run .

tunnel-init:
	cloudflared tunnel login
	cloudflared tunnel create dev-afonsodev-api
	cloudflared tunnel route dns --overwrite-dns dev-afonsodev-api dev-afonsodev-api.afonso.dev

tunnel-run:
	cloudflared tunnel --config=.devcontainer/cloudflared/config.yml run dev-afonsodev-api

docker-build:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--tag registry.gitlab.com/afonsodemori/afonsodev-api:latest \
		--push \
		.
