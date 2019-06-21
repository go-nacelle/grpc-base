package main

import (
	"github.com/go-nacelle/grpcbase"
	"github.com/go-nacelle/nacelle"
)

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	processes.RegisterInitializer(NewRedisInitializer(), nacelle.WithInitializerName("redis"))
	processes.RegisterProcess(grpcbase.NewServer(NewServerInitializer()), nacelle.WithProcessName("grpc-server"))
	return nil
}

func main() {
	nacelle.NewBootstrapper("grpcbase-example", setup).BootAndExit()
}
