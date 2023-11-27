.DEFAULT_GOAL := build

IMAGE ?= telegram-bot-test:latest

.PHONY: build
build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--progress=plain \
		--output "type=docker,push=false" \
		--tag $(IMAGE) \
		--file build/Dockerfile \
		.

.PHONY: run
run:
	docker run -e MODEL_PATH="/model" -e TELEGRAM_BOT_TOKEN="6853139137:AAHt9ab7_9dmcC-zq0cF4CbCrAaSxQKTPu8" telegram-bot-test:latest run