#syntax=docker/dockerfile:1
FROM golang:1.24 AS llvm

RUN apt-get update && \
    apt-get install -y apt-utils git make cmake clang-15 ninja-build && \
    rm -rf \
        /var/lib/apt/lists/* \
        /var/log/* \
        /var/tmp/* \
        /tmp/*

# COPY lib/llvm/Makefile /lang/Makefile
RUN git clone --depth=1 --branch llvmorg-19.1.7 https://github.com/llvm/llvm-project.git

FROM llvm AS llvm-build
WORKDIR /llvm-build
RUN cmake -G Ninja /llvm-project \
		"-DLLVM_TARGETS_TO_BUILD=X86;ARM;AArch64;WebAssembly"

FROM llvm-build AS build

WORKDIR /lang

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ir ./
RUN go build -tags=byollvm ./
