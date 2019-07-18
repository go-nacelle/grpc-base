# Nacelle Base gRPC Process [![GoDoc](https://godoc.org/github.com/go-nacelle/grpcbase?status.svg)](https://godoc.org/github.com/go-nacelle/grpcbase) [![CircleCI](https://circleci.com/gh/go-nacelle/grpcbase.svg?style=svg)](https://circleci.com/gh/go-nacelle/grpcbase) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/grpcbase/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/grpcbase?branch=master)

Abstract gRPC server process for nacelle.

---

For a more full-featured gRPC server framework built on nacelle, see [scarf](https://nacelle.dev/docs/frameworks/scarf).

### Usage

The supplied server process is an abstract gRPC server whose behavior is determined by a supplied `ServerInitializer` interface. This interface has only an `Init` method that receives application config as well as the gRPC server instance, allowing server implementations to be registered before the server accepts clients. There is an [example](./example) included in this repository.

The following options can be supplied to the server constructor to tune its behavior.

- **WithTagModifiers** registers the tag modifiers to be used when loading process configuration (see [below](#Configuration)). This can be used to change default hosts and ports, or prefix all target environment variables in the case where more than one gRPC server is registered per application (e.g. health server and application server, data plane and control plane server).
- **WithServerOptions** registers options to be supplied directly to the gRPC server constructor.

### Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Default | Description |
| -------------------- | ------- | ----------- |
| GRPC_HOST            | 0.0.0.0 | The host on which to accept connections. |
| GRPC_PORT            | 5000    | The port on which to accept connections. |
