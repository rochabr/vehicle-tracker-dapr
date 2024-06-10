RELEASE_VERSION  = 0.1.1
SERVICE_NAME = shipment-handler
IMAGE_NAME = shipment-handler
REGISTRY=rochabr
APP_PORT=5100
DAPR_HTTP_PORT=5180
DAPR_GRPC_PORT=5101

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependencies 
	go mod tidy



.PHONY: run
run: tidy ## Runs uncompiled code in Dapr
	dapr run \
		--app-id $(SERVICE_NAME) \
		--app-port ${APP_PORT} \
		--app-protocol http \
		--dapr-http-port ${DAPR_HTTP_PORT} \
		--dapr-grpc-port ${DAPR_GRPC_PORT} \
		--resources-path ../../components/local \
		go run ./cmd/api

.PHONY: build_image
build_image: tidy ## Builds local release binary
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${REGISTRY}/${IMAGE_NAME}:${RELEASE_VERSION} --push .