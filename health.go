package grpcbase

type healthToken string

func (t healthToken) String() string {
	return "grpc-init"
}
