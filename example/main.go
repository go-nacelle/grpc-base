package main

import (
	"context"
	"strings"

	"github.com/go-nacelle/grpcbase"
	"github.com/go-nacelle/grpcbase/internal"
	"github.com/go-nacelle/nacelle"
	"google.golang.org/grpc"
)

type ServerInitializer struct{}

func (si *ServerInitializer) Init(config nacelle.Config, server *grpc.Server) error {
	internal.RegisterTestServiceServer(server, &upperService{})
	return nil
}

type upperService struct{}

func (us *upperService) ToUpper(ctx context.Context, r *internal.UpperRequest) (*internal.UpperResponse, error) {
	return &internal.UpperResponse{Text: strings.ToUpper(r.GetText())}, nil
}

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	processes.RegisterProcess(grpcbase.NewServer(&ServerInitializer{}))
	return nil
}

func main() {
	nacelle.NewBootstrapper("grpcbase-example", setup).BootAndExit()
}
