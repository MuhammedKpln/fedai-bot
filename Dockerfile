FROM golang:1.23.1-bookworm AS pluginbuilder
WORKDIR /app
COPY . .
RUN go mod download
ENV ENV=PRODUCTION
RUN bash -c /app/scripts/compile_plugins.sh

FROM golang:1.23.1-bookworm AS serverbuilder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build --trimpath -o /fedai main.go

FROM debian:stable AS server
WORKDIR /app
RUN apt update && apt install -y ffmpeg
COPY --from=serverbuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=serverbuilder /fedai .
COPY --from=serverbuilder /app/scripts/start.sh .
COPY --from=pluginbuilder /app/pl/* ./pl/*

CMD [ "/fedai" ]
