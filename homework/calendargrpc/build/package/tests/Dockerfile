# Build docker image for Calendar services integration tests
FROM deployments_builder:latest
MAINTAINER Artem Orlov <artem.orlov@gmail.com>

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

# Add Git package
RUN apk add git

# Clean APK cache
RUN rm -rf /var/cache/apk/*

# Set working dir
WORKDIR $GOPATH/src/github.com/AcroManiac/otus-go/homework/calendargrpc/test/integration

## Run integration tests
#ENV CGO_ENABLED=0
#ENV GOOS=linux
#CMD ["dockerize", "-wait", "tcp://postgres:5432", "-wait", "tcp://rabbitmq:15672", "-timeout", "10s", "-wait-retry-interval", "10s", "go", "test"]
