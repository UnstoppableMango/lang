#syntax=docker/dockerfile:1
FROM ubuntu:noble AS llvm

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
