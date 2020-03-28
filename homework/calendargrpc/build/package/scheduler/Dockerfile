# Build application from sources
FROM golang:latest as builder
LABEL maintainer="Artem Orlov <artem.orlov@gmail.com>"

RUN go get -v github.com/AcroManiac/otus-go/...
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/cmd/calendar_scheduler
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/calendar_scheduler

# Build docker image for application
FROM alpine:latest
WORKDIR /calendargrpc
ENV GOPATH /go
COPY --from=builder $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/bin/calendar_scheduler .
COPY --from=builder $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/configs/calendar_scheduler.yaml .
CMD ["./calendar_scheduler --config=calendar_scheduler.yaml"]
