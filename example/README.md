# gRPC Base Example

A trivial example application to showcase the [grpcbase](https://github.com/go-nacelle/grpcbase) library.

## Overview

This example application uses Redis to provide a simple string get/set API over gRPC. The **main** function boots [nacelle](https://github.com/go-nacelle/nacelle) with a initializer that dials Redis and a server initializer for the process provided by this library. The connection created by the former is injected into the later.

## Building and Running

If running in Docker, simply run `docker-compose up`. This will compile the example application via a multi-stage build and start a container for the API as well as a container for the Redis dependency.

If running locally, simply build with `go build` (using Go 1.12 or above) and invoke with `REDIS_ADDR=redis://{your_redis_host}:6379 ./example`.

For reference, the protobuf files were generated via `protoc -I ./ ./keyvalue.proto --go_out=plugins=grpc:./`.

## Usage

This example comes with a simple looping client that accepts commands of the form `get {key}` and `set {key} {value}`, as followes.

```bash
$ go run client/main.go
> get example-key
payload
> set example-key payload
> get example-key
payload
```
