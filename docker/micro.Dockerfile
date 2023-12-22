FROM golang:1.21-alpine AS build-stage

WORKDIR /build

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG APP_ID=undefined
RUN test -d ./apps/$APP_ID || (echo "APP_ID is not set to a valid app ID" && exit 1)

ARG APP_KIND=main
RUN test -d ./apps/$APP_ID/cmd/$APP_KIND || (echo "APP_KIND is not set to a valid app kind" && exit 1)

RUN CGO_ENABLED=0 go build -o /app -ldflags="-w -s -X 'main.version=$(git describe --tags --always --dirty)' -X 'main.commit=$(git rev-parse HEAD)' -X 'main.date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" ./apps/$APP_ID/cmd/$APP_KIND

FROM scratch AS run-stage

WORKDIR /

COPY --from=build-stage /app /app
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT [ "/app" ]
