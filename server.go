package grpcbase

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/go-nacelle/nacelle/v2"
	"github.com/go-nacelle/process/v2"
	"github.com/go-nacelle/service/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type (
	Server struct {
		Config        *nacelle.Config           `service:"config"`
		Logger        nacelle.Logger            `service:"logger"`
		Services      *nacelle.ServiceContainer `service:"services"`
		Health        *nacelle.Health           `service:"health"`
		tagModifiers  []nacelle.TagModifier
		initializer   ServerInitializer
		listener      *net.TCPListener
		server        *grpc.Server
		once          *sync.Once
		stopped       chan struct{}
		host          string
		port          int
		serverOptions []grpc.ServerOption
		healthToken   healthToken
		healthStatus  *process.HealthComponentStatus
	}

	ServerInitializer interface {
		Init(context.Context, *grpc.Server) error
	}

	ServerInitializerFunc func(context.Context, *grpc.Server) error
)

func (f ServerInitializerFunc) Init(ctx context.Context, server *grpc.Server) error {
	return f(ctx, server)
}

func NewServer(initializer ServerInitializer, configs ...ConfigFunc) *Server {
	options := getOptions(configs)

	return &Server{
		tagModifiers:  options.tagModifiers,
		initializer:   initializer,
		once:          &sync.Once{},
		stopped:       make(chan struct{}),
		serverOptions: options.serverOptions,
		healthToken:   healthToken(uuid.New().String()),
	}
}

func (s *Server) Init(ctx context.Context) (err error) {
	healthStatus, err := s.Health.Register(s.healthToken)
	if err != nil {
		return err
	}
	s.healthStatus = healthStatus

	grpcConfig := &Config{}
	if err = s.Config.Load(grpcConfig, s.tagModifiers...); err != nil {
		return err
	}

	s.listener, err = makeListener(grpcConfig.GRPCHost, grpcConfig.GRPCPort)
	if err != nil {
		return
	}

	if err := service.Inject(ctx, s.Services, s.initializer); err != nil {
		return err
	}

	s.host = grpcConfig.GRPCHost
	s.port = grpcConfig.GRPCPort
	s.server = grpc.NewServer(s.serverOptions...)
	err = s.initializer.Init(ctx, s.server)
	return
}

func (s *Server) Start(ctx context.Context) error {
	defer s.listener.Close()

	s.healthStatus.Update(true)

	s.Logger.Info("Serving gRPC on %s:%d", s.host, s.port)

	if err := s.server.Serve(s.listener); err != nil {
		select {
		case <-s.stopped:
		default:
			return err
		}
	}

	s.Logger.Info("No longer serving gRPC on %s:%d", s.host, s.port)
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.once.Do(func() {
		s.Logger.Info("Shutting down gRPC server")
		close(s.stopped)
		s.server.GracefulStop()
	})

	return nil
}

func makeListener(host string, port int) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}
