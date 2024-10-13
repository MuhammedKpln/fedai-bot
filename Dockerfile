FROM golang:latest

COPY . /opt/fedai

RUN apt update && apt install ffmpeg -y


WORKDIR /opt/fedai
RUN go run ./scripts/compile_plugins.go
RUN go build main.go