# Build application from sources
FROM golang:latest as builder
LABEL maintainer="Artem Orlov <artem.orlov@gmail.com>"

RUN go get -v github.com/AcroManiac/otus-go/...
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/cmd/calendar_sender
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/calendar_sender

# Build docker image for application
FROM alpine:latest
WORKDIR /calendargrpc
ENV GOPATH /go
COPY --from=builder $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/bin/calendar_sender .
CMD ["./calendar_sender"]
