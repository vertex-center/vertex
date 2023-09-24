FROM golang:1.20-alpine AS build-stage

WORKDIR /app

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /vertex -ldflags="-X 'main.version=$(git describe --tags --always --dirty)' -X 'main.commit=$(git rev-parse HEAD)' -X 'main.date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'"

FROM build-stage AS test-stage
RUN go test -v ./...

FROM alpine AS run-stage

WORKDIR /

COPY --from=build-stage /vertex /vertex

EXPOSE 6130

CMD /vertex
