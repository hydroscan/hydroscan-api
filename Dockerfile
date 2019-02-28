FROM golang:1.11-alpine AS builder
WORKDIR $home/hydroscan/hydroscan-api
COPY . .
# install gcc for ethereum
RUN apk add --no-cache gcc musl-dev
# Build the binary using go modules vendor
RUN go build -mod=vendor -o /go/bin/subscriber cmd/subscriber/subscriber.go

FROM scratch
COPY --from=builder /go/bin/subscriber /go/bin/subscriber
