FROM golang:1.23-alpine as builder

ARG REVISION

WORKDIR /ava

COPY . .

RUN go mod download

RUN go run github.com/steebchen/prisma-client-go generate

RUN CGO_ENABLED=0 go build -ldflags "-s -w \
    -X github.com/matthisholleville/ava/pkg/version.REVISION=${REVISION}" \
    -a -o bin/ava

FROM alpine:3.20@sha256:77726ef6b57ddf65bb551896826ec38bc3e53f75cdde31354fbffb4f25238ebd

ARG BUILD_DATE
ARG VERSION
ARG REVISION

LABEL maintainer="matthis.holleville"

RUN addgroup -S app \
    && adduser -S -G app app \
    && apk --no-cache add \
    curl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /ava/bin/ava .
RUN chown -R app:app ./

WORKDIR /home/app/config

COPY --from=builder /ava/config .

WORKDIR /home/app

USER app

EXPOSE 8080

ENTRYPOINT ["./ava"]