FROM golang:1.23 AS base
ARG VERSION=latest

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive
ENV VERSION=$VERSION

WORKDIR /go/app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0

RUN echo "Building version: ${VERSION}"

# Application dependencies
COPY . .

RUN go mod download \
    && go mod verify

ENV APP_ENV=docker
RUN go build -o app

### Production
FROM alpine:3.16
WORKDIR /usr/local/bin

ARG VERSION=latest
ENV VERSION=$VERSION

# Setup non-root user
ENV user=swadm
ENV gid=1001
ENV uid=1001

RUN addgroup --gid 1001 $user &&  \
    adduser --disabled-password --uid $uid -G $user $user

# Copy executable
COPY --from=base --chown=$user:$user /go/app/app /usr/local/bin/app

USER $user
ENV APP_ENV=dev

EXPOSE ${PORT} 
CMD ["./app"]
