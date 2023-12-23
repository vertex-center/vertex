FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS build-stage

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG APP_KIND=main
RUN test -d ./cmd/"$APP_KIND" || (echo "APP_KIND is not set to a valid app kind" && exit 1)

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app -ldflags="-w -s -X 'main.version=$(git describe --tags --always --dirty)' -X 'main.commit=$(git rev-parse HEAD)' -X 'main.date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" ./cmd/"$APP_KIND"

FROM scratch AS run-stage

WORKDIR /

COPY --from=build-stage /app /app
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80 7500 7501 7502 7504 7505 7506 7508 7510 7512 7514 7516 7518

ENTRYPOINT [ "/app" ]
