package api

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type InProcessHealthClient struct {
	Server health_v1.HealthServer //for grpc type health check
	health_v1.HealthClient  //for http health check
}

func (client *InProcessHealthClient) Check(ctx context.Context, req *health_v1.HealthCheckRequest, opts ...grpc.CallOption) (*health_v1.HealthCheckResponse, error) {
	// we ignore call options since it is in-process
	return client.Server.Check(ctx, req)
}

type HealthServer struct {
	health_v1.UnimplementedHealthServer
}

func (m *HealthServer) Check(ctx context.Context, req *health_v1.HealthCheckRequest) (*health_v1.HealthCheckResponse, error) {
	// we ignore call options since it is in-process
	// return m.Server.Check(ctx, req)
	return &health_v1.HealthCheckResponse{Status: health_v1.HealthCheckResponse_SERVING}, nil
}

func (m *HealthServer) Watch(_ *health_v1.HealthCheckRequest, _ health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}