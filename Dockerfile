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

RUN echo -e "http://mirrors.aliyun.com/alpine/v3.11/main\nhttp://mirrors.aliyun.com/alpine/v3.11/community" > /etc/apk/repositories \
    && apk add -U tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 
    
WORKDIR /

COPY --from=build /app/k8s-api-service /k8s-api-service
COPY --from=build /app/config.yaml /config.yaml

EXPOSE 8080

ENTRYPOINT ["/k8s-api-service"]