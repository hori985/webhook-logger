FROM golang:1.23.7-alpine3.21 AS builder

RUN mkdir /home/webhook-logger

WORKDIR /home/webhook-logger

COPY . .

RUN go mod tidy && \
    go mod verify   && \
    go mod download && \
    go build -o ./bin/ ./...

FROM alpine:3.21

RUN addgroup -g 9345 webhook-logger && \
    adduser -g 9345 -u 9345 --system --home /home webhook-logger

COPY --from=builder --chown=webhook-logger:webhook-logger /home/webhook-logger/bin/ /home/webhook-logger/bin/

WORKDIR /home/webhook-logger
USER webhook-logger

EXPOSE 8080
ENTRYPOINT [ "bin/webhook-logger" ]