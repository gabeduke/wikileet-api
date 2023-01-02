.PHONY: docker
docker:
	docker buildx build --push --platform linux/arm64,linux/amd64 -t dukeman/wikileet-api:latest .

.PHONY: helm
helm:
	helm upgrade --install -n wikileet-test wikileet-test deploy/wikileet

.PHONY: docs
docs:
	swag fmt ./...
	swag init