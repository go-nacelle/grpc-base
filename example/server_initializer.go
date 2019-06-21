package main

import (
	"example/proto"
	"github.com/garyburd/redigo/redis"
	"github.com/go-nacelle/grpcbase"
	"github.com/go-nacelle/nacelle"
	"google.golang.org/grpc"
)

type ServerInitializer struct {
	Logger nacelle.Logger `service:"logger"`
	Redis  redis.Conn     `service:"redis"`
}

func NewServerInitializer() grpcbase.ServerInitializer {
	return &ServerInitializer{}
}

func (si *ServerInitializer) Init(config nacelle.Config, server *grpc.Server) error {
	proto.RegisterKeyValueServiceServer(server, NewKeyValueService(si.Logger, si.Redis))
	return nil
}
