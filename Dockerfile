FROM golang:1.23.7-alpine3.20 as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# RELEASE_VERSION is baked into the binary via ldflags. Override via
# `docker build --build-arg RELEASE_VERSION=2.184.0-yandex.1`.
ARG RELEASE_VERSION=unspecified

RUN apk add --no-cache make git

WORKDIR /go/src/github.com/supabase/auth

# Pulling dependencies
COPY ./Makefile ./go.* ./
RUN make deps

# Building stuff
COPY . /go/src/github.com/supabase/auth

RUN RELEASE_VERSION=${RELEASE_VERSION} make build

# Always use alpine:3 so the latest version is used. This will keep CA certs more up to date.
FROM alpine:3
RUN adduser -D -u 1000 supabase

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/supabase/auth/auth /usr/local/bin/auth
COPY --from=build /go/src/github.com/supabase/auth/migrations /usr/local/etc/auth/migrations/
RUN ln -s /usr/local/bin/auth /usr/local/bin/gotrue

ENV GOTRUE_DB_MIGRATIONS_PATH /usr/local/etc/auth/migrations

USER supabase
CMD ["auth"]
