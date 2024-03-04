FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS build-stage

WORKDIR /build

COPY . .

RUN apk add --no-cache git

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

ARG APP_KIND=main

EXPOSE 80 7500 7502 7504 7505 7506 7508 7510 7512 7514 7516 7518

ENTRYPOINT CompileDaemon -build='go build -o binary ./cmd/'"$APP_KIND" -directory=. -command=./binary
