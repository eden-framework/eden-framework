FROM alpine:3

RUN sed -i "s|http://dl-cdn.alpinelinux.org|http://mirrors.aliyun.com|g" /etc/apk/repositories

RUN apk add --no-cache curl git openssh wget unzip \
	&& cd /tmp \
	&& wget https://releases.rancher.com/cli2/v2.4.3/rancher-linux-amd64-v2.4.3.tar.gz \
	&& tar -zxvf rancher-linux-amd64-v2.4.3.tar.gz \
	&& mv rancher-v2.4.3/rancher /bin/rancher \
	&& chmod +x /bin/rancher

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl \
	&& mv ./kubectl /bin/kubectl \
	&& chmod +x /bin/kubectl

ADD ./scripts/kube_config/config /root/.kube/config
