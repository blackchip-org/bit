FROM ubuntu:xenial

RUN apt-get update
RUN apt-get install -y \
    build-essential \
    bison \
    m4 \
    texinfo \
    gawk \
    wget
RUN update-alternatives --install /bin/sh sh /bin/bash 100
