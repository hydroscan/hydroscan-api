FROM golang:1.11-alpine AS builder
WORKDIR $home/hydroscan/hydroscan-api
COPY . .
# install gcc for ethereum
RUN apk add --no-cache gcc musl-dev
# Build the binary using go modules vendor
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod=vendor -o /bin/subscriber cmd/subscriber/subscriber.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod=vendor -o /bin/cron cmd/cron/cron.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod=vendor -o /bin/server cmd/server/server.go

FROM scratch
COPY --from=builder /bin/subscriber /bin/subscriber
COPY --from=builder /bin/cron /bin/cron
COPY --from=builder /bin/server /bin/server