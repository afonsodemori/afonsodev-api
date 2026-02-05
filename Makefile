.PHONY: run
run:
	@set -a; \
	[ -f .env ] && . ./.env; \
	set +a; \
	go run main.go

.PHONY: gemini
gemini:
	docker sandbox run gemini .
