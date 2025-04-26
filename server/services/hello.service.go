package services

import (
	"context"

	pb "modular_monolith/protobuf/api"
)

type HelloService struct {
	pb.UnimplementedHelloWorldServer
}

func (s *HelloService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello, " + req.Name + "! (⁄ ⁄>⁄ ▽ ⁄<⁄ ⁄)",
	}, nil
}
