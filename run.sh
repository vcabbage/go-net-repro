#!/bin/bash

docker build -t go-net-repro . && \
    docker run -it --rm -v "$(pwd):/go/src/repro" --privileged go-net-repro