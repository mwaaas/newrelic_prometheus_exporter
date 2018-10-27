FROM golang:1.11.0-stretch as builder

WORKDIR /usr/src/app
COPY . .

RUN mkdir -p ./build \
    && go build  -o dist/exporter

CMD go run *.go


## running container
#FROM alpine
#
#EXPOSE 4100
#
#COPY --from=builder /usr/src/app/dist/exporter /
#ENTRYPOINT ["/exporter"]
