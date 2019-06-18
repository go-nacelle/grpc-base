package grpcbase

import (
	"github.com/go-nacelle/nacelle"
	"google.golang.org/grpc"
)

type (
	options struct {
		tagModifiers  []nacelle.TagModifier
		serverOptions []grpc.ServerOption
	}

	// ConfigFunc is a function used to configure an instance of
	// a gRPC Server.
	ConfigFunc func(*options)
)

// WithTagModifiers applies the given tag modifiers on config load.
func WithTagModifiers(modifiers ...nacelle.TagModifier) ConfigFunc {
	return func(o *options) { o.tagModifiers = append(o.tagModifiers, modifiers...) }
}

// WithServerOptions sets gRPC options on the underlying server.
func WithServerOptions(opts ...grpc.ServerOption) ConfigFunc {
	return func(o *options) { o.serverOptions = append(o.serverOptions, opts...) }
}

func getOptions(configs []ConfigFunc) *options {
	options := &options{}
	for _, f := range configs {
		f(options)
	}

	return options
}
