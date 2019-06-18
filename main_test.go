package grpcbase

import (
	"testing"

	"github.com/aphistic/sweet"
	junit "github.com/aphistic/sweet-junit"
	"github.com/go-nacelle/nacelle"
	"github.com/go-nacelle/nacelle/mocks"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ServerSuite{})
	})
}

//
// Config

type emptyConfig struct{}

func makeConfig(base *Config) nacelle.Config {
	cfg := mocks.NewMockConfig()
	cfg.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...nacelle.TagModifier) error {
		c := target.(*Config)
		c.GRPCPort = base.GRPCPort
		return nil
	})

	return cfg
}

//
//  Injection

type A struct{ X int }
type B struct{ X float64 }

func makeBadContainer() nacelle.ServiceContainer {
	container := nacelle.NewServiceContainer()
	container.Set("A", &B{})
	return container
}
