.PHONY: docker
docker:
	docker buildx build --push --platform linux/arm64,linux/amd64 -t dukeman/wikileet-api:latest .
	docker buildx build --push --platform linux/arm64,linux/amd64 -t dukeman/wikileet-ui:latest ./frontend/wikileet-ui

.PHONY: helm
helm:
	helm upgrade --install -n wikileet wikileet deploy/wikileet

.PHONY: docs
docs:
	swag fmt ./...
	swag init