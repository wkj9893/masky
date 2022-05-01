FROM golang:alpine AS build
WORKDIR /masky
COPY . .
RUN apk add make curl && make build-client

FROM alpine:latest
WORKDIR /
COPY --from=build /masky/masky-client /
EXPOSE 1080 
ENTRYPOINT [ "/masky-client" ]