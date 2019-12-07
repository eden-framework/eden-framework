FROM registry.profzone.net:5000/profzone/golang:latest

RUN sed -i "s|http://dl-cdn.alpinelinux.org|http://mirrors.aliyun.com|g" /etc/apk/repositories

RUN apk add --no-cache curl git openssh wget unzip gcc libc-dev

ENV GODEBUG=netdns=cgo GOPROXY=https://goproxy.io
