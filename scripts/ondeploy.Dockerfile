FROM golang:1.14-alpine

COPY ./scripts/proc_id.go.patch /proc_id.go

RUN cd $(go env GOROOT)/src/runtime \
    && mv /proc_id.go . \
    && go install

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories

RUN apk add --no-cache curl git openssh wget unzip
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl \
	&& mv ./kubectl /bin/kubectl \
	&& chmod +x /bin/kubectl

ENV CGO_ENABLED 0
ENV GOSUMDB off
ENV GOPROXY https://goproxy.cn

COPY ./ /go/src/eden-framework/eden-framework
RUN cd /go/src/eden-framework/eden-framework/cmd/eden && go install

ADD ./scripts/eden/.eden.yaml /root/.eden.yaml
ADD ./scripts/kube_config/config /root/.kube/config
ADD ./scripts/kube_config/config_staging /root/.kube/config_staging
ADD ./scripts/kube_config/config_test /root/.kube/config_test
ADD ./scripts/kube_config/config_demo /root/.kube/config_demo
