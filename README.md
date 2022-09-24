# proto-go-sample

[![build sample code](https://github.com/go-training/proto-go-sample/actions/workflows/go.yml/badge.svg)](https://github.com/go-training/proto-go-sample/actions/workflows/go.yml)

Sample Code for [proto connect](https://github.com/bufbuild/connect-go). Connect is a slim library for building browser and gRPC-compatible HTTP APIs.

## Feature

* Support gRPC and HTTP RESTful API using [connect-go](https://github.com/bufbuild/connect-go) library
* Support OpenTelemetry protocol (OTLP) and APM like [uptrace](https://uptrace.dev/) and [signoz](https://signoz.io/)

## Install signoz

Install SigNoz Using Docker Compose

```sh
git clone -b main https://github.com/SigNoz/signoz.git && cd signoz/deploy/
docker-compose -f docker/clickhouse-setup/docker-compose.yaml up -d
```

Ensure that the ports `3301`, `4317` and `4318` are open on the machine where you install SigNoz.

## Build Server and Client

```sh
make build
```

## Start server

Support the following server.

* [Chi](https://github.com/go-chi/chi)
* [Gin](https://github.com/gin-gonic/gin)

Run three service

```sh
# service 01
./bin/gin-server -env-file .env.server01
# service 02
./bin/gin-server -env-file .env.server02
# service 03
./bin/chi-server -env-file .env.server03
```

## Start client

run three client with [Go](https://go.dev)

```sh
./bin/client -env-file .env.server01
./bin/client -env-file .env.server02
./bin/client -env-file .env.server03
```

## Testing with other command

### Use curl

run client with `curl` command

```sh
curl --header "Content-Type: application/json" \
    --data '{"name": "foobar"}' \
    http://localhost:8080/gitea.v1.GiteaService/Gitea
```

health check

```sh
curl --header "Content-Type: application/json" \
    --data '{"service": "gitea.v1.GiteaService"}' \
    http://localhost:8080/grpc.health.v1.Health/Check
```

### Use grpcurl

run client with [grpcurl](https://github.com/fullstorydev/grpcurl) command

```sh
grpcurl \
  -plaintext \
  -d '{"name": "foobar"}' \
  localhost:8080 \
  gitea.v1.GiteaService/Gitea
```

health check

```sh
grpcurl \
  -plaintext \
  -d '{"service": "gitea.v1.GiteaService"}' \
  localhost:8080 \
  grpc.health.v1.Health/Check
```

### Use grpcui

See the [details page for grpcui](https://github.com/fullstorydev/grpcui)

```sh
grpcui -plaintext localhost:8080
```

![page](./images/grpcui01.png)

### Use Postman

![page](./images/postman01.png)
