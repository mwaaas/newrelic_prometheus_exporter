FROM golang:1.11.1-alpine3.8 as build-env

RUN apk add --update \
    && apk add ca-certificates git gcc musl-dev \
    && mkdir -p ./build

WORKDIR /usr/src/app

COPY . .

RUN  go build  -o dist/exporter

CMD go run *.go


# running container
FROM alpine:3.8
RUN apk add --update \
      ca-certificates
EXPOSE 80
WORKDIR /app
COPY --from=build-env /usr/src/app/dist/exporter /app/
ENTRYPOINT ./exporter

