FROM  golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN  CGO_ENABLED=0  go build  -o invoiceApp ./cmd

RUN chmod +x /app/invoiceApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

# COPY invoiceApp /app
COPY  --from=builder /app/invoiceApp /app
COPY --from=builder /app/env /app

RUN cd /app && ls -al

# RUN cd /app && ls
# RUN ldd /app/invoiceApp
# RUN cd /lib64/ld-linux-x86-64.so.2

CMD [ "/app/invoiceApp"]