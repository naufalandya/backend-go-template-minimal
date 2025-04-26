package services

import (
	"context"

	pb "modular_monolith/protobuf/healthpb"
)

type HealthService struct {
	pb.UnimplementedHealthServer
}

func (s *HealthService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "Healthy nya~! (๑˃ᴗ˂)ﻭ",
	}, nil
}
