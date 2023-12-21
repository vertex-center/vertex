FROM golang:1.21-alpine AS build-stage

WORKDIR /app
COPY . ./

ARG APP_ID=undefined
RUN test -d ./apps/$APP_ID || (echo "APP_ID is not set to a valid app ID" && exit 1)

RUN apk add git

RUN go mod download
RUN go build -o /$APP_ID -ldflags="-X 'main.version=$(git describe --tags --always --dirty)' -X 'main.commit=$(git rev-parse HEAD)' -X 'main.date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" ./apps/$APP_ID/cmd/main

FROM alpine AS run-stage

WORKDIR /

ARG APP_ID
ENV APP_ID=$APP_ID

COPY --from=build-stage /$APP_ID /$APP_ID

EXPOSE 8080

ENTRYPOINT /$APP_ID
