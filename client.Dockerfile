FROM golang:alpine AS build
WORKDIR /masky
COPY . .
ENV http_proxy=http://172.20.29.172:1080
ENV https_proxy=http://172.20.29.172:1080
RUN apk add make curl && make build-client

FROM alpine:latest
WORKDIR /
COPY --from=build /masky/Country.mmdb /
COPY --from=build /masky/masky-client /
COPY --from=build /masky/web/build /web/build
EXPOSE 1080 1081
ENTRYPOINT [ "/masky-client" ]