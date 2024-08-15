FROM registry-cn-hangzhou.ack.aliyuncs.com/dev/alpine:3.18-update@sha256:f2f97bbc6f3173044292cbf8d306ca31b4356a5f2eb54701fc97f64065b412cf
WORKDIR /
COPY benchmark /benchmark
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add tzdata
RUN echo "hosts: files dns" > /etc/nsswitch.conf
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
RUN apk del tzdata
ENTRYPOINT ["/benchmark"]
