FROM golang:1.23.1-bookworm as pluginBuilder

WORKDIR /app

COPY . .

RUN go build --buildmode=plugin --trimpath -o /help.so /app/pl/help.go
RUN go build --buildmode=plugin --trimpath -o /voicy.so /app/pl/voicy.go
RUN go build --buildmode=plugin --trimpath -o /plugin.so /app/pl/plugin.go
RUN go build --buildmode=plugin --trimpath -o /plugins.so /app/pl/plugins.go


FROM golang:1.23.1-bookworm as serverBuilder
WORKDIR /app
COPY . .
RUN go build --trimpath -o /fedai main.go

FROM debian:stable AS server
WORKDIR /app

RUN apt update && apt install -y ffmpeg
COPY --from=serverBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=pluginBuilder /help.so ./pl/help.so
COPY --from=pluginBuilder /voicy.so ./pl/voicy.so
COPY --from=pluginBuilder /plugin.so ./pl/plugin.so
COPY --from=pluginBuilder /plugins.so ./pl/plugins.so
COPY --from=serverBuilder /fedai .
COPY --from=serverBuilder /app/start.sh .

CMD [ "/fedai" ]