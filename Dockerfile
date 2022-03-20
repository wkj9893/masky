FROM golang:alpine AS build
WORKDIR /masky
ENV http_proxy=http://172.20.29.172:1080
ENV https_proxy=http://172.20.29.172:1080
COPY . .
RUN apk add make && make build-server

FROM alpine:latest
WORKDIR /
COPY --from=build /masky/masky-server .
EXPOSE 1080/udp
ENTRYPOINT [ "/masky-server" ]