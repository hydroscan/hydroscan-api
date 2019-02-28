FROM golang:1.11-alpine AS builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
# install gcc for ethereum
RUN apk add --no-cache gcc musl-dev
# Build the binary using go modules vendor
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod=vendor -o /bin/subscriber cmd/subscriber/subscriber.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod=vendor -o /bin/cron cmd/cron/cron.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod=vendor -o /bin/server cmd/server/server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder /bin/subscriber /bin/subscriber
COPY --from=builder /bin/cron /bin/cron
COPY --from=builder /bin/server /bin/server