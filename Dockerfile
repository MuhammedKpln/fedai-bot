FROM golang:latest

RUN sudo apt install ffmpeg

RUN go run ./scripts/compile_plugins.go
RUN go build main.go