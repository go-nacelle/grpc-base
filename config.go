package grpcbase

type Config struct {
	GRPCHost string `env:"grpc_host" default:"0.0.0.0"`
	GRPCPort int    `env:"grpc_port" default:"5000"`
}
