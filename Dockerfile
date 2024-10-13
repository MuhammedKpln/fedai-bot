FROM golang:latest as pluginBuilder

WORKDIR /app

COPY . .

RUN go build --buildmode=plugin --trimpath -o /help.so /app/pl/help.go
RUN go build --buildmode=plugin --trimpath -o /voicy.so /app/pl/voicy.go


FROM golang:latest as serverBuilder
WORKDIR /app
COPY . .
RUN go build --trimpath -o /fedai main.go

FROM debian:stable AS server
WORKDIR /app

COPY --from=serverBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=pluginBuilder /help.so ./pl/help.so
COPY --from=pluginBuilder /voicy.so ./pl/voicy.so
COPY --from=serverBuilder /fedai .

CMD [ "/fedai" ]