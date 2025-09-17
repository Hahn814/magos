FROM golang:1-bookworm

ARG GO_LOG_VERBOSITY_LEVEL=99
ARG GO_LOG_SEVERITY_LEVEL=info

RUN useradd -ms /bin/bash magos \
    && apt update \
    && apt install -y protobuf-compiler

USER magos

WORKDIR /usr/src/app
COPY . .

RUN mkdir -p /home/magos/go/bin \
    && go install github.com/spf13/cobra-cli@latest \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH=/home/magos/go/bin:$PATH
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=${GO_LOG_VERBOSITY_LEVEL}
ENV GRPC_GO_LOG_SEVERITY_LEVEL=${GO_LOG_SEVERITY_LEVEL}