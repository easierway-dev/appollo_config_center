FROM golang:1.15.3 as builder

ARG SSH_KEY

WORKDIR /data
ENV GOPATH=/root/go
ENV GOCACHE=/root/gocache/.cache/go-build
ENV THIRD_DIR=/data
ENV GO111MODULE=on
ENV GOPROXY=direct
ENV GOSUMDB=off

# set prv key 
RUN mkdir -p /root/.ssh && \
        chmod 0700 /root/.ssh && \
        ssh-keyscan github.com > /root/.ssh/known_hosts && \
        ssh-keyscan gitlab.mobvista.com >> /root/.ssh/known_hosts; \
        echo "$SSH_KEY" > /root/.ssh/id_rsa && \
        chmod 600 /root/.ssh/* 
# set git config
RUN git config --global url."git@gitlab.mobvista.com:".insteadOf "http://gitlab.mobvista.com" ; 

COPY . /data/appollo_config_center/

RUN cd /data/appollo_config_center/ && make build; 

FROM alpine:3.11
WORKDIR /data/appollo_config_center
ENV PATH="${PATH}:/sbin"
RUN apk update && apk upgrade && apk add bash && apk add --no-cache libc6-compat
COPY --from=builder /data/appollo_config_center/deployments ./
CMD ["./run"]
