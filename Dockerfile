# syntax=docker/dockerfile:1

FROM golang:latest AS build-binary

WORKDIR "/go/src/github.com/lukasbischof/luk4s.dev/"

COPY ["go.mod", "go.sum", "*.go", "./"]
RUN go mod download && go mod verify
COPY app ./app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o /go/bin/luk4s.dev .

FROM node:latest AS build-assets

WORKDIR "/var/build"

COPY assets/ ./assets
COPY bin/ ./bin
COPY public/ ./public
COPY ["package.json", "package-lock.json", "./"]
RUN npm install
RUN npm run build

FROM alpine:3.15.4

MAINTAINER Lukas Bischof

RUN apk --no-cache add ca-certificates
WORKDIR /opt/luk4s.dev
COPY --from=build-binary /go/bin/luk4s.dev ./
COPY --from=build-assets /var/build/public ./public
COPY "views/" "./views"

ENV HCAPTCHA_SITE_KEY="${HCAPTCHA_SITE_KEY}"
ENV HCAPTCHA_SECRET_KEY="${HCAPTCHA_SECRET_KEY}"

EXPOSE 3000
CMD ["sh", "-c", "/opt/luk4s.dev/luk4s.dev"]
