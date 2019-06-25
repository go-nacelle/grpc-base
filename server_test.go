package grpcbase

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	"github.com/go-nacelle/grpcbase/internal/proto"
	"github.com/go-nacelle/nacelle"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type ServerSuite struct{}

var testConfig = nacelle.NewConfig(nacelle.NewTestEnvSourcer(map[string]string{
	"grpc_port": "0",
}))

func (s *ServerSuite) TestServeAndStop(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		proto.RegisterTestServiceServer(server, &upperService{})
		return nil
	})

	err := server.Init(testConfig)
	Expect(err).To(BeNil())

	go server.Start()
	defer server.Stop()

	// Hack internals to get the dynamic port (don't bind to one on host)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", getDynamicPort(server.listener)), grpc.WithInsecure())
	Expect(err).To(BeNil())
	defer conn.Close()

	client := proto.NewTestServiceClient(conn)

	resp, err := client.ToUpper(context.Background(), &proto.UpperRequest{Text: "foobar"})
	Expect(err).To(BeNil())
	Expect(resp.GetText()).To(Equal("FOOBAR"))
}

func (s *ServerSuite) TestBadInjection(t sweet.T) {
	server := NewServer(&badInjectionInitializer{})
	server.Services = makeBadContainer()
	server.Health = nacelle.NewHealth()

	err := server.Init(testConfig)
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *ServerSuite) TestTagModifiers(t sweet.T) {
	server := NewServer(
		ServerInitializerFunc(func(config nacelle.Config, server *grpc.Server) error {
			return nil
		}),
		WithTagModifiers(nacelle.NewEnvTagPrefixer("prefix")),
	)

	server.Logger = nacelle.NewNilLogger()
	server.Services = nacelle.NewServiceContainer()
	server.Health = nacelle.NewHealth()

	err := server.Init(nacelle.NewConfig(nacelle.NewTestEnvSourcer(map[string]string{
		"prefix_grpc_port": "1234",
	})))

	Expect(err).To(BeNil())
	Expect(server.port).To(Equal(1234))
}

func (s *ServerSuite) TestServerOptions(t sweet.T) {
	handler := NewMockHandler()
	handler.TagRPCFunc.SetDefaultHook(func(ctx context.Context, info *stats.RPCTagInfo) context.Context { return ctx })
	handler.TagConnFunc.SetDefaultHook(func(ctx context.Context, info *stats.ConnTagInfo) context.Context { return ctx })

	server := NewServer(
		ServerInitializerFunc(func(config nacelle.Config, server *grpc.Server) error {
			proto.RegisterTestServiceServer(server, &upperService{})
			return nil
		}),
		WithServerOptions(grpc.StatsHandler(handler)),
	)

	server.Logger = nacelle.NewNilLogger()
	server.Services = nacelle.NewServiceContainer()
	server.Health = nacelle.NewHealth()

	err := server.Init(testConfig)
	Expect(err).To(BeNil())

	go server.Start()
	defer server.Stop()

	// Hack internals to get the dynamic port (don't bind to one on host)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", getDynamicPort(server.listener)), grpc.WithInsecure())
	Expect(err).To(BeNil())
	defer conn.Close()

	client := proto.NewTestServiceClient(conn)

	resp, err := client.ToUpper(context.Background(), &proto.UpperRequest{Text: "foobar"})
	Expect(err).To(BeNil())
	Expect(resp.GetText()).To(Equal("FOOBAR"))

	// Ensure stats handler was registered
	Expect(handler.TagRPCFunc).To(BeCalled())
	Expect(handler.HandleRPCFunc).To(BeCalled())
	Expect(handler.TagConnFunc).To(BeCalled())
	Expect(handler.HandleConnFunc).To(BeCalled())
}

func (s *ServerSuite) TestInitError(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		return fmt.Errorf("oops")
	})

	err := server.Init(testConfig)
	Expect(err).To(MatchError("oops"))
}

//
// Helpers

func makeGRPCServer(initializer func(nacelle.Config, *grpc.Server) error) *Server {
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

func (i *badInjectionInitializer) Init(nacelle.Config, *grpc.Server) error {
	return nil
}

func makeBadContainer() nacelle.ServiceContainer {
	container := nacelle.NewServiceContainer()
	container.Set("A", &B{})
	return container
}
