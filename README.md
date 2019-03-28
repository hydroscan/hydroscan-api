# HydroScan API

[![CircleCI](https://circleci.com/gh/hydroscan/hydroscan-api.svg?style=svg)](https://circleci.com/gh/hydroscan/hydroscan-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/hydroscan/hydroscan-api)](https://goreportcard.com/report/github.com/hydroscan/hydroscan-api)
[![Uptime Robot status](https://img.shields.io/uptimerobot/status/m782290067-dff6909eff9f905729fcfc92.svg)](https://hydroscan.io)

> The HydroScan backend built with Go. Collects data about Hydro Protocol exchanges and provides methods for a client to consume the data.

## Components

### Subscriber

Subscribes to Hydro Protocol related event logs on the blockchain for up to date trade information.

### Cron

A scheduler to pull a variety of data about running Exchanges, including current token price, trade volume, etc.

### Server

Provides an API for the HydroScan client to query for data.

## Requirements

- Go version >= 1.11
- Postgres
- Redis

## Getting Started

1. Clone this repo to somewhere outside of GOPATH
2. Copy .config.sample.yml to .config.yml and Set valid urls for Postgres and Redis
3. Install the dependencies:

```
go mod download
```

If you get an error like this:

```
go: modules disabled inside GOPATH/src by GO111MODULE=auto; see 'go help modules'
```

It means you have cloned the repo inside GOPATH. To solve this, you can try cloning this repo to somewhere outside of GOPATH. Alternatively, if you want to work in your GOPATH:

```
$ export GO111MODULE=on    # manually active module mode
```

### Run Subscriber

```
go run ./cmd/subscriber/subscriber.go
```

### Run Cron

```
go run ./cmd/cron/cron.go
```

### Run API server

```
go run ./cmd/server/server.go    # Started on port 8080 by default
```

## Contributing

1. Fork it (<https://github.com/hydroscan/hydroscan-api/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request
