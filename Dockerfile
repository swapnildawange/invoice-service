FROM  golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN  CGO_ENABLED=0 && go build -o invoiceApp ./cmd

RUN chmod +x /app/invoiceApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY invoiceApp /app

CMD [ "/app/invoiceApp"]