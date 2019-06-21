package main

import (
	"context"
	"example/proto"
	"github.com/garyburd/redigo/redis"
	"github.com/go-nacelle/nacelle"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KeyValueService struct {
	logger nacelle.Logger
	redis  redis.Conn
}

func NewKeyValueService(logger nacelle.Logger, redis redis.Conn) *KeyValueService {
	return &KeyValueService{
		logger: logger,
		redis:  redis,
	}
}

func (kvs *KeyValueService) Get(ctx context.Context, r *proto.GetRequest) (*proto.GetResponse, error) {
	reply, err := redis.String(kvs.redis.Do("GET", r.GetKey()))
	if err != nil {
		if err == redis.ErrNil {
			return nil, status.Error(codes.NotFound, "key is not set")
		}

		kvs.logger.Error("Failed to perform GET (%s)", err.Error())
		return nil, err
	}

	kvs.logger.Debug("Retrieved key %s", r.GetKey())
	return &proto.GetResponse{Value: reply}, nil
}

func (kvs *KeyValueService) Set(ctx context.Context, r *proto.SetRequest) (*proto.SetResponse, error) {
	if _, err := kvs.redis.Do("SET", r.GetKey(), r.GetValue()); err != nil {
		kvs.logger.Error("Failed to perform SET (%s)", err.Error())
		return nil, err
	}

	kvs.logger.Debug("Set key %s", r.GetKey())
	return &proto.SetResponse{}, nil
}
