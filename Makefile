.DEFAULT_GOAL := build

IMAGE ?= telegram-bot-test:latest

.PHONY: build
build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--output "type=docker,push=false" \
		--tag $(IMAGE) \
		--file build/Dockerfile \
		.