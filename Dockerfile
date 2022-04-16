# syntax=docker/dockerfile:1

FROM golang:latest AS build

WORKDIR "/go/src/github.com/lukasbischof/luk4s.dev/"

COPY ["go.mod", "go.sum", "*.go", "./"]
RUN go mod download && go mod verify
COPY app ./app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o /go/bin/luk4s.dev .

FROM alpine:3.15.4

MAINTAINER Lukas Bischof

RUN apk --no-cache add ca-certificates
WORKDIR /opt/luk4s.dev
COPY --from=build /go/bin/luk4s.dev ./
COPY "views/" "./views"

EXPOSE 3000
CMD ["/opt/luk4s.dev/luk4s.dev"]
