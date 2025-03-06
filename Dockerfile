#syntax=docker/dockerfile:1
FROM golang:1.24 AS llvm

RUN apt-get update && \
    apt-get install -y apt-utils git make cmake clang-15 ninja-build && \
    rm -rf \
        /var/lib/apt/lists/* \
        /var/log/* \
        /var/tmp/* \
        /tmp/*

COPY lib/llvm/Makefile /lang/Makefile
RUN make -C /lang llvm-project

FROM llvm AS llvm-build
RUN make -C /lang llvm-build

FROM llvm-build AS build

WORKDIR /lang

COPY go.mod go.sum ./

COPY cmd/ir ./
RUN go build ./
