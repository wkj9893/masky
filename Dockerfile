FROM golang:alpine AS build
WORKDIR /masky
COPY . .
RUN apk add make && make build-server

FROM alpine:latest
WORKDIR /
COPY --from=build /masky/masky-server .
EXPOSE 1080/udp
ENTRYPOINT [ "/masky-server" ]