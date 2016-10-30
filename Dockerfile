FROM golang:latest

RUN apt-get update && apt-get install strace && apt-get clean

CMD go build -o repro repro && while strace -f -e close,network ./repro; do :; done