ARG GO_VERSION=1.16.5
ARG NODE_VERSION=14.17.0
ARG ALPINE_VERSION=3.13.5

FROM node:${NODE_VERSION}-alpine AS node-builder
WORKDIR /app
COPY nextjs/package.json nextjs/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY nextjs/ .
ENV NEXT_TELEMETRY_DISABLED=1
RUN yarn run export

FROM golang:${GO_VERSION}-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum main.go ./
COPY --from=node-builder /app/dist ./nextjs/dist

RUN go mod download

# Perform the build
COPY . .

RUN apk update
RUN apk add --no-cache musl-dev gcc ca-certificates

RUN CGO_ENABLED=1 GOOS=linux go build -tags netgo -a -v
#RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w"
#RUN CGO_ENABLED=1 go build --tags "libsqlite3 linux sqlite_fts5"

FROM alpine:${ALPINE_VERSION}
WORKDIR /app
COPY --from=go-builder /app/golang-nextjs-portable .

VOLUME /db

ENTRYPOINT ["./golang-nextjs-portable"]

EXPOSE 8080