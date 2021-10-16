FROM golang:1.14-alpine

COPY ./scripts/proc_id.go.patch /proc_id.go

RUN cd $(go env GOROOT)/src/runtime \
    && mv /proc_id.go . \
    && go install

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories

RUN apk add --no-cache curl git openssh wget unzip gcc libc-dev

ENV CGO_ENABLED 0
ENV GOSUMDB off
ENV GOPROXY https://goproxy.cn

COPY ./ /go/src/eden-framework/eden-framework
RUN cd /go/src/eden-framework/eden-framework/cmd/eden && go install

ADD ./scripts/eden/.eden.yaml /root/.eden.yaml
