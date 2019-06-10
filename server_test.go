package grpcbase

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"

	"github.com/go-nacelle/grpcbase/internal"
	"github.com/go-nacelle/nacelle"
)

type ServerSuite struct{}

func (s *ServerSuite) TestServeAndStop(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		internal.RegisterTestServiceServer(server, &upperService{})

		return nil
	})

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(&Config{}))
	Expect(err).To(BeNil())

	go server.Start()
	defer server.Stop()

	// Hack internals to get the dynamic port (don't bind to one on host)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", getDynamicPort(server.listener)), grpc.WithInsecure())
	Expect(err).To(BeNil())
	defer conn.Close()

	client := internal.NewTestServiceClient(conn)

	resp, err := client.ToUpper(context.Background(), &internal.UpperRequest{Text: "foobar"})
	Expect(err).To(BeNil())
	Expect(resp.GetText()).To(Equal("FOOBAR"))
}

func (s *ServerSuite) TestBadInjection(t sweet.T) {
	server := NewServer(&badInjectionInitializer{})
	server.Services = makeBadContainer()
	server.Health = nacelle.NewHealth()

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(&Config{}))
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *ServerSuite) TestInitError(t sweet.T) {
	server := makeGRPCServer(func(config nacelle.Config, server *grpc.Server) error {
		return fmt.Errorf("utoh")
	})

	os.Setenv("GRPC_PORT", "0")
	defer os.Clearenv()

	err := server.Init(makeConfig(&Config{}))
	Expect(err).To(MatchError("utoh"))
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
// Service Impl

type upperService struct{}

func (us *upperService) ToUpper(ctx context.Context, r *internal.UpperRequest) (*internal.UpperResponse, error) {
	return &internal.UpperResponse{Text: strings.ToUpper(r.GetText())}, nil
}

//
// Bad Injection

type badInjectionInitializer struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionInitializer) Init(nacelle.Config, *grpc.Server) error {
	return nil
}
