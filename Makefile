.PHONY: docker
docker:
	docker buildx build --push --platform linux/arm64,linux/amd64 -t dukeman/wikileet:latest .

.PHONY: helm
helm:
	helm upgrade --install -n wikileet-api wikileet deploy/wikileet

.PHONY: docs
docs:
	swag fmt ./...
	swag init