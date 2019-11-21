FROM golang:1.13 AS builder
COPY . /app
WORKDIR /app
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/cron cmd/cron/cron.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -tags=jsoniter -o /bin/server cmd/server/server.go

FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/resource /resource/
COPY --from=builder /bin/cron /bin/cron
COPY --from=builder /bin/server /bin/server
RUN addgroup -S appuser && adduser -S -G appuser appuser
USER appuser