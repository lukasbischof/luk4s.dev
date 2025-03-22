# syntax=docker/dockerfile:1

FROM golang:latest AS build-binary

WORKDIR "/go/src/github.com/lukasbischof/luk4s.dev/"

COPY ["go.mod", "go.sum", "*.go", "./"]
RUN go mod download && go mod verify
COPY app ./app
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -v -o /go/bin/luk4s.dev .

FROM oven/bun:1 AS build-assets

WORKDIR /usr/src/app

COPY assets/ ./assets
COPY bin/ ./bin
COPY public/ ./public
COPY ["package.json", "bun.lockb", "./"]
RUN bun install --frozen-lockfile --production
RUN bun run build

FROM alpine:3.18.0

LABEL org.opencontainers.image.authors="Lukas Bischof <me@luk4s.dev>"
LABEL org.opencontainers.image.url="https://github.com/lukasbischof/luk4s.dev"

RUN apk --no-cache add ca-certificates
WORKDIR /opt/luk4s.dev
COPY --from=build-binary /go/bin/luk4s.dev ./
COPY --from=build-assets /usr/src/app/public ./public
COPY "views/" "./views"

EXPOSE 3000
CMD ["sh", "-c", "/opt/luk4s.dev/luk4s.dev"]
