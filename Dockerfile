FROM golang:1.20-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /vertex

FROM build-stage AS test-stage
RUN go test -v ./...

FROM alpine AS run-stage

WORKDIR /

COPY --from=build-stage /vertex /vertex

EXPOSE 6130

CMD /vertex
