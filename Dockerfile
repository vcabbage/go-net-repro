FROM golang:latest

RUN apt-get update && apt-get install strace && apt-get clean

CMD go build -o repro repro && while strace -e network ./repro; do :; done