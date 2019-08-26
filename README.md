# Nacelle Base gRPC Process [![GoDoc](https://godoc.org/github.com/go-nacelle/grpcbase?status.svg)](https://godoc.org/github.com/go-nacelle/grpcbase) [![CircleCI](https://circleci.com/gh/go-nacelle/grpcbase.svg?style=svg)](https://circleci.com/gh/go-nacelle/grpcbase) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/grpcbase/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/grpcbase?branch=master)

Abstract gRPC server process for nacelle.

---

This library supplies an abstract gRPC server [process](https://nacelle.dev/docs/core/process) whose behavior can be configured by implementing a `ServerInitializer` interface. For a more full-featured gRPC server framework built on nacelle, see [scarf](/docs/frameworks/scarf).

You can see an additional example of a gRPC process in the [example repository](https://github.com/go-nacelle/example), specifically the [server initializer](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/cmd/grpc-api/server_initializer.go#L21).

### Process

A gRPC process is created by supplying an initializer, described [below](https://nacelle.dev/docs/base-processes/grpcbase#server-initializer), that controls its behavior.

```go
server := grpcbase.NewServer(NewServerInitializer(), options...)
```

### Server Initializer

A server initializer is a struct with an `Init` method that takes a config object and an [grpc.Server](https://godoc.org/google.golang.org/grpc#Server) as parameters.  This method may return an error value, which signals a fatal error to the process that runs it. This method provides an extension point to register services to the server instance before the process accepts connections.

The following example registers a gRPC service to the server that will atomically increment a request counter and return it in a payload defined in the `proto` package that also contains the service definition.

```go
type Initializer struct {}

func (i *Initializer) Init(config nacelle.Config, server *http.Server) error {
    proto.RegisterRequestCounterServiceServer(server, &RequestCounterService{})
    return nil
}

type RequestCounterService {
    requests uint
}

func (kvs *RequestCounterService) Get(ctx context.Context, r *proto.Request) (*proto.Response, error) {
    value := atomic.AddUint32(&i.requests)
    return &proto.Response{count: value}, nil
}
```

#### Initializer Function

A simple server initializer stuct that does not need additional methods, state, or dependency instances injected via a service container can use the server initializer function wrapper instead.

```go
ServerInitializerFunc(func(config nacelle.Config, server *grpc.Server) error {
    proto.RegisterRequestCounterServiceServer(server, &RequestCounterService{})
    return nil
})
```

### Server Process Options

The following options can be supplied to the server constructor to tune its behavior.

<dl>
  <dt>WithTagModifiers</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/grpcbase#WithTagModifiers">WithTagModifiers</a> registers the tag modifiers to be used when loading process configuration (see <a href="https://godoc.org/github.com/go-nacelle/grpcbase#Configuration">below</a>). This can be used to change default hosts and ports, or prefix all target environment variables in the case where more than one gRPC server is registered per application (e.g. health server and application server, data plane and control plane server).</dd>

  <dt>WithServerOptions</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/grpcbase#WithServerOptions">WithServerOptions</a> registers options to be supplied directly to the gRPC server constructor.</dd>
</dl>

### Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Default | Description |
| -------------------- | ------- | ----------- |
| GRPC_HOST            | 0.0.0.0 | The host on which to accept connections. |
| GRPC_PORT            | 5000    | The port on which to accept connections. |
