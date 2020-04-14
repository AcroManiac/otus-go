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
COPY --from=deployments_builder:latest $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/bin/calendar_api .
EXPOSE 8888
CMD ["dockerize", "-wait", "tcp://postgres:5432", "-timeout", "10s", "-wait-retry-interval", "5s", "./calendar_api"]
