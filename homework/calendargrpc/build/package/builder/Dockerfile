# Build applications from sources
FROM golang:alpine
MAINTAINER Artem Orlov <artem.orlov@gmail.com>

# Update packages and install dependecies
RUN apk update && apk upgrade \
    && apk add ca-certificates && update-ca-certificates \
    && apk add openssl \
    && apk add git \
    && rm -rf /var/cache/apk/*

# Get sources from repository
RUN go get -v github.com/AcroManiac/otus-go/...
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc
RUN go mod vendor

# Build Calendar API application
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/cmd/calendar_api
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/calendar_api

# Build Scheduler application
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/cmd/calendar_scheduler
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/calendar_scheduler

# Build Sender application
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/cmd/calendar_sender
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/calendar_sender
