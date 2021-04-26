FROM golang:1.15.3 as gobase
ARG ssh_prv_key
ENV GOPATH=/root/go
WORKDIR /data
# 获取密钥用于依赖拉取代码
COPY dsp_server_test.go ./
RUN mkdir -p /root/.ssh && \
        chmod 0700 /root/.ssh && \
        ssh-keyscan github.com > /root/.ssh/known_hosts && \
        ssh-keyscan gitlab.mobvista.com >> /root/.ssh/known_hosts; \
        echo "$ssh_prv_key" > /root/.ssh/id_rsa && \
        chmod 600 /root/.ssh/*;\
        git clone git@gitlab.mobvista.com:mvbjqa/appollo_config_center.git && \
        cd appollo_config_center && \
        sed -i '/golang:/,+9d' Makefile && \
        make build;

FROM alpine:3.11
WORKDIR /data
COPY --from=gobase /data/appollo_config_center/deployments   /data/appollo_config_center/deployments

