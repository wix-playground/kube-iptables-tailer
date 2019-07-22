FROM golang:1.11.5 as builder

ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/wix-playground/kube-iptables-tailer
COPY . $GOPATH/src/github.com/wix-playground/kube-iptables-tailer
RUN make build

FROM ubuntu
LABEL maintainer="Saifuding Diliyaer <sdiliyaer@box.com>"
WORKDIR /root/
COPY --from=builder /go/src/github.com/wix-playground/kube-iptables-tailer/kube-iptables-tailer /kube-iptables-tailer
