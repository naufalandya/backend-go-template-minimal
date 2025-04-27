package client

import (
	"log"
	"modular_monolith/protobuf/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Global variable to store the gRPC client
var HelloWorldService api.HelloWorldClient

// Connect to the gRPC service and initialize the client
func Connect() error {
	conn, err := grpc.NewClient("localhost:3551", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}

	HelloWorldService = api.NewHelloWorldClient(conn)
	return nil
}

// // ServiceClients holds all your service clients
// type ServiceClients struct {
// 	HelloWorldClient pbapi.HelloWorldClient
// 	UserServiceClient pbapi.UserServiceClient  // Another service
// 	ProductServiceClient pbapi.ProductServiceClient // Example for Product Service
// }

// // Connect initializes and returns service clients
// func Connect() (*ServiceClients, error) {
// 	// First connection to HelloWorld service
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Printf("Failed to connect to HelloWorld service: %v", err)
// 		return nil, err
// 	}

// 	// You can add more connections for different services (e.g., UserService, ProductService)
// 	connUser, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Printf("Failed to connect to UserService: %v", err)
// 		return nil, err
// 	}

// 	connProduct, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Printf("Failed to connect to ProductService: %v", err)
// 		return nil, err
// 	}

// 	// Returning a struct containing all clients
// 	return &ServiceClients{
// 		HelloWorldClient: pbapi.NewHelloWorldClient(conn),
// 		UserServiceClient: pbapi.NewUserServiceClient(connUser),
// 		ProductServiceClient: pbapi.NewProductServiceClient(connProduct),
// 	}, nil
// }
