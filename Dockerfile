FROM golang:1.23.1-bookworm as pluginBuilder
WORKDIR /app
COPY . .
RUN go mod download
RUN bash -c /app/scripts/compile_plugins.sh

FROM golang:1.23.1-bookworm as serverBuilder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build --trimpath -o /fedai main.go

FROM debian:stable AS server
WORKDIR /app
RUN apt update && apt install -y ffmpeg
COPY --from=serverBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=serverBuilder /fedai .
COPY --from=serverBuilder /app/scripts/start.sh .
COPY --from=pluginBuilder /app/pl/* ./pl/*

CMD [ "/fedai" ]
