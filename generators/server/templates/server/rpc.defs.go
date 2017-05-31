package server

import (
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"

	pb "jaxf-github.fanatics.corp/apparel/<%= appname %>/protocol"
)

type server struct{}

func (s *server) SayHello(c context.Context, m *pb.HelloRequest) (*pb.HelloReply, error) {
	grpclog.Printf("rpc request SayHello(%q)\n", m.Name)
	return &pb.HelloReply{Message: "Hello " + m.Name}, nil
}

func (s *server) Version(c context.Context, m *pb.Empty) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{Version: os.Getenv("VERSION")}, nil
}

func newServer() *server {
	return new(server)
}
