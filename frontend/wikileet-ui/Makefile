.PHONY: docker
docker:
	# docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx build --push --platform linux/arm64 -t dukeman/wikileet-ui:latest .
