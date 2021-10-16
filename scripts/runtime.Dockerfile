FROM golang:1.14-alpine

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories

RUN apk add --no-cache --virtual .build-deps \
	    curl bash vim htop

ENV GODEBUG=netdns=cgo

WORKDIR /var/services
