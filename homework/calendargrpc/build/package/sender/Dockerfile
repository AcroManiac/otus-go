# Build application from sources
FROM golang:latest as builder
MAINTAINER Artem Orlov <artem.orlov@gmail.com>

# Get sources from repository
RUN go get -v github.com/AcroManiac/otus-go/...
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/cmd/calendar_sender
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/calendar_sender

# Build docker image for application
FROM alpine:latest
RUN apk update && apk upgrade \
    && apk add ca-certificates && update-ca-certificates \
    && apk add openssl

# Add dokerize utility
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# Set Moscow timezone
RUN apk add --update tzdata
ENV TZ=Europe/Moscow

# Clean APK cache
RUN rm -rf /var/cache/apk/*

WORKDIR /calendargrpc
ENV GOPATH /go
COPY --from=builder $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/bin/calendar_sender .
CMD ["dockerize", "-wait", "tcp://postgres:5432", "-wait", "tcp://rabbitmq:15672", "-timeout", "10s", "-wait-retry-interval", "5s", "./calendar_sender"]
