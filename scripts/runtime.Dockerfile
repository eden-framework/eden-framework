FROM golang:1.14-alpine

RUN sed -i "s|http://dl-cdn.alpinelinux.org|http://mirrors.aliyun.com|g" /etc/apk/repositories

RUN apk add --no-cache --virtual .build-deps \
	    curl bash vim htop

ENV GODEBUG=netdns=cgo

WORKDIR /var/services
