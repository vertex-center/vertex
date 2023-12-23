FROM golang:1.21-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -buildvcs=false -o /vertex ./cmd/main
RUN go build -buildvcs=false -o /vertex-kernel ./cmd/kernel

FROM alpine AS run-stage

RUN apk add --no-cache git sudo
RUN addgroup -g 1000 nexa && \
    adduser -u 1000 -G nexa -s /bin/sh -D nexa && \
    echo "nexa ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers && \
    mkdir -p /home/nexa/vertex && \
    chown -R nexa:nexa /home/nexa

USER nexa:nexa

WORKDIR /home/nexa/vertex

COPY --from=build-stage /vertex /home/nexa/vertex/vertex
COPY --from=build-stage /vertex-kernel /home/nexa/vertex/vertex-kernel

EXPOSE 80 6130 6131

CMD sudo ./vertex-kernel -user nexa
