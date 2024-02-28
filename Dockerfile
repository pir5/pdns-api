FROM golang:latest AS build-env

ENV CGO_ENABLED=0
WORKDIR /go/src/pdns-api
COPY . .

RUN mkdir -p /build
RUN go build -a  -ldflags="-s -w -extldflags \"-static\"" -o=/build/pdns-api main.go

FROM alpine:3
# Timezone = Tokyo
RUN apk --no-cache add tzdata zlib && \
    apk add --upgrade --no-cache && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

COPY --from=build-env /build/pdns-api /build/pdns-api
RUN chmod u+x /build/pdns-api

ENTRYPOINT ["/build/pdns-api", "server"]
