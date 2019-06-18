# Nacelle Base gRPC Process [![GoDoc](https://godoc.org/github.com/go-nacelle/grpcbase?status.svg)] [![CircleCI](https://circleci.com/gh/go-nacelle/grpcbase.svg?style=svg)](https://circleci.com/gh/go-nacelle/grpcbase) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/grpcbase/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/grpcbase?branch=master)

Abstract gRPC server process for nacelle.

---

For a more full-featured gRPC server framework built on nacelle, see [scarf](https://github.com/go-nacelle/scarf).

### Usage

The supplied server process is an abstract gRPC server whose behavior is determined by a supplied `ServerInitializer` interface. This interface has only an `Init` method that receives application config as well as the gRPC server instance, allowing server implementations to be registered before the server accepts clients. There is an [example](./example) included in this repository.


The following options can be supplied to the server constructor to tune its behavior.

- **WithTagModifiers** registers the tag modifiers to be used when loading process configuration (see [below](#Configuration)). This can be used to change default hosts and ports, or prefix all target environment variables in the case where more than one HTTP server is registered per application (e.g. health server and application server, data plane and control plane server).
- **WithServerOptions** registers options to be supplied directly to the gRPC server constructor.

### Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Default | Description |
| -------------------- | ------- | ----------- |
| GRPC_HOST            | 0.0.0.0 | The host on which to accept connections. |
| GRPC_PORT            | 5000    | The port on which to accept connections. |















To use the server, initialize a process by passing a Server Initializer to the `NewServer`
constructor. A server initializer is an object with an `Init` method that takes a nacelle
config object (as all process initializer methods do) as well as a `*grpc.Server`. This
*hook* is provided so that services can be registered to the gRPC server before it begins
accepting clients.

The server initializer will have services injected and will receive the nacelle config
object on initialization as if it were a process.

To get a better understanding of the full usage, see the
[example](https://github.com/go-nacelle/tree/master/examples/grpc).

## Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Default | Description |
| -------------------- | ------- | ----------- |
| GRPC_HOST            | 0.0.0.0 | The host on which the server accepts clients. |
| GRPC_PORT            | 6000    | The port on which the server accepts clients. |

## Using Multiple Servers

In order to run multiple gRPC servers, tag modifiers can be applied during config
registration. For more details on how to do this, see the
[example](https://github.com/go-nacelle/tree/master/examples/multi-grpc).

Remember that multiple services can be registered to the same grpc.Server instance, so
multiple processes may not even be necessary depending on your use case.
