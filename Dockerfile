# syntax = docker/dockerfile:experimental
FROM ubuntu:20.04
# MAINTAINER Wukun 1491580498@qq.com
RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && sed -i s@/security.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    # && apt-get update --fixing-missing \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt install software-properties-common -y \
    && add-apt-repository ppa:longsleep/golang-backports \
    && apt-get install gcc libc6-dev git lrzsz -y \
    && apt-get install bash \
    # # # # 安装Go环境
    && DEBIAN_FRONTEND=noninteractive apt-get -y -qq install --no-install-recommends golang 

# # # 配置环境变量 
ENV GOROOT=/usr/lib/go
ENV PATH=$PATH:/usr/lib/go/bin
ENV GOPATH=/root/go
ENV PATH=$GOPATH/bin/:$PATH
ENV GOPROXY=https://goproxy.cn

RUN DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends protobuf-compiler -y \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


# RUN go get github.com/hashicorp/memberlist \
#     && go get stathat.com/c/consistent
# # # 定制工作目录
WORKDIR /root/go/src/geecache/
# # # 装载项目
COPY ./main.go /root/go/src/geecache/
COPY ./cache/* /root/go/src/geecache/cache/
COPY ./http/* /root/go/src/geecache/http/
COPY ./proto/* /root/go/src/geecache/proto/
COPY ./run.sh /root/go/src/geecache/
RUN go mod init geecache && go mod tidy && chmod +x /root/go/src/geecache/run.sh \
    && cd /root/go/src/geecache/proto/ && protoc --go_out=. --go-grpc_out=. node.proto
# EXPOSE 12345
# EXPOSE 12346
# EXPOSE 12347
# # # # # 运行
ENTRYPOINT ["/root/go/src/geecache/run.sh"]

