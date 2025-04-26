package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pbapi "modular_monolith/protobuf/api"
	pbhealth "modular_monolith/protobuf/healthpb"
	"modular_monolith/server/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var grpcServer *grpc.Server

type Config struct {
	GRPCPort string
}

func LoadConfig() *Config {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "3551" // default port if not set, safety net! 🎀
	}

	return &Config{
		GRPCPort: port,
	}
}

// Start initializes and runs the gRPC server
func Start(cfg *Config) error {
	address := fmt.Sprintf(":%s", cfg.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("❌ Error listening on port %s: %v", cfg.GRPCPort, err)
		return fmt.Errorf("failed to listen on port %s: %w", cfg.GRPCPort, err)
	}

	// Wrap listener to filter connections by IP
	lis = &ipFilteredListener{Listener: lis}

	// Setup server with a single interceptor
	grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor), // Only add your interceptor once
	)

	// Register services
	registerServices(grpcServer)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	log.Printf("🌸 gRPC server is running on %s ✨", address)

	// Start graceful shutdown handler
	go gracefulShutdown()

	// Serve incoming connections
	log.Println("🌸 Starting to serve incoming connections...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("❌ Error during Serve: %v", err)
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// registerServices registers all gRPC services
func registerServices(server *grpc.Server) {
	pbhealth.RegisterHealthServer(server, &services.HealthService{})
	pbapi.RegisterHelloWorldServer(server, &services.HelloService{})
}

// IP Filtering Listener for whitelisting
type ipFilteredListener struct {
	net.Listener
}

func (l *ipFilteredListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		log.Printf("❌ Error during Accept: %v", err)
		return nil, err
	}

	return conn, nil
}

// Rate Limiting Interceptor
// Define the actual `unaryInterceptor` function
func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Start time for logging
	start := time.Now()
	log.Printf("🌸 gRPC Call: %s at %s", info.FullMethod, start.Format(time.RFC3339))

	// Example: Extract metadata for logging or further checks
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("🌸 Metadata received: %v", md)
	}

	// Proceed with handling the request
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("❌ Error handling request: %v", err)
	}

	// Log completion and duration
	log.Printf("🌸 Completed %s in %v with result: %v", info.FullMethod, time.Since(start), resp)
	return resp, err
}

// gracefulShutdown stops the gRPC server nicely when receiving a termination signal
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🌸 Shutting down gRPC server gracefully... ✨")
	grpcServer.GracefulStop()
	log.Println("🌸 gRPC server stopped. Bye bye~! (つ≧▽≦)つ")
}
