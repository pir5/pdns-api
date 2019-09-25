FROM debian:stable-slim

ADD ./build/linux/amd64/pdns-api /pdns-api

EXPOSE 8080 8080/tcp
ENTRYPOINT ["/pdns-api", "server"]
