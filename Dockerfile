# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-alpine AS build

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . /app/

RUN go build -o /app/k8s-api-service

## Deploy
FROM alpine

WORKDIR /

COPY --from=build /app/k8s-api-service /k8s-api-service
COPY --from=build /app/config.yaml /config.yaml

EXPOSE 8080

ENTRYPOINT ["/k8s-api-service"]