FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS build-stage

WORKDIR /vertex

RUN apk add --no-cache git

ARG APP_KIND=main
ARG TARGETOS
ARG TARGETARCH

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    cd server && \
    CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -o /app -ldflags="-w -s -X 'main.version=$(git describe --tags --always --dirty)' -X 'main.commit=$(git rev-parse HEAD)' -X 'main.date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" ./cmd/"$APP_KIND"

FROM scratch AS run-stage

WORKDIR /

COPY --from=build-stage /app /app
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80 7500 7502 7504 7505 7506 7508 7510 7512 7514 7516 7518

ENTRYPOINT [ "/app" ]
