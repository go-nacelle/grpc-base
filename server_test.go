package grpcbase

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/go-nacelle/grpcbase/internal/proto"
	"github.com/go-nacelle/nacelle/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

var testConfig = nacelle.NewConfig(nacelle.NewTestEnvSourcer(map[string]string{
	"grpc_port": "0",
}))

func TestServeAndStop(t *testing.T) {
	server := makeGRPCServer(func(ctx context.Context, server *grpc.Server) error {
		proto.RegisterTestServiceServer(server, &upperService{})
		return nil
	})
	server.Config = testConfig

	ctx := context.Background()
	err := server.Init(ctx)
	assert.Nil(t, err)

	go server.Start()
	defer server.Stop()

	// Hack internals to get the dynamic port (don't bind to one on host)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", getDynamicPort(server.listener)), grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()

	client := proto.NewTestServiceClient(conn)

	resp, err := client.ToUpper(context.Background(), &proto.UpperRequest{Text: "foobar"})
	assert.Nil(t, err)
	assert.Equal(t, "FOOBAR", resp.GetText())
}

func TestBadInjection(t *testing.T) {
	server := NewServer(&badInjectionInitializer{})
	server.Services = makeBadContainer()
	server.Health = nacelle.NewHealth()
	server.Config = testConfig

	ctx := context.Background()
	err := server.Init(ctx)
	assert.Contains(t, err.Error(), "ServiceA")
}

func TestTagModifiers(t *testing.T) {
	server := NewServer(
		ServerInitializerFunc(func(ctx context.Context, server *grpc.Server) error {
			return nil
		}),
		WithTagModifiers(nacelle.NewEnvTagPrefixer("prefix")),
	)

	server.Logger = nacelle.NewNilLogger()
	server.Services = nacelle.NewServiceContainer()
	server.Health = nacelle.NewHealth()
	server.Config = nacelle.NewConfig(nacelle.NewTestEnvSourcer(map[string]string{
		"prefix_grpc_port": "1234",
	}))

	ctx := context.Background()
	err := server.Init(ctx)

	assert.Nil(t, err)
	assert.Equal(t, 1234, server.port)
}

func TestServerOptions(t *testing.T) {
	handler := NewMockHandler()
	handler.TagRPCFunc.SetDefaultHook(func(ctx context.Context, info *stats.RPCTagInfo) context.Context { return ctx })
	handler.TagConnFunc.SetDefaultHook(func(ctx context.Context, info *stats.ConnTagInfo) context.Context { return ctx })

	server := NewServer(
		ServerInitializerFunc(func(ctx context.Context, server *grpc.Server) error {
			proto.RegisterTestServiceServer(server, &upperService{})
			return nil
		}),
		WithServerOptions(grpc.StatsHandler(handler)),
	)

	server.Logger = nacelle.NewNilLogger()
	server.Services = nacelle.NewServiceContainer()
	server.Health = nacelle.NewHealth()
	server.Config = testConfig

	ctx := context.Background()
	err := server.Init(ctx)
	assert.Nil(t, err)

	go server.Start()
	defer server.Stop()

	// Hack internals to get the dynamic port (don't bind to one on host)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", getDynamicPort(server.listener)), grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()

	client := proto.NewTestServiceClient(conn)

	resp, err := client.ToUpper(context.Background(), &proto.UpperRequest{Text: "foobar"})
	assert.Nil(t, err)
	assert.Equal(t, "FOOBAR", resp.GetText())

	// Ensure stats handler was registered
	mockassert.Called(t, handler.TagRPCFunc)
	mockassert.Called(t, handler.HandleRPCFunc)
	mockassert.Called(t, handler.TagConnFunc)
	mockassert.Called(t, handler.HandleConnFunc)
}

func TestInitError(t *testing.T) {
	server := makeGRPCServer(func(ctx context.Context, server *grpc.Server) error {
		return fmt.Errorf("oops")
	})
	server.Config = testConfig

	ctx := context.Background()
	err := server.Init(ctx)
	assert.EqualError(t, err, "oops")
}

//
// Helpers

func makeGRPCServer(initializer func(context.Context, *grpc.Server) error) *Server {
	server := NewServer(ServerInitializerFunc(initializer))
	server.Logger = nacelle.NewNilLogger()
	server.Services = nacelle.NewServiceContainer()
	server.Health = nacelle.NewHealth()
	return server
}

func getDynamicPort(listener net.Listener) int {
	return listener.Addr().(*net.TCPAddr).Port
}

//
// Service Implementation

type upperService struct{}

func (us *upperService) ToUpper(ctx context.Context, r *proto.UpperRequest) (*proto.UpperResponse, error) {
	return &proto.UpperResponse{Text: strings.ToUpper(r.GetText())}, nil
}

//
// Bad Injection

type A struct{ X int }
type B struct{ X float64 }

type badInjectionInitializer struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionInitializer) Init(context.Context, *grpc.Server) error {
	return nil
}

func makeBadContainer() *nacelle.ServiceContainer {
	container := nacelle.NewServiceContainer()
	container.Set("A", &B{})
	return container
}
