FROM --platform=$BUILDPLATFORM golang:latest AS build
WORKDIR /src
COPY . .
RUN go mod download
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -installsuffix cgo -o vehicleHandlerService ./cmd/api

FROM --platform=$BUILDPLATFORM alpine:latest AS final
WORKDIR /root/

COPY --from=build /src/vehicleHandlerService .

RUN mkdir product-definitions

COPY --from=build /src/cmd/api/product-definitions ./product-definitions

EXPOSE 80

CMD ["./vehicleHandlerService"]